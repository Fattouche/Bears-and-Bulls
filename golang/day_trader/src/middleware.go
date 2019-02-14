package main

import (
	"context"
	"errors"
	"log"

	pb "github.com/Fattouche/DayTrader/golang/protobuff"
	"google.golang.org/grpc"
)

type logKey string

func withServerUnaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(serverInterceptor)
}

func serverInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	if command, ok := req.(*pb.Command); ok {
		ctx = populateContextWithLogInfo(ctx, command)
		err := logUserCommand(ctx, command)
		if err != nil {
			log.Println("Error logging: ", err)
			return nil, err
		}
		if err := checks(command); err != nil {
			log.Println("Error creating user: ", err)
			return nil, err
		}
	}
	// Calls the handler
	h, err := handler(ctx, req)
	return h, err
}

// Populates the context with logging info, so that it can be extracted whenever
// a log event occurs
func populateContextWithLogInfo(ctx context.Context, command *pb.Command) context.Context {
	var key logKey

	// TODO: get the server name from an environment variable
	rawLog := pb.Log{
		TransactionNum: command.TransactionId,
		Username:       command.UserId,
		ServerName:     "Beaver_1",
		Command:        command.Name,
	}
	key = "log"
	ctx = context.WithValue(ctx, key, rawLog)
	return ctx
}

func logUserCommand(ctx context.Context, command *pb.Command) error {
	log, err := makeLogFromContext(ctx)
	if err != nil {
		return err
	}
	log.StockSymbol = command.Symbol
	log.Filename = command.Filename
	log.Funds = command.Amount
	logEvent := &logObj{log: &log, funcName: "LogUserCommand"}
	logChan <- logEvent
	return nil
}

// make sure user exists
func checks(req *pb.Command) error {
	createUser(req.UserId)
	if req.Name != "DUMPLOG" && req.Name != "DISPLAY_SUMMARY" {
		if req.UserId == "" {
			return errors.New("No user Id specified")
		}
	}
	return nil
}
