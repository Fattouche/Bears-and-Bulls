package main

import (
	"database/sql"
	"time"
)

var createTableStatements = []string{
	`CREATE DATABASE IF NOT EXISTS ` + DB_NAME,
	`USE ` + DB_NAME,
	`CREATE TABLE IF NOT EXISTS User(
		Id varchar(32) NOT NULL,
		Balance float DEFAULT 0, 
		Name varchar(32) NULL,
		PRIMARY KEY (id)
	)`,
	`CREATE TABLE IF NOT EXISTS Sell(
		Id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		Price float DEFAULT 0,
		StockSymbol varchar(3) NULL,
		UserId varchar(32) NOT NULL,
		IntendedCashAmount float DEFAULT 0,
		ActualCashAmount float DEFAULT 0,
		StockSoldAmount int DEFAULT 0,
		PRIMARY KEY (Id),
		FOREIGN KEY (UserId) REFERENCES User(Id) ON DELETE CASCADE
	)`,
	`CREATE TABLE IF NOT EXISTS Buy(
		Id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		Price float DEFAULT 0,
		StockSymbol varchar(3) NULL,
		UserId varchar(32) NOT NULL,
		IntendedCashAmount float DEFAULT 0,
		ActualCashAmount float DEFAULT 0,
		StockBoughtAmount int DEFAULT 0,
		PRIMARY KEY (Id),
		FOREIGN KEY (UserId) REFERENCES User(Id) ON DELETE CASCADE
	)`,
	`CREATE TABLE IF NOT EXISTS Sell_Trigger(
		UserId varchar(32) NOT NULL,
		SellId INT UNSIGNED NOT NULL,
		Active BOOLEAN DEFAULT false,
		UNIQUE(UserId, SellId),
		FOREIGN KEY (UserId) REFERENCES User(Id) ON DELETE CASCADE,
		FOREIGN KEY (SellId) REFERENCES Sell(Id) ON DELETE CASCADE
	)`,
	`CREATE TABLE IF NOT EXISTS Buy_Trigger(
		UserId varchar(32) NOT NULL,
		BuyId INT UNSIGNED NOT NULL,
		Active BOOLEAN DEFAULT false,
		UNIQUE(UserId, BuyId),
		FOREIGN KEY (UserId) REFERENCES User(Id) ON DELETE CASCADE,
		FOREIGN KEY (BuyId) REFERENCES Buy(Id) ON DELETE CASCADE
	)`,
	`CREATE TABLE IF NOT EXISTS User_Stock(
		UserId varchar(32) NOT NULL,
		StockSymbol varchar(3) NULL,
		Amount INT DEFAULT 0,
		UNIQUE(UserId, StockSymbol),
		FOREIGN KEY (UserId) REFERENCES User(Id) ON DELETE CASCADE
	)`,
}
var db *sql.DB

func createAndOpenDB() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(daytrader_db:3306)/")
	if err != nil {
		panic(err)
	}

	for _, stmt := range createTableStatements {
		_, err := db.Exec(stmt)
		if err != nil {
			if err != nil {
				panic(err)
			}
		}
	}
	db, err = sql.Open("mysql", "root@tcp(daytrader_db:3306)/daytrader")
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Second * 0)
	db.SetMaxIdleConns(10000)
}
