// A in-memory KV store with WAL(write-ahead log)
// This is the data plane
package worker

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eyeKill/KV/common"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

const (
	LOG_FILENAME              = "log.txt"
	SLOT_FILENAME             = "slots.json"
	SLOT_TMP_FILENAME_PATTERN = "slots.*.json"
	LOG_TMP_FILENAME_PATTERN  = "log.*.txt"
)

var (
	ENOENT    = errors.New("entry does not exist")
	EINVTRANS = errors.New("invalid transaction id")
	ENOTRANS  = errors.New("no available transaction id, try again later")
)

const (
	TRANSACTION_COUNT = 16
)

// interface for a kv store
type KVStore interface {
	Get(key string, transactionId int) (value string, err error)
	Put(key string, value string, transactionId int) (version uint64, err error)
	Delete(key string, transactionId int) (version uint64, err error)
	StartTransaction() (transactionId int, err error)
	Rollback(transactionId int) error
	Commit(transactionId int) error
	// persist kv store
	Flush()
	Checkpoint() error
	PrepareExtract()
	// Extract all values for keys that satisfies the divider function at the time this method is called.
	// This method should not block. When doing calculation, the KVStore should continue to serve on other threads.
	Extract(divider func(key string) bool) map[string]string
}

type TransactionStruct struct {
	Lock  sync.RWMutex
	Layer map[string]*string
}

// and our implementation
// KV map is separated into two kind of maps. Base map records those contained in "slots.json",
// and latest maps contains those indicated by WAL. The array of latest maps forms a log-like data structure,
// and provide atomic undo/redo. Note that values could have been moved, so latest map could contain nil values.
type SimpleKV struct {
	// permutation is: base <- layers[0] <- layers[1] <-...<- lastLayer
	// base and layers can only be read, lastLayer can be read & write
	// so only lastLayer is locked with RWMutex
	base           map[string]string
	layers         []map[string]*string
	transactions   []*TransactionStruct
	tLock          sync.RWMutex // for usedTransactions array
	checkpointLock sync.Mutex
	path           string
	version        uint64
	logFile        *os.File
}

func (kv *SimpleKV) getTransaction(transactionId int) *TransactionStruct {
	if transactionId < 0 || transactionId >= TRANSACTION_COUNT {
		return nil
	}
	kv.tLock.RLock()
	defer kv.tLock.RUnlock()
	return kv.transactions[transactionId]
}

func (kv *SimpleKV) Get(key string, transactionId int) (string, error) {
	// get does not require logging
	// lookup layer by layer
	t := kv.getTransaction(transactionId)
	if t == nil {
		return "", EINVTRANS
	}
	t.Lock.RLock()
	if v, ok := t.Layer[key]; ok {
		t.Lock.RUnlock()
		if v == nil {
			return "", ENOENT
		} else {
			return *v, nil
		}
	}
	t.Lock.RUnlock()

	// go through layers
	for i := range kv.layers {
		ii := len(kv.layers) - i - 1
		if v, ok := kv.layers[ii][key]; ok {
			if v == nil {
				return "", ENOENT
			} else {
				return *v, nil
			}
		}
	}
	value, ok := kv.base[key]
	if ok {
		return value, nil
	} else {
		return "", ENOENT
	}
}

func (kv *SimpleKV) Put(key string, value string, transactionId int) (uint64, error) {
	t := kv.getTransaction(transactionId)
	if t == nil {
		return 0, EINVTRANS
	}
	t.Lock.Lock()
	defer t.Lock.Unlock()
	if transactionId == 0 {
		kv.version += 1
		kv.writeLog("put", key, value, "0", strconv.FormatUint(kv.version, 16))
		t.Layer[key] = &value
		return kv.version, nil
	} else {
		kv.writeLog("put", key, value, strconv.FormatInt(int64(transactionId), 16))
		t.Layer[key] = &value
		return 0, nil
	}
}

