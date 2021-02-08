package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type RegisterStockRequest struct {
	UserID    int64   `json:"userId"`
	StockCode string  `json:"stockCode"`
	Stradegy  string  `json:"stradegy"`
	Price     float64 `json:"price"`
}

type StockName struct {
	StockID   int    `db:"id" json:"-"`
	StockCode string `db:"code" json:"stockCode"`
	StockName string `db:"name" json:"stockName"`
}

func printErr(w http.ResponseWriter, err error) {
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
}

func successResponse(w http.ResponseWriter, content string) {
	fmt.Fprintf(w, content)
}

// RegisterStockRequest
func stockInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		lock := make(chan string, 1)
		req := RegisterStockRequest{}
		body, err := ioutil.ReadAll(r.Body)
		printErr(w, err)
		err = json.Unmarshal(body, &req)
		printErr(w, err)

		var result StockName
		err = dbm.Get(&result, `SELECT * FROM stock WHERE code=?`, req.StockCode)
		// insert stock info if not exists
		if err != nil {
			log.Printf("Cannot find stock : %s \ninsert data in to database", err.Error())
			stockName, err := parseStockName(req.StockCode)
			if err != nil {
				printErr(w, err)
				return
			}
			result = StockName{0, req.StockCode, stockName}
			tx := dbm.MustBegin()
			_, err = tx.NamedExec("INSERT INTO stock(code, name) VALUES (:code, :name)", &result)
			err = tx.Commit()
			lock <- req.StockCode
			if err != nil {
				printErr(w, err)
				return
			}
			go func(chan string) {
				code := <-lock
				parseStockDetail(code)
			}(lock)
		}

		jsonResponse(w, result)

	}
}
