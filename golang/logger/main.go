package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/Fattouche/DayTrader/golang/protobuff"
)

type server struct{}

func (s *server) LogUserCommand(ctx context.Context, req *pb.Log) (*pb.Response, error) {
	return &pb.Response{Message: "YEE"}, nil
}

func (s *server) LogQuoteServerEvent(ctx context.Context, req *pb.Log) (*pb.Response, error) {
	return &pb.Response{Message: "YEE"}, nil
}

func (s *server) LogAccountTransaction(ctx context.Context, req *pb.Log) (*pb.Response, error) {
	return &pb.Response{Message: "YEE"}, nil
}

func (s *server) LogSystemEvent(ctx context.Context, req *pb.Log) (*pb.Response, error) {
	return &pb.Response{Message: "YEE"}, nil
}

func (s *server) LogErrorEvent(ctx context.Context, req *pb.Log) (*pb.Response, error) {
	return &pb.Response{Message: "YEE"}, nil
}

func (s *server) LogDebugEvent(ctx context.Context, req *pb.Log) (*pb.Response, error) {
	return &pb.Response{Message: "YEE"}, nil
}

// Starts a generic GRPC server
func startGRPCServer() {
	lis, err := net.Listen("tcp", GRPC_PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(withServerUnaryInterceptor())
	pb.RegisterLoggerServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	createAndOpenDB()
	startGRPCServer()
}
