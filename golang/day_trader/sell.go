package main

import (
	"fmt"
	"log"
	"math"
	"time"
)

func createSell(intendedCashAmount float32, symbol, userID string) (*Sell, error) {
	stock, err := quote(userID, symbol)
	if err != nil {
		return nil, err
	}
	sell := &Sell{Price: stock.Price, StockSymbol: symbol, UserId: userID}
	err = sell.updateCashAmount(intendedCashAmount)
	if err != nil {
		return nil, err
	}
	err = sell.updatePrice(stock.Price)
	if err != nil {
		return nil, err
	}
	return sell, err
}

func (sell *Sell) updateCashAmount(amount float32) error {
	stock, _ := quote(sell.UserId, sell.StockSymbol)
	userStock := getOrCreateUserStock(sell.UserId, sell.StockSymbol)
	stockSoldAmount := int(math.Floor(float64(amount / stock.Price)))
	if stockSoldAmount > userStock.Amount {
		return fmt.Errorf("Not enough stock, have %d need %d", userStock.Amount, stockSoldAmount)
	}
	sell.IntendedCashAmount = amount
	return nil
}

func (sell *Sell) updatePrice(stockPrice float32) error {
	userStock := getOrCreateUserStock(sell.UserId, sell.StockSymbol)
	updateSoldAmount := int(math.Min(math.Floor(float64(sell.IntendedCashAmount/stockPrice)), float64(userStock.Amount+sell.StockSoldAmount)))
	updated := updateSoldAmount - sell.StockSoldAmount
	sell.StockSoldAmount += updated
	sell.ActualCashAmount = float32(sell.StockSoldAmount) * stockPrice
	sell.Timestamp = time.Now()
	sell.Price = stockPrice
	userStock.updateStockAmount(updated * -1)
	return nil
}

func (sell *Sell) commit(update bool) (err error) {
	user := getUser(sell.UserId)
	user.updateUserBalance(sell.ActualCashAmount)
	if update {
		err = sell.updateSell()
	} else {
		_, err = sell.insertSell()
	}
	return
}

func (sell *Sell) cancel() {
	userStock := getOrCreateUserStock(sell.UserId, sell.StockSymbol)
	userStock.updateStockAmount(sell.StockSoldAmount)
}

func (sell *Sell) updateSell() error {
	_, err := db.Exec("update Sell set IntendedCashAmount=?, Price=?, ActualCashAmount=?, StockSoldAmount = ? where Id=?", sell.IntendedCashAmount, sell.Price, sell.ActualCashAmount, sell.StockSoldAmount, sell.Id)
	if err != nil {
		return err
	}
	return err
}

func (sell *Sell) insertSell() (*Sell, error) {
	res, err := db.Exec("insert into Sell(Price,StockSymbol,UserId,IntendedCashAmount,ActualCashAmount,StockSoldAmount) values(?,?,?,?,?,?)", sell.Price, sell.StockSymbol, sell.UserId, sell.IntendedCashAmount, sell.ActualCashAmount, sell.StockSoldAmount)
	if err != nil {
		return sell, err
	}
	sell.Id, err = res.LastInsertId()
	return sell, err
}

func getSell(id int64) *Sell {
	sell := &Sell{}
	err := db.QueryRow("Select * from Sell where Id=?", id).Scan(&sell.Id, &sell.Price, &sell.StockSymbol, &sell.UserId, &sell.IntendedCashAmount, &sell.ActualCashAmount, &sell.StockSoldAmount)
	if err != nil {
		log.Println(err)
	}
	return sell
}