// Make sure that key is removed from KVStore, regardless of whether it exists beforehand or not.
// You should check if the key exists in the KV beforehand, otherwise this API could thrash the KV.
func (kv *SimpleKV) Delete(key string, transactionId int) (uint64, error) {
	t := kv.getTransaction(transactionId)
	if t == nil {
		return 0, EINVTRANS
	}
	t.Lock.Lock()
	defer t.Lock.Unlock()
	if transactionId == 0 {
		kv.version += 1
		kv.writeLog("del", key, "0", strconv.FormatUint(kv.version, 16))
		t.Layer[key] = nil
		return kv.version, nil
	} else {
		kv.writeLog("del", key, strconv.FormatInt(int64(transactionId), 16))
		t.Layer[key] = nil
		return 0, nil
	}
}

func (kv *SimpleKV) StartTransaction() (transactionId int, err error) {
	// find a valid transaction id
	kv.tLock.Lock()
	defer kv.tLock.Unlock()
	for i, u := range kv.transactions {
		if u == nil {
			kv.writeLog("start", strconv.FormatInt(int64(transactionId), 16))
			kv.transactions[i] = &TransactionStruct{
				Lock:  sync.RWMutex{},
				Layer: make(map[string]*string),
			}
			return i, nil
		}
	}
	return 0, ENOTRANS
}

func (kv *SimpleKV) Rollback(transactionId int) error {
	if transactionId < 1 || transactionId >= TRANSACTION_COUNT {
		return EINVTRANS
	}
	kv.tLock.Lock()
	defer kv.tLock.Unlock()
	if kv.transactions[transactionId] != nil {
		kv.writeLog("rollback", strconv.FormatInt(int64(transactionId), 16))
		kv.transactions[transactionId] = nil
		return nil
	} else {
		return EINVTRANS
	}
}

// transaction zero can also be committed, but committing zero does not perform any operation
func (kv *SimpleKV) Commit(transactionId int) error {
	if transactionId < 0 || transactionId >= TRANSACTION_COUNT {
		return EINVTRANS
	}
	if transactionId == 0 {
		return nil
	}
	kv.tLock.RLock()
	t := kv.transactions[transactionId]
	kv.tLock.RUnlock()
	if t != nil {
		// merge it into transaction zero
		kv.transactions[0].Lock.Lock()
		kv.writeLog("commit", strconv.FormatInt(int64(transactionId), 16))
		t.Lock.RLock()
		for k, v := range t.Layer {
			kv.transactions[0].Layer[k] = v
		}
		t.Lock.RUnlock()
		kv.transactions[0].Lock.Unlock()
		// remove this transaction
		kv.tLock.Lock()
		kv.transactions[transactionId] = nil
		kv.tLock.Unlock()
		return nil
	} else {
		return EINVTRANS
	}
}

// squish last layer into layers array
// NOTICE you'll have to manually handle locking when calling this method
func (kv *SimpleKV) saveLayer() {
	kv.tLock.Lock()
	defer kv.tLock.Unlock()
	t := kv.transactions[0]
	t.Lock.Lock()
	defer t.Lock.Unlock()
	if len(t.Layer) == 0 {
		return
	}
	kv.layers = append(kv.layers, t.Layer)
	t.Layer = make(map[string]*string)
}

