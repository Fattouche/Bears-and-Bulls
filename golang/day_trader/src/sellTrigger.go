package main

import (
	"context"
	"errors"
	"log"
)

func (trigger *SellTrigger) toString() string {
	if trigger == nil {
		return ""
	}
	bytes, _ := trigger.MarshalJSON()
	return string(bytes)
}

func (trigger *SellTrigger) updateCashAmount(ctx context.Context, amount float32, user *User) error {
	sell := getSell(ctx, trigger.SellId)
	err := sell.updateCashAmount(ctx, amount, user)
	if err != nil {
		return err
	}
	sell.updateSell(ctx)
	return err
}

func (trigger *SellTrigger) updatePrice(ctx context.Context, price float32, user *User) error {
	sell := getSell(ctx, trigger.SellId)
	err := sell.updatePrice(ctx, price, user)
	if err == nil {
		trigger.Active = true
		sell.updateSell(ctx)
		db.Exec("UPDATE Sell_Trigger set Active=true where UserId=? and SellId=?", trigger.UserId, trigger.SellId)
	}
	return err
}

func (trigger *SellTrigger) cancel(ctx context.Context, user *User) {
	sell := getSell(ctx, trigger.SellId)
	sell.cancel(ctx, user)
	db.Exec("DELETE from Sell_Trigger UserId=? and SellId=?", trigger.UserId, trigger.SellId)
	db.Exec("DELETE From Sell where Id=?", trigger.SellId)
}

func getSellTrigger(ctx context.Context, userID, symbol string) (*SellTrigger, error) {
	sellTrigger := &SellTrigger{UserId: userID, SellId: -1}
	db.QueryRow("SELECT Sell.Id, Sell_Trigger.Active from Sell inner join Sell_Trigger on Sell_Trigger.SellId=Sell.Id where Sell_Trigger.UserId=? and Sell.StockSymbol=?", sellTrigger.UserId, symbol).Scan(&sellTrigger.SellId, &sellTrigger.Active)
	if sellTrigger.SellId == -1 {
		return nil, errors.New("No sell trigger found")
	}
	return sellTrigger, nil
}

func createSellTrigger(ctx context.Context, userID, symbol string, sellID int64, amount float32) *SellTrigger {
	_, err := db.Exec("insert into Sell_Trigger(UserId,SellId) values(?,?)", userID, sellID)
	if err != nil {
		log.Println(err)
	}
	sellTrigger := &SellTrigger{UserId: userID, SellId: sellID, Active: false}
	return sellTrigger
}

func checkSellTriggers() {
	rows, err := db.Query("SELECT Sell.Id, Sell.StockSymbol, Sell.IntendedCashAmount,Sell.StockSoldAmount, Sell.Price,Sell.UserId from Sell inner join Sell_Trigger on Sell_Trigger.SellId=Sell.Id where Sell_Trigger.Active=true")
	if err != nil {
		log.Println(err)
	}
	sells := make([]*Sell, 0)
	for rows.Next() {
		sell := &Sell{}
		err = rows.Scan(&sell.Id, &sell.StockSymbol, &sell.IntendedCashAmount, &sell.StockSoldAmount, &sell.Price, &sell.UserId)
		if err != nil {
			log.Println("Error scanning trigger: ", err)
		}
		sells = append(sells, sell)
	}
	rows.Close()
	for _, sell := range sells {
		stock, _ := quote(context.Background(), sell.UserId, sell.StockSymbol)
		if sell.Price <= stock.Price {
			user := getUser(sell.UserId)
			sell.updatePrice(context.Background(), stock.Price, user)
			sell.commit(context.Background(), true, user)
			db.Exec("Delete From Sell_Trigger where SellId=? and UserId=?", sell.Id, sell.UserId)
		}
	}
}
