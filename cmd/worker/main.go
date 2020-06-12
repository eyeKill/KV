// Worker for the distributed KV store
// Worker is data node. It stores the actual KV hashmap and responses to
// clients' requests
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/eyeKill/KV/common"
	pb "github.com/eyeKill/KV/proto"
	"github.com/eyeKill/KV/worker"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/samuel/go-zookeeper/zk"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	hostname  = flag.String("hostname", "localhost", "The server's hostname")
	port      = flag.Int("port", 7900, "The server port")
	path      = flag.String("path", "data/", "Path for persistent log and slot file.")
	zkServers = strings.Fields(*flag.String("zk-servers", "localhost:2181",
		"Zookeeper server cluster, separated by space"))
	zkNodeRoot = "/kv/nodes"
	zkNodeName = "worker"
)

var (
	conn   *zk.Conn
	kv     *worker.KVStore
	server *grpc.Server
	log    *zap.Logger
)

type WorkerServer struct {
	pb.UnimplementedKVWorkerServer
}

type WorkerInternalServer struct {
	pb.UnimplementedKVWorkerInternalServer
}

func (s *WorkerServer) Put(_ context.Context, pair *pb.KVPair) (*pb.PutResponse, error) {
	kv.Put(pair.Key, pair.Value)
	return &pb.PutResponse{Status: pb.Status_OK}, nil
}

func (s *WorkerServer) Get(_ context.Context, key *pb.Key) (*pb.GetResponse, error) {
	value, ok := kv.Get(key.Key)
	if ok {
		return &pb.GetResponse{
			Status: pb.Status_OK,
			Value:  value,
		}, nil
	} else {
		return &pb.GetResponse{
			Status: pb.Status_ENOENT,
			Value:  "",
		}, nil
	}
}

func (s *WorkerServer) Delete(_ context.Context, key *pb.Key) (*pb.DeleteResponse, error) {
	ok := kv.Delete(key.Key)
	if ok {
		return &pb.DeleteResponse{Status: pb.Status_OK}, nil
	} else {
		return &pb.DeleteResponse{Status: pb.Status_ENOENT}, nil
	}
}

func (s *WorkerInternalServer) Flush(_ context.Context, _ *empty.Empty) (*pb.FlushResponse, error) {
	if err := kv.Flush(); err != nil {
		log.Error("KV flush failed.", zap.Error(err))
		return &pb.FlushResponse{Status: pb.Status_EFAILED}, nil
	}
	return &pb.FlushResponse{Status: pb.Status_OK}, nil
}

func registerToZk(conn *zk.Conn) error {
	// don't have to ensure that the path exist here
	// since we're merely a worker
	nodePath := zkNodeRoot + "/" + zkNodeName
	exists, _, err := conn.Exists(zkNodeRoot)
	if err != nil {
		log.Panic("Failed to check whether root node exists.", zap.Error(err))
	} else if !exists {
		log.Panic("Root node does not exist.", zap.Error(err))
	}
	data := common.GetWorkerNodeData(*hostname, *port)
	b, err := json.Marshal(data)
	if err != nil {
		log.Panic("Failed to marshall into json object.", zap.Error(err))
	}
	name, err := conn.CreateProtectedEphemeralSequential(nodePath, b, zk.WorldACL(zk.PermAll))
	if err != nil {
		log.Panic("Failed to register itself to zookeeper.", zap.Error(err))
	}
	log.Info("Registration complete.", zap.String("name", name))
	return nil
}

func getGrpcServer() *grpc.Server {
	if server == nil {
		var opts []grpc.ServerOption
		server = grpc.NewServer(opts...)
	}
	return server
}

func runGrpcServer(server *grpc.Server) error {
	address := fmt.Sprintf("localhost:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	log.Info("Starting gRPC server...", zap.Int("port", *port))
	go func() {
		if err := server.Serve(listener); err != nil {
			log.Info("Error from gRPC server.", zap.Error(err))
		}
	}()
	return nil
}

// handle ctrl-c gracefully
func setupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Info("Ctrl-C captured.")
		if server != nil {
			log.Info("Gracefully stopping gRPC server...")
			server.GracefulStop()
		}
		if conn != nil {
			log.Info("Closing zookeeper connection...")
			conn.Close()
		}
		os.Exit(1)
	}()
}

func main() {
	setupCloseHandler()

	log = common.Log()

	flag.Parse()
	//if len(*hostname) == 0 {
	//	n, err := os.Hostname()
	//	if err != nil {
	//		log.Fatalf("Cannot get default hostname. Try to specify it in command line.")
	//	}
	//	hostname = &n
	//}
	// by default we bind to an arbitrary port
	// this behavior could be changed under environment like docker

	// initialize kv store
	kvStore, err := worker.NewKVStore(*path)
	if err != nil {
		log.Panic("Failed to create KVStore object.", zap.Error(err))
	}
	kv = kvStore

	// connect to zookeeper & register itself
	c, err := common.ConnectToZk(zkServers)
	if err != nil {
		log.Panic("Failed to connect too zookeeper cluster.", zap.Error(err))
	}
	log.Info("Connected to zookeeper cluster.", zap.String("server", c.Server()))
	conn = c // transfer to global scope

	defer conn.Close()
	if err := registerToZk(conn); err != nil {
		log.Panic("Failed to register to zookeeper cluster.", zap.Error(err))
	}

	// setup gRPC server & run it
	s := getGrpcServer()
	pb.RegisterKVWorkerServer(s, &WorkerServer{})
	pb.RegisterKVWorkerInternalServer(s, &WorkerInternalServer{})
	if err := runGrpcServer(s); err != nil {
		log.Panic("Failed to run gRPC server.", zap.Error(err))
	}

	// May you rest in a deep and restless slumber
	for {
		time.Sleep(10 * time.Second)
	}
}