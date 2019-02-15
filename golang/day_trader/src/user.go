package main

import (
	"context"
	"log"
)

func (user *User) toString() string {
	if user == nil {
		return ""
	}
	bytes, _ := user.MarshalJSON()
	return string(bytes)
}

func (user *User) popFromBuyStack() *Buy {
	if len(user.BuyStack) == 0 {
		return nil
	}
	buy := user.BuyStack[len(user.BuyStack)-1]
	user.BuyStack = user.BuyStack[:len(user.BuyStack)-1]
	user.setCache()
	return buy
}

func (user *User) popFromSellStack() *Sell {
	if len(user.SellStack) == 0 {
		return nil
	}
	sell := user.SellStack[len(user.SellStack)-1]
	user.SellStack = user.SellStack[:len(user.SellStack)-1]
	user.setCache()
	return sell
}

func getUser(userID string) *User {
	user, err := getCacheUser(userID)
	if err != nil {
		db.QueryRow("SELECT Balance from User where Id=?", user.Id).Scan(&user.Balance)
		user.SellStack = make([]*Sell, 0)
		user.BuyStack = make([]*Buy, 0)
		user.setCache()
	}
	return user
}

func createUser(userID string) error {
	user := &User{Id: userID, Balance: 0, Name: ""}
	user.SellStack = make([]*Sell, 0)
	user.BuyStack = make([]*Buy, 0)
	user.setCache()
	_, err := db.Exec("insert into User(Id) values(?)", userID)
	return err
}

func (user *User) updateUserBalance(ctx context.Context, amount float32) (*User, error) {
	var accountAction string
	if amount < 0 {
		accountAction = "Remove"
	} else {
		accountAction = "Add"
	}
	user.Balance += amount
	_, err := db.Exec("update User set Balance=? where Id=?", user.Balance, user.Id)
	if err != nil {
		return user, err
	}
	pbLog, err := makeLogFromContext(ctx)
	if err != nil {
		log.Println(err)
	}
	pbLog.AccountAction = accountAction
	pbLog.Funds = user.Balance
	logEvent := &logObj{log: &pbLog, funcName: "LogAccountTransaction"}
	logChan <- logEvent
	user.setCache()
	return user, nil
}
