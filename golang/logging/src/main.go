package main

import (
	"context"
	"log"
	"net"
	"strconv"

	pb "github.com/Fattouche/DayTrader/golang/protobuff"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// This server implements the protobuff Node type
type server struct{}

var (
	grpcPort = ":40000"
)

func LogBase(ctx context.Context, req *pb.Log) int {
	res, err := db.Exec("INSERT INTO BaseLog(LogType, Server, TransactionNum) VALUES(?,?,?)", "userCommand", req.ServerName, req.TransactionNum)
	if err != nil {
		log.Println(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
	}
	return int(id)
}

func (s *server) LogUserCommand(ctx context.Context, req *pb.Log) (*pb.Response, error) {
	baseLogID := LogBase(ctx, req)
	_, err := db.Exec(
		"INSERT INTO UserCommandLog VALUES(?,?,?,?,?,?)",
		baseLogID, req.Command, req.Username, req.StockSymbol, req.Filename,
		req.Funds,
	)
	if err != nil {
		log.Println(err)
	}
	return &pb.Response{Message: strconv.Itoa(baseLogID)}, err
}

func (s *server) LogQuoteServerEvent(ctx context.Context, req *pb.Log) (*pb.Response, error) {
	baseLogID := LogBase(ctx, req)
	_, err := db.Exec(
		"INSERT INTO QuoteServerLog VALUES(?,?,?,?,?,?)",
		baseLogID, req.Price, req.StockSymbol, req.Username,
		req.QuoteServerTime, req.CryptoKey,
	)
	if err != nil {
		log.Println(err)
	}
	return &pb.Response{Message: strconv.Itoa(baseLogID)}, err
}

func (s *server) LogAccountTransaction(ctx context.Context, req *pb.Log) (*pb.Response, error) {
	baseLogID := LogBase(ctx, req)
	_, err := db.Exec(
		"INSERT INTO AccountTransactionLog VALUES(?,?,?,?)",
		baseLogID, req.AccountAction, req.Username, req.Funds,
	)
	if err != nil {
		log.Println(err)
	}
	return &pb.Response{Message: strconv.Itoa(baseLogID)}, err
}

func (s *server) LogSystemEvent(ctx context.Context, req *pb.Log) (*pb.Response, error) {
	baseLogID := LogBase(ctx, req)
	_, err := db.Exec(
		"INSERT INTO SystemEventLog VALUES(?,?,?,?,?,?)",
		baseLogID, req.Command, req.Username, req.StockSymbol,
		req.Filename, req.Funds,
	)
	if err != nil {
		log.Println(err)
	}
	return &pb.Response{Message: strconv.Itoa(baseLogID)}, err
}

func (s *server) LogErrorEvent(ctx context.Context, req *pb.Log) (*pb.Response, error) {
	baseLogID := LogBase(ctx, req)
	_, err := db.Exec(
		"INSERT INTO SystemEventLog VALUES(?,?,?,?,?,?,?)",
		baseLogID, req.Command, req.Username, req.StockSymbol,
		req.Filename, req.Funds, req.ErrorMessage,
	)
	if err != nil {
		log.Println(err)
	}
	return &pb.Response{Message: strconv.Itoa(baseLogID)}, err
}

func (s *server) LogDebugEvent(ctx context.Context, req *pb.Log) (*pb.Response, error) {
	baseLogID := LogBase(ctx, req)
	_, err := db.Exec(
		"INSERT INTO SystemEventLog VALUES(?,?,?,?,?,?,?)",
		baseLogID, req.Command, req.Username, req.StockSymbol,
		req.Filename, req.Funds, req.DebugMessage,
	)
	if err != nil {
		log.Println(err)
	}
	return &pb.Response{Message: strconv.Itoa(baseLogID)}, err
}

func (s *server) DumpLogs(ctx context.Context, req *pb.Command) (*pb.Response, error) {
	return &pb.Response{Message: "yee dump"}, nil
}

func startGRPCServer() {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterLoggerServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	log.Println("Running logging server")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	createAndOpenDB()
	startGRPCServer()
}
