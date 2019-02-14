package main

import (
	"database/sql"
	"time"
)

var (
	dbName = "logging" // Name of database to store logs
	db     *sql.DB
)

var createTableStatements = []string{
	`CREATE DATABASE IF NOT EXISTS ` + dbName,
	`USE ` + dbName,
	`CREATE TABLE IF NOT EXISTS BaseLog(
		ID BIGINT NOT NULL AUTO_INCREMENT,
		LogType VARCHAR(32) NOT NULL,
		Timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, 
		Server VARCHAR(32) NOT NULL,
		TransactionNum BIGINT NOT NULL,
		PRIMARY KEY (ID)
	)`,
	`CREATE TABLE IF NOT EXISTS UserCommandLog(
		BaseLogID BIGINT NOT NULL,
		Command VARCHAR(32) NOT NULL,
		Username VARCHAR(32),
		StockSymbol VARCHAR(3),
		Filename VARCHAR(255),
		Funds DECIMAL(32, 2),
		FOREIGN KEY (BaseLogID) REFERENCES BaseLog(ID) ON DELETE CASCADE,
		PRIMARY KEY (BaseLogID)
	)`,
	`CREATE TABLE IF NOT EXISTS QuoteServerLog(
		BaseLogID BIGINT NOT NULL,
		Price DECIMAL(32, 2) NOT NULL,
		StockSymbol VARCHAR(3) NOT NULL,
		Username VARCHAR(32) NOT NULL,
		QuoteServerTime BIGINT NOT NULL,
		CryptoKey VARCHAR(128) NOT NULL,
		FOREIGN KEY (BaseLogID) REFERENCES BaseLog(ID) ON DELETE CASCADE,
		PRIMARY KEY (BaseLogID)
	)`,
	`CREATE TABLE IF NOT EXISTS AccountTransactionLog(
		BaseLogID BIGINT NOT NULL,
		Action VARCHAR(32) NOT NULL,
		Username VARCHAR(32) NOT NULL,
		Funds DECIMAL(32, 2) NOT NULL,
		FOREIGN KEY (BaseLogID) REFERENCES BaseLog(ID) ON DELETE CASCADE,
		PRIMARY KEY (BaseLogID)
	)`,
	`CREATE TABLE IF NOT EXISTS SystemEventLog(
		BaseLogID BIGINT NOT NULL,
		Command VARCHAR(16) NOT NULL,
		Username VARCHAR(32),
		StockSymbol VARCHAR(3),
		Filename VARCHAR(255),
		Funds DECIMAL(32, 2),
		FOREIGN KEY (BaseLogID) REFERENCES BaseLog(ID) ON DELETE CASCADE,
		PRIMARY KEY (BaseLogID)
	)`,
	`CREATE TABLE IF NOT EXISTS ErrorEventLog(
		BaseLogID BIGINT NOT NULL,
		Command VARCHAR(16) NOT NULL,
		Username VARCHAR(32),
		StockSymbol VARCHAR(3),
		Filename VARCHAR(255),
		Funds DECIMAL(32, 2),
		ErrorMessage TINYTEXT,
		FOREIGN KEY (BaseLogID) REFERENCES BaseLog(ID) ON DELETE CASCADE,
		PRIMARY KEY (BaseLogID)
	)`,
	`CREATE TABLE IF NOT EXISTS DebugEventLog(
		BaseLogID BIGINT NOT NULL,
		Command VARCHAR(16) NOT NULL,
		Username VARCHAR(32),
		StockSymbol VARCHAR(3),
		Filename VARCHAR(255),
		Funds DECIMAL(32, 2),
		DebugMessage TINYTEXT,
		FOREIGN KEY (BaseLogID) REFERENCES BaseLog(ID) ON DELETE CASCADE,
		PRIMARY KEY (BaseLogID)
	)`,
}

func createAndOpenDB() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(logging_db:3306)/")
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
	db, err = sql.Open("mysql", "root@tcp(logging_db:3306)/logging")
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Second * 0)
	db.SetMaxIdleConns(10000)
}
