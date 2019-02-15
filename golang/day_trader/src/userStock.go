package main

import (
	"context"
)

func (userStock *UserStock) toString() string {
	if userStock == nil {
		return ""
	}
	bytes, _ := userStock.MarshalJSON()
	return string(bytes)
}

func getOrCreateUserStock(ctx context.Context, userID, symbol string, user *User) *UserStock {
	if amount, ok := user.StockMap[symbol]; ok {
		return &UserStock{UserId: userID, StockSymbol: symbol, Amount: amount}
	}
	amount := 0
	db.QueryRow("Select Amount from User_Stock where UserId=?", user.Id).Scan(&amount)
	userStock := &UserStock{UserId: userID, StockSymbol: symbol, Amount: amount}
	user.setCache()
	return userStock
}

func (userStock *UserStock) updateStockAmount(ctx context.Context, amount int, user *User) {
	userStock.Amount += amount
	user.StockMap[userStock.StockSymbol] = userStock.Amount
	user.setCache()
}
