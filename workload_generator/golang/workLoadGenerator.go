package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var getMap = map[string]int{"QUOTE": 1, "DUMPLOG": 1, "DISPLAY_SUMMARY": 1}
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
var userMap = make(map[string][]*http.Request)
var wg sync.WaitGroup
var baseURL string
var client = &http.Client{}

func parseCommands(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var userID string
	for scanner.Scan() {
		totalCommand := strings.Split(scanner.Text(), " ")
		userCommands := strings.Split(totalCommand[1], ",")
		if userCommands[0] == "DUMPLOG" && len(userCommands) == 2 {
			userID = "admin"
		} else {
			userID = userCommands[1]
		}
		req := generateRequest(userID, userCommands)
		if _, ok := userMap[userID]; !ok {
			userMap[userID] = make([]*http.Request, 0)
		}
		userMap[userID] = append(userMap[userID], req)
	}
}

func generateRequest(userID string, commands []string) *http.Request {
	requestType := "POST"
	if _, ok := getMap[commands[0]]; ok {
		requestType = "GET"
	}
	params := make(map[string]string)
	params["user_id"] = userID
	if commands[0] == "DUMPLOG" {
		if len(commands) > 2 {
			params["filename"] = commands[2]
		} else {
			params["filename"] = commands[1]
		}
	}
	if _, ok := symbolCommands[commands[0]]; ok {
		if len(commands) > 2 {
			params["symbol"] = commands[2]
		}
	}
	if _, ok := amountCommands[commands[0]]; ok {
		if commands[0] == "ADD" {
			if len(commands) > 2 {
				params["amount"] = commands[2]
			}
		} else {
			if len(commands) > 3 {
				params["amount"] = commands[3]
			}
		}
	}
	url := baseURL + strings.ToLower(commands[0])
	if requestType == "POST" {
		jsonValue, _ := json.Marshal(params)
		req, err := http.NewRequest(requestType, url, bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			log.Println("Error creating request: ", err)
			os.Exit(1)
		}
		return req
	}
	req, err := http.NewRequest(requestType, url, nil)
	if err != nil {
		log.Println("Error creating request: ", err)
		os.Exit(1)
	}
	q := req.URL.Query()
	for key, val := range params {
		q.Add(key, val)
	}
	req.URL.RawQuery = q.Encode()
	return req
}

func makeRequest(requests []*http.Request) {
	for _, req := range requests {
		resp, err := client.Do(req)
		if err != nil {
			log.Println("ERROR: ", err)
		}
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
	wg.Done()
}

func main() {
	fileName := flag.String("f", "1userWorkLoad", "The name of the workload file")
	tempBaseURL := flag.String("url", "http://day_trader_lb/", "The url of the web server")
	flag.Parse()
	baseURL = *tempBaseURL
	parseCommands(*fileName)
	wg.Add(len(userMap))
	start := time.Now()
	for _, requests := range userMap {
		go makeRequest(requests)
	}
	wg.Wait()
	fmt.Println("Time taken: ", time.Since(start))
}
