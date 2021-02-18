package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type StockName struct {
	StockCode string `db:"id" json:"stockCode"`
	StockName string `db:"name" json:"stockName"`
}

type StockAlert struct {
	UserID  int64   `db:"user_id" json:"user_id"`
	StockId string  `db:"stock_id" json:"stock_id"`
	BuySell bool    `db:"buy_sell" json:"buy_sell"`
	Cross   bool    `db:"cross_direction" json:"cross_direction"`
	Price   float64 `db:"price" json:"price"`
}

// Register for alert
func registerAlert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		req := StockAlert{}
		body, err := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(body, &req)
		if err != nil {
			errorResponse(w, err.Error())
			return
		}
		// check whether stock is in db, if not add it
		result, err := checkStockInfo(req.StockId)
		if err != nil {
			errorResponse(w, err.Error())
			return
		}
		// add alert data into db
		_, err = dbm.NamedExec(`INSERT INTO stock_alert (user_id, stock_id, buy_sell, cross_direction, price) 
		VALUES(:user_id, :stock_id, :buy_sell, :cross_direction, :price)`, req)
		if err != nil {
			errorResponse(w, err.Error())
			return
		}
		successResponse(w, result)
	}
}

// check whether stock in db, if not insert and parse the daily data
func checkStockInfo(code string) (StockName, error) {
	lock := make(chan string, 1)
	var result StockName
	err := dbm.Get(&result, "SELECT * FROM stock WHERE id=?", code)
	// insert stock info if not exists
	if err != nil {
		log.Printf("Cannot find stock : %s \ninsert data in to database", err.Error())
		stockName, err := parseStockName(code)
		if err != nil {
			return result, err
		}
		result = StockName{code, stockName}
		tx := dbm.MustBegin()
		_, err = tx.NamedExec("INSERT INTO stock(id, name) VALUES (:id, :name)", &result)
		err = tx.Commit()
		lock <- code
		if err != nil {
			return result, err
		}
		go func(chan string) {
			code := <-lock
			parseStockDetail(code)
		}(lock)
	}
	return result, nil
}
