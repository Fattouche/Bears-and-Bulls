package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	pb "github.com/Fattouche/DayTrader/golang/protobuff"
	"google.golang.org/grpc"
)

var symbolCommands = map[string]int{
	"QUOTE":            1,
	"BUY":              1,
	"SELL":             1,
	"SET_BUY_AMOUNT":   1,
	"CANCEL_SET_BUY":   1,
	"SET_BUY_TRIGGER":  1,
	"SET_SELL_AMOUNT":  1,
	"SET_SELL_TRIGGER": 1,
	"CANCEL_SET_SELL":  1,
}
var amountCommands = map[string]int{
	"ADD":              1,
	"BUY":              1,
	"SELL":             1,
	"SET_BUY_AMOUNT":   1,
	"SET_BUY_TRIGGER":  1,
	"SET_SELL_AMOUNT":  1,
	"SET_SELL_TRIGGER": 1,
}

var baseURL string
var userMap = make(map[string][]*pb.Command)
var dumpData *pb.Command
var wg sync.WaitGroup

func completeCall(command *pb.Command, client pb.DayTraderClient) {
	ctx := context.Background()
	var err error
	var resp *pb.Response
	switch command.Name {
	case "ADD":
		resp, err = client.Add(ctx, command)
	case "QUOTE":
		resp, err = client.Quote(ctx, command)
	case "BUY":
		resp, err = client.Buy(ctx, command)
	case "SELL":
		resp, err = client.Sell(ctx, command)
	case "COMMIT_BUY":
		resp, err = client.CommitBuy(ctx, command)
	case "COMMIT_SELL":
		resp, err = client.CommitSell(ctx, command)
	case "CANCEL_BUY":
		resp, err = client.CancelBuy(ctx, command)
	case "CANCEL_SELL":
		resp, err = client.CancelSell(ctx, command)
	case "SET_BUY_AMOUNT":
		resp, err = client.SetBuyAmount(ctx, command)
	case "SET_SELL_AMOUNT":
		resp, err = client.SetSellAmount(ctx, command)
	case "SET_BUY_TRIGGER":
		resp, err = client.SetBuyTrigger(ctx, command)
	case "SET_SELL_TRIGGER":
		resp, err = client.SetSellTrigger(ctx, command)
	case "CANCEL_SET_BUY":
		resp, err = client.CancelSetBuy(ctx, command)
	case "CANCEL_SET_SELL":
		resp, err = client.CancelSetSell(ctx, command)
	case "DUMPLOG":
		resp, err = client.DumpLog(ctx, command)
	case "DISPLAY_SUMMARY":
		resp, err = client.DisplaySummary(ctx, command)
	}
	_, _ = resp, err
	// UNCOMMENT FOR RESPONSES
	// fmt.Print(command.Name + " ")
	// log.Print("RESP: " + resp.String())
	// if err != nil {
	// 	fmt.Println(" ERROR: " + err.Error())
	// }
}

func parseCommands(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var userID string
	transactionNum := int32(1)
	for scanner.Scan() {
		totalCommand := strings.Split(scanner.Text(), " ")
		userCommands := strings.Split(totalCommand[1], ",")
		if userCommands[0] == "DUMPLOG" && len(userCommands) == 2 {
			dumpData = &pb.Command{Filename: userCommands[1], Name: "DUMPLOG"}
			continue
		} else {
			userID = userCommands[1]
		}
		command := generateRequest(userID, userCommands, transactionNum)
		transactionNum++
		if _, ok := userMap[userID]; !ok {
			userMap[userID] = make([]*pb.Command, 0)
		}
		userMap[userID] = append(userMap[userID], command)
	}
}

func generateRequest(userID string, commands []string, transactionNum int32) *pb.Command {
	command := &pb.Command{}
	if len(userID) > 0 {
		command.UserId = userID
	}
	if commands[0] == "DUMPLOG" {
		if len(commands) > 2 {
			command.Filename = commands[2]
		} else {
			command.Filename = commands[1]
		}
	}
	if _, ok := symbolCommands[commands[0]]; ok {
		if len(commands) > 2 {
			command.Symbol = commands[2]
		}
	}
	if _, ok := amountCommands[commands[0]]; ok {
		if commands[0] == "ADD" {
			if len(commands) > 2 {
				amount, _ := strconv.ParseFloat(commands[2], 32)
				command.Amount = float32(amount)
			}
		} else {
			if len(commands) > 3 {
				amount, _ := strconv.ParseFloat(commands[3], 32)
				command.Amount = float32(amount)
			}
		}
	}
	command.TransactionId = transactionNum
	command.Name = commands[0]
	return command
}

func makeRequest(commands []*pb.Command) {
	conn, err := grpc.Dial(baseURL, grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to dial to %s with %v", baseURL, err)
	}
	defer conn.Close()
	client := pb.NewDayTraderClient(conn)
	for _, command := range commands {
		completeCall(command, client)
	}
	wg.Done()
}

func makeDumpRequest() {
	conn, err := grpc.Dial(baseURL, grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to dial to %s with %v", baseURL, err)
	}
	client := pb.NewDayTraderClient(conn)
	completeCall(dumpData, client)
}

func main() {
	fileName := flag.String("f", "./workload_files/45_user_workload.workload", "The name of the workload file")
	tempBaseURL := flag.String("url", "daytrader_lb:80", "The url of the web server")
	flag.Parse()
	baseURL = *tempBaseURL
	parseCommands(*fileName)
	start := time.Now()
	wg.Add(len(userMap))
	for _, requests := range userMap {
		go makeRequest(requests)
	}
	wg.Wait()
	makeDumpRequest()
	fmt.Println("Time taken: ", time.Since(start))
}
