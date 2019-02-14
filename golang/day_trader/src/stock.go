package main

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func (stock *Stock) toString() string {
	if stock == nil {
		return ""
	}
	bytes, _ := stock.MarshalJSON()
	return string(bytes)
}

func quote(ctx context.Context, userID, symbol string) (*Stock, error) {
	stock, err := getCacheStock(symbol)
	if err != nil || (stock != nil && stock.isExpired()) {
		var quoteServerTimestamp int64
		//TODO: Do something with this hash
		stock = &Stock{Symbol: symbol}
		stock.Price, quoteServerTimestamp, stock.Hash, err = executeRequest(ctx, userID, symbol)
		err = logQuoteServerEvent(ctx, stock.Price, userID, symbol, stock.Hash, quoteServerTimestamp)
		stock.TimeStamp = time.Now()
		stock.setCache()
		return stock, err
	}
	return stock, err
}

func logQuoteServerEvent(ctx context.Context, price float32, userID string,
	symbol string, hash string, quoteServerTimestamp int64) error {
	pbLog, err := makeLogFromContext(ctx)
	if err != nil {
		return err
	}
	pbLog.Price = price
	pbLog.StockSymbol = symbol
	pbLog.Username = userID
	pbLog.QuoteServerTime = quoteServerTimestamp
	pbLog.CryptoKey = hash
	logEvent := &logObj{log: &pbLog, funcName: "logQuoteServerEvent"}
	logChan <- logEvent
	return nil
}

func (stock *Stock) isExpired() bool {
	duration := time.Since(stock.TimeStamp)
	if duration > time.Second*60 {
		return true
	}
	return false
}

func executeRequest(ctx context.Context, userID, symbol string) (float32, int64, string, error) {
	ln, err := net.Dial("tcp", QUOTE_HOST+QUOTE_PORT)
	defer ln.Close()
	if err != nil {
		return 0, 0, "", err
	}
	buf := make([]byte, 300)
	str := fmt.Sprintf("%s,%s\r", symbol, userID)
	ln.Write([]byte(str))
	len, err := ln.Read(buf)
	if err != nil {
		return 0, 0, "", err
	}
	info := string(buf[:len])
	infoArr := strings.Split(info, ",")
	price, err := strconv.ParseFloat(infoArr[0], 32)
	if err != nil {
		return 0, 0, "", err
	}
	quoteServerTimestamp, err := strconv.ParseInt(infoArr[3], 10, 64)
	if err != nil {
		return 0, 0, "", err
	}
	return float32(price), quoteServerTimestamp, infoArr[4], nil
}
