package main

import (
	"errors"
	"log"
)

type BuyTrigger struct {
	UserId string
	BuyId  int64
	Active bool
}

func getBuyTrigger(userID, symbol string) (*BuyTrigger, error) {
	buyTrigger := &BuyTrigger{UserId: userID, BuyId: -1}
	db.QueryRow("SELECT Buy.Id,Buy_Trigger.Active from Buy_Trigger inner join Buy on Buy_Trigger.BuyId=Buy.Id where Buy_Trigger.UserId=? and Buy.StockSymbol=?", buyTrigger.UserId, symbol).Scan(&buyTrigger.BuyId, &buyTrigger.Active)
	if buyTrigger.BuyId == -1 {
		return nil, errors.New("No Buy trigger found")
	}
	return buyTrigger, nil
}

func (trigger *BuyTrigger) updateCashAmount(amount float32) error {
	buy := getBuy(trigger.BuyId)
	err := buy.updateCashAmount(amount)
	if err != nil {
		return err
	}
	buy.updateBuy()
	return err
}

func (trigger *BuyTrigger) updatePrice(price float32) {
	buy := getBuy(trigger.BuyId)
	buy.updatePrice(price)
	buy.updateBuy()
	db.Exec("UPDATE Buy_Trigger set Active=true where UserId=? and BuyId=?", trigger.UserId, trigger.BuyId)
}

func (trigger *BuyTrigger) cancel() {
	buy := getBuy(trigger.BuyId)
	buy.cancel()
	trigger.Active = false
	db.Exec("UPDATE Buy_Trigger set Active=false where UserId=? and SellId=?", trigger.UserId, trigger.BuyId)
}

func createBuyTrigger(userID, symbol string, buyID int64, amount float32) *BuyTrigger {
	_, err := db.Exec("insert into Buy_Trigger(UserId,BuyId) values(?,?)", userID, buyID)
	if err != nil {
		log.Println(err)
	}
	buyTrigger := &BuyTrigger{UserId: userID, BuyId: buyID, Active: false}
	return buyTrigger
}

func checkBuyTriggers() {
	rows, err := db.Query("SELECT Buy.Id, Buy.StockSymbol, Buy.UserId from Buy inner join Buy_Trigger on Buy_Trigger.BuyId=Buy.Id where Buy_Trigger.Active=true")
	if err != nil {
		log.Println(err)
	}
	buys := make([]*Buy, 0)
	for rows.Next() {
		buy := &Buy{}
		err = rows.Scan(&buy.Id, buy.UserId)
		if err != nil {
			log.Println("Error scanning trigger: ", err)
		}
		buys = append(buys, buy)
	}
	rows.Close()
	//TODO: Improve this preformance
	for _, buy := range buys {
		stock, _ := quote(buy.UserId, buy.StockSymbol)
		if buy.Price <= stock.Price {
			buy.updatePrice(stock.Price)
			buy.commit()
			db.Exec("Update Buy_Trigger set Active=false where BuyId=? and UserId=?", buy.Id, buy.UserId)
		}
	}
}
