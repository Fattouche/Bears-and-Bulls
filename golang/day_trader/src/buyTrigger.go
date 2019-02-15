package main

import (
	"context"
	"errors"
	"log"
)

func (trigger *BuyTrigger) toString() string {
	if trigger == nil {
		return ""
	}
	bytes, _ := trigger.MarshalJSON()
	return string(bytes)
}

func getBuyTrigger(ctx context.Context, userID, symbol string) (*BuyTrigger, error) {
	buyTrigger := &BuyTrigger{UserId: userID, BuyId: -1}
	db.QueryRow("SELECT Buy.Id,Buy_Trigger.Active from Buy_Trigger inner join Buy on Buy_Trigger.BuyId=Buy.Id where Buy_Trigger.UserId=? and Buy.StockSymbol=?", buyTrigger.UserId, symbol).Scan(&buyTrigger.BuyId, &buyTrigger.Active)
	if buyTrigger.BuyId == -1 {
		return nil, errors.New("No Buy trigger found")
	}
	return buyTrigger, nil
}

func (trigger *BuyTrigger) updateCashAmount(ctx context.Context, amount float32, user *User) error {
	buy := getBuy(ctx, trigger.BuyId)
	err := buy.updateCashAmount(ctx, amount, user)
	if err != nil {
		return err
	}
	buy.updateBuy(ctx)
	return err
}

func (trigger *BuyTrigger) updatePrice(ctx context.Context, price float32) {
	buy := getBuy(ctx, trigger.BuyId)
	buy.updatePrice(price)
	buy.updateBuy(ctx)
	trigger.Active = true
	db.Exec("UPDATE Buy_Trigger set Active=true where UserId=? and BuyId=?", trigger.UserId, trigger.BuyId)
}

func (trigger *BuyTrigger) cancel(ctx context.Context, user *User) {
	buy := getBuy(ctx, trigger.BuyId)
	buy.cancel(ctx, user)
	db.Exec("DELETE From Buy_Trigger UserId=? and BuyId=?", trigger.UserId, trigger.BuyId)
	db.Exec("DELETE From Buy where Id=?", trigger.BuyId)
}

func createBuyTrigger(ctx context.Context, userID, symbol string, buyID int64, amount float32) *BuyTrigger {
	_, err := db.Exec("insert into Buy_Trigger(UserId,BuyId) values(?,?)", userID, buyID)
	if err != nil {
		log.Println(err)
	}
	buyTrigger := &BuyTrigger{UserId: userID, BuyId: buyID, Active: false}
	return buyTrigger
}

func checkBuyTriggers() {
	rows, err := db.Query("SELECT Buy.Id, Buy.StockSymbol, Buy.IntendedCashAmount,Buy.StockBoughtAmount, Buy.Price, Buy.UserId from Buy inner join Buy_Trigger on Buy_Trigger.BuyId=Buy.Id where Buy_Trigger.Active=true")
	if err != nil {
		log.Println(err)
	}
	buys := make([]*Buy, 0)
	for rows.Next() {
		buy := &Buy{}
		err = rows.Scan(&buy.Id, &buy.StockSymbol, &buy.IntendedCashAmount, &buy.StockBoughtAmount, &buy.Price, &buy.UserId)
		if err != nil {
			log.Println("Error scanning trigger: ", err)
		}
		buys = append(buys, buy)
	}
	rows.Close()
	for _, buy := range buys {
		stock, _ := quote(context.Background(), buy.UserId, buy.StockSymbol)
		if buy.Price >= stock.Price {
			buy.updatePrice(stock.Price)
			buy.commit(context.Background(), getUser(userId), true)
			db.Exec("Delete From Buy_Trigger where BuyId=? and UserId=?", buy.Id, buy.UserId)
		}
	}
}