// Clear log entries, flush current kv in memory to slots.
// Call this when log file is getting too large.
func (kv *SimpleKV) Checkpoint() error {
	// only one checkpoint operation can happen at a time
	kv.checkpointLock.Lock()
	defer kv.checkpointLock.Unlock()
	// make current last layer immutable
	kv.saveLayer()
	// calculate new base
	b := make(map[string]string)
	for k, v := range kv.base {
		b[k] = v
	}
	for _, l := range kv.layers {
		for k, v := range l {
			if v == nil {
				delete(b, k)
			} else {
				b[k] = *v
			}
		}
	}
	// and add the newest layer to it...
	// essentially, we are collecting the new changes made during calculation of new base
	// lock both read & write
	t := kv.transactions[0]
	t.Lock.Lock()
	for k, v := range t.Layer {
		if v == nil {
			delete(b, k)
		} else {
			b[k] = *v
		}
	}
	// redirect log to another temporary file
	tmpLogFile, err := ioutil.TempFile(kv.path, LOG_TMP_FILENAME_PATTERN)
	if err != nil {
		t.Lock.Unlock()
		return err
	}
	if err := kv.logFile.Close(); err != nil {
		t.Lock.Unlock()
		return err
	}
	kv.logFile = tmpLogFile
	t.Lock.Unlock()
	// COMMIT POINT

	// flush into temporary slot file
	bin, err := json.Marshal(b)
	if err != nil {
		return err
	}
	tmpSlotFile, err := ioutil.TempFile(kv.path, SLOT_TMP_FILENAME_PATTERN)
	if err != nil {
		return err
	}
	if _, err := tmpSlotFile.Write(bin); err != nil {
		return err
	}
	t.Lock.Lock()
	defer t.Lock.Unlock()
	// rename temporary log file to actual log file
	if err := os.Rename(tmpLogFile.Name(), path.Join(kv.path, LOG_FILENAME)); err != nil {
		return err
	}
	// rename temporary slot file to actual slot file
	slotFileName := path.Join(kv.path, SLOT_FILENAME)
	if err := os.Rename(tmpSlotFile.Name(), slotFileName); err != nil {
		return err
	}
	return nil
}

// write to log(only to OS buffer)
func (kv *SimpleKV) writeLog(op string, args ...string) {
	// convert to quoted strings
	var quoted []string
	for _, v := range args {
		quoted = append(quoted, strconv.Quote(v))
	}
	l := fmt.Sprintf("%s %s\n", op, strings.Join(quoted, " "))
	_, err := kv.logFile.WriteString(l)
	if err != nil {
		common.Log().Error("Failed to write log",
			zap.String("op", op), zap.Strings("args", args), zap.Error(err))
	}
}

// flush log
func (kv *SimpleKV) Flush() {
	if err := kv.logFile.Sync(); err != nil {
		common.Log().Error("Failed to flush log.", zap.Error(err))
	}
}

func NewKVStore(pathString string) (*SimpleKV, error) {
	log := common.Log()
	// create the path if not exist
	if _, err := os.Stat(pathString); os.IsNotExist(err) {
		if err := os.Mkdir(pathString, 0755); err != nil {
			return nil, err
		}
	}

	// create a log file & slot file
	logFileName := path.Join(pathString, LOG_FILENAME)
	slotFileName := path.Join(pathString, SLOT_FILENAME)
	var logFile, slotFile *os.File

	var version uint64 = 0
	base := make(map[string]string)
	latest := make(map[string]*string)

	// open / create slot file
	if _, err := os.Stat(slotFileName); os.IsNotExist(err) {
		slotFile, err = os.Create(slotFileName)
		if err != nil {
			return nil, err
		}
		// initialize slot file
		if _, err := slotFile.Write([]byte("{}")); err != nil {
			return nil, err
		}
		if err := slotFile.Close(); err != nil {
			return nil, err
		}
		log.Info("Created new slot file.", zap.String("path", slotFileName))
	} else {
		// read slot file
		slotFile, err = os.Open(slotFileName)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(slotFile)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(b, &base); err != nil {
			return nil, err
		}
		log.Info("Recovered base from previous slot file.")
		if err := slotFile.Close(); err != nil {
			return nil, err
		}
	}
	// open / create log file
	if _, err := os.Stat(logFileName); os.IsNotExist(err) {
		logFile, err = os.Create(logFileName)
		if err != nil {
			return nil, err
		}
		log.Info("Created new log file.", zap.String("path", logFileName))
	} else {
		// read log file & parse it
		logFile, err = os.OpenFile(logFileName, os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}
		scanner := bufio.NewScanner(logFile)
		latest, _, err = ReadLog(scanner)
		if err != nil {
			return nil, err
		}
		log.Info("Recovered base from log entries.")
	}
	// transaction zero is always used and valid
	ts := make([]*TransactionStruct, TRANSACTION_COUNT)
	ts[0] = &TransactionStruct{
		Lock:  sync.RWMutex{},
		Layer: latest,
	}
	// others are nil
	return &SimpleKV{
		base:         base,
		layers:       make([]map[string]*string, 0),
		transactions: ts,
		path:         pathString,
		version:      version,
		logFile:      logFile,
	}, nil
}

