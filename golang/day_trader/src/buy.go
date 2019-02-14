package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"time"
)

func (buy *Buy) toString() string {
	if buy == nil {
		return ""
	}
	bytes, _ := buy.MarshalJSON()
	return string(bytes)
}

func createBuy(ctx context.Context, intendedCashAmount float32, symbol, userID string) (*Buy, error) {
	stock, err := quote(ctx, userID, symbol)
	if err != nil {
		return nil, err
	}
	buy := &Buy{Price: stock.Price, StockSymbol: symbol, UserId: userID}
	user, err := buy.updateCashAmount(ctx, intendedCashAmount)
	if err != nil {
		return nil, err
	}
	buy.updatePrice(stock.Price)
	buy.Timestamp = time.Now()
	user.BuyStack = append(user.BuyStack, buy)
	user.setCache()
	return buy, err
}

func (buy *Buy) updateCashAmount(ctx context.Context, amount float32) (*User, error) {
	user := getUser(buy.UserId)
	if amount > user.Balance {
		msg := fmt.Sprintf("Not enough balance, have %f need %f", user.Balance, amount)
		return nil, errors.New(msg)
	}
	updatedAmount := buy.IntendedCashAmount - amount
	user.updateUserBalance(ctx, updatedAmount)
	buy.IntendedCashAmount = float32(math.Abs(float64(updatedAmount)))
	return user, nil
}

func (buy *Buy) updatePrice(stockPrice float32) {
	buy.Price = stockPrice
	buy.StockBoughtAmount = int(math.Floor(float64(buy.IntendedCashAmount / buy.Price)))
	buy.ActualCashAmount = float32(buy.StockBoughtAmount) * buy.Price
}

func (buy *Buy) commit(ctx context.Context, update bool) (*UserStock, error) {
	var err error
	if update {
		err = buy.updateBuy(ctx)
	} else {
		_, err = buy.insertBuy(ctx)
	}
	userStock := getOrCreateUserStock(ctx, buy.UserId, buy.StockSymbol)
	userStock.updateStockAmount(ctx, buy.StockBoughtAmount)
	return userStock, err
}

func (buy *Buy) cancel(ctx context.Context) {
	user := getUser(buy.UserId)
	user.updateUserBalance(ctx, buy.IntendedCashAmount)
}

func (buy *Buy) updateBuy(ctx context.Context) error {
	_, err := db.Exec("update Buy set IntendedCashAmount=?, Price=?, ActualCashAmount=?, StockBoughtAmount = ? where Id=?", buy.IntendedCashAmount, buy.Price, buy.ActualCashAmount, buy.StockBoughtAmount, buy.Id)
	if err != nil {
		return err
	}
	return err
}

func (buy *Buy) insertBuy(ctx context.Context) (*Buy, error) {
	res, err := db.Exec("insert into Buy(Price,StockSymbol,UserId,IntendedCashAmount,ActualCashAmount,StockBoughtAmount) values(?,?,?,?,?,?)", buy.Price, buy.StockSymbol, buy.UserId, buy.IntendedCashAmount, buy.ActualCashAmount, buy.StockBoughtAmount)
	if err != nil {
		return buy, err
	}
	buy.Id, err = res.LastInsertId()
	return buy, err
}

func getBuy(ctx context.Context, id int64) *Buy {
	buy := &Buy{}
	err := db.QueryRow("Select * from Buy where Id=?", id).Scan(&buy.Id, &buy.Price, &buy.StockSymbol, &buy.UserId, &buy.IntendedCashAmount, &buy.ActualCashAmount, &buy.StockBoughtAmount)
	if err != nil {
		log.Println(err)
	}
	return buy
}
