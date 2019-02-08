package main

import (
	"errors"
	"log"
)

type SellTrigger struct {
	UserId string
	SellId int64
	Active bool
}

func (trigger *SellTrigger) updateCashAmount(amount float32) error {
	sell := getSell(trigger.SellId)
	err := sell.updateCashAmount(amount)
	if err != nil {
		return err
	}
	sell.updateSell()
	return err
}

func (trigger *SellTrigger) updatePrice(price float32) error {
	sell := getSell(trigger.SellId)
	err := sell.updatePrice(price)
	if err == nil {
		trigger.Active = true
		db.Exec("UPDATE Sell_Trigger set Active=true where UserId=? and SellId=?", trigger.UserId, trigger.SellId)
	}
	return err
}

func (trigger *SellTrigger) cancel() {
	sell := getSell(trigger.SellId)
	sell.cancel()
	trigger.Active = false
	db.Exec("UPDATE Sell_Trigger set Active=false where UserId=? and SellId=?", trigger.UserId, trigger.SellId)
}

func getSellTrigger(userID, symbol string) (*SellTrigger, error) {
	sellTrigger := &SellTrigger{UserId: userID, SellId: -1}
	db.QueryRow("SELECT Sell.Id, Sell_Trigger.Active from Sell inner join Sell_Trigger on Sell_Trigger.SellId=Sell.Id where Sell_Trigger.UserId=? and Sell.StockSymbol=?", sellTrigger.UserId, symbol).Scan(&sellTrigger.SellId, &sellTrigger.Active)
	if sellTrigger.SellId == -1 {
		return nil, errors.New("No sell trigger found")
	}
	return sellTrigger, nil
}

func createSellTrigger(userID, symbol string, sellID int64, amount float32) *SellTrigger {
	_, err := db.Exec("insert into Sell_Trigger(UserId,SellId) values(?,?)", userID, sellID)
	if err != nil {
		log.Println(err)
	}
	sellTrigger := &SellTrigger{UserId: userID, SellId: sellID, Active: false}
	return sellTrigger
}

func checkSellTriggers() {
	rows, err := db.Query("SELECT Sell.Id, Sell.StockSymbol, Sell.UserId from Sell inner join Sell_Trigger on Sell_Trigger.SellId=Sell.Id where Sell_Trigger.Active=true")
	if err != nil {
		log.Println(err)
	}
	sells := make([]*Sell, 0)
	for rows.Next() {
		sell := &Sell{}
		err = rows.Scan(&sell.Id, &sell.StockSymbol, sell.UserId)
		if err != nil {
			log.Println("Error scanning trigger: ", err)
		}
		sells = append(sells, sell)
	}
	rows.Close()
	//TODO: Improve this preformance
	for _, sell := range sells {
		stock, _ := quote(sell.UserId, sell.StockSymbol)
		if sell.Price >= stock.Price {
			sell.updatePrice(stock.Price)
			sell.commit()
			db.Exec("Update Sell_Trigger set Active=false where SellId=? and UserId=?", sell.Id, sell.UserId)
		}
	}
}