func readString(quoted string) string {
	ret, err := strconv.Unquote(quoted)
	if err != nil {
		panic(err)
	}
	return ret
}

func readNum(quoted string, base int) uint64 {
	s := readString(quoted)
	ret, err := strconv.ParseUint(s, base, 64)
	if err != nil {
		panic(err)
	}
	return ret
}

func ReadLog(scanner *bufio.Scanner) (map[string]*string, uint64, error) {
	log := common.Log()
	trans := make([]map[string]*string, TRANSACTION_COUNT)
	trans[0] = make(map[string]*string)
	var tokens []string
	var version uint64
	defer func() {
		if x := recover(); x != nil {
			log.Error("Failed to parse log.", zap.Any("error", x), zap.Strings("tokens", tokens))
		}
	}()
	for scanner.Scan() {
		tokens = strings.Fields(scanner.Text())
		if len(tokens) == 0 {
			continue
		}
		switch tokens[0] {
		case "put":
			if len(tokens) < 4 {
				panic(nil)
			}
			key := readString(tokens[1])
			value := readString(tokens[2])
			transNum := int(readNum(tokens[3], 16))
			if trans[transNum] == nil {
				panic(nil)
			}
			if transNum == 0 {
				// get version number
				if len(tokens) < 5 {
					panic(nil)
				}
				version = readNum(tokens[4], 16)
			}
			trans[transNum][key] = &value
		case "del":
			if len(tokens) < 3 {
				panic(nil)
			}
			key := readString(tokens[1])
			transNum := int(readNum(tokens[2], 16))
			if trans[transNum] == nil {
				panic(nil)
			}
			if transNum == 0 {
				if len(tokens) < 4 {
					panic(nil)
				}
				version = readNum(tokens[3], 16)
			}
			trans[transNum][key] = nil
		case "start":
			// start transaction
			if len(tokens) < 2 {
				panic(nil)
			}
			transNum := int(readNum(tokens[1], 16))
			if trans[transNum] != nil {
				panic(nil)
			}
			trans[transNum] = make(map[string]*string)
		case "commit":
			// commit transaction
			if len(tokens) != 3 {
				panic(nil)
			}
			transNum := int(readNum(tokens[1], 16))
			if trans[transNum] == nil {
				panic(nil)
			}
			// merge into trans[0]
			for k, v := range trans[transNum] {
				trans[0][k] = v
			}
			trans[transNum] = nil
			version = readNum(tokens[2], 16)
		case "rollback":
			if len(tokens) != 2 {
				panic(nil)
			}
			transNum := int(readNum(tokens[1], 16))
			if trans[transNum] == nil {
				panic(nil)
			}
			trans[transNum] = nil
		default:
			panic(nil)
		}
	}
	return trans[0], version, nil
}

func (kv *SimpleKV) PrepareExtract() {
	// save layer
	kv.saveLayer()
}

func (kv *SimpleKV) Extract(divider func(key string) bool) map[string]string {
	kv.saveLayer()
	// extract content out
	b := make(map[string]string)
	for k, v := range kv.base {
		if divider(k) {
			b[k] = v
		}
	}
	for _, l := range kv.layers {
		for k, v := range l {
			if divider(k) {
				if v == nil {
					delete(b, k)
				} else {
					b[k] = *v
				}
			}
		}
	}
	return b
}

func (kv *SimpleKV) Close() {
	log := common.Log()
	kv.Flush()
	if err := kv.logFile.Close(); err != nil {
		log.Panic("Failed to close.", zap.Error(err))
	}
}
