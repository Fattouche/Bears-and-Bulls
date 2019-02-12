package main

import (
	"context"

	pb "github.com/Fattouche/DayTrader/golang/protobuff"
	"google.golang.org/grpc"
)

func withServerUnaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(serverInterceptor)
}

func serverInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	if log, ok := req.(*pb.Log); ok {
		if err := checks(log); err != nil {
			return nil, err
		}
	}
	// Calls the handler
	h, err := handler(ctx, req)
	return h, err
}

// make sure user exists
func checks(req *pb.Log) error {
	//Do checks in here etc
	return nil
}
