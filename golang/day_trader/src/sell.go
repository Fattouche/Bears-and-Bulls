package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"
)

func (sell *Sell) toString() string {
	if sell == nil {
		return ""
	}
	bytes, _ := sell.MarshalJSON()
	return string(bytes)
}

func createSell(ctx context.Context, intendedCashAmount float32, symbol string, user *User) (*Sell, error) {
	stock, err := quote(ctx, user.Id, symbol)
	if err != nil {
		return nil, err
	}
	sell := &Sell{Price: stock.Price, StockSymbol: symbol, UserId: user.Id}
	err = sell.updateCashAmount(ctx, intendedCashAmount, user)
	if err != nil {
		return nil, err
	}
	err = sell.updatePrice(ctx, stock.Price, user)
	if err != nil {
		return nil, err
	}
	return sell, err
}

func (sell *Sell) updateCashAmount(ctx context.Context, amount float32, user *User) error {
	stock, _ := quote(ctx, sell.UserId, sell.StockSymbol)
	userStock := getOrCreateUserStock(ctx, sell.UserId, sell.StockSymbol, user)
	stockSoldAmount := int(math.Floor(float64(amount / stock.Price)))
	if stockSoldAmount > userStock.Amount {
		return fmt.Errorf("Not enough stock, have %d need %d", userStock.Amount, stockSoldAmount)
	}
	sell.IntendedCashAmount = amount
	return nil
}

func (sell *Sell) updatePrice(ctx context.Context, stockPrice float32, user *User) error {
	userStock := getOrCreateUserStock(ctx, sell.UserId, sell.StockSymbol, user)
	updateSoldAmount := int(math.Min(math.Floor(float64(sell.IntendedCashAmount/stockPrice)), float64(userStock.Amount+sell.StockSoldAmount)))
	updated := updateSoldAmount - sell.StockSoldAmount
	sell.StockSoldAmount += updated
	sell.ActualCashAmount = float32(sell.StockSoldAmount) * stockPrice
	sell.Timestamp = time.Now()
	sell.Price = stockPrice
	userStock.updateStockAmount(ctx, updated*-1, user)
	return nil
}

func (sell *Sell) commit(ctx context.Context, update bool, user *User) (err error) {
	err = user.updateStockBalance(ctx, sell.StockSymbol)
	if err != nil {
		return err
	}
	user.updateUserBalance(ctx, sell.ActualCashAmount, true)
	if update {
		err = sell.updateSell(ctx)
	} else {
		_, err = sell.insertSell(ctx)
	}
	return
}

func (sell *Sell) cancel(ctx context.Context, user *User) {
	userStock := getOrCreateUserStock(ctx, sell.UserId, sell.StockSymbol, user)
	userStock.updateStockAmount(ctx, sell.StockSoldAmount, user)
}

func (sell *Sell) updateSell(ctx context.Context) error {
	_, err := db.Exec("update Sell set IntendedCashAmount=?, Price=?, ActualCashAmount=?, StockSoldAmount = ? where Id=?", sell.IntendedCashAmount, sell.Price, sell.ActualCashAmount, sell.StockSoldAmount, sell.Id)
	if err != nil {
		return err
	}
	return err
}

func (sell *Sell) insertSell(ctx context.Context) (*Sell, error) {
	res, err := db.Exec("insert into Sell(Price,StockSymbol,UserId,IntendedCashAmount,ActualCashAmount,StockSoldAmount) values(?,?,?,?,?,?)", sell.Price, sell.StockSymbol, sell.UserId, sell.IntendedCashAmount, sell.ActualCashAmount, sell.StockSoldAmount)
	if err != nil {
		return sell, err
	}
	sell.Id, err = res.LastInsertId()
	return sell, err
}

func getSell(ctx context.Context, id int64) *Sell {
	sell := &Sell{}
	err := db.QueryRow("Select * from Sell where Id=?", id).Scan(&sell.Id, &sell.Price, &sell.StockSymbol, &sell.UserId, &sell.IntendedCashAmount, &sell.ActualCashAmount, &sell.StockSoldAmount)
	if err != nil {
		log.Println(err)
	}
	return sell
}
