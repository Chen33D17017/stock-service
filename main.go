package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/robfig/cron"
)

var cronManager *cron.Cron
var dbm *sqlx.DB

const DateFormat = "2006-01-02"

var wg *sync.WaitGroup

func init() {
	cronManager = cron.New()
	dbm = NewDBManager()

}

func main() {
	http.HandleFunc("/getStock", test)
	http.HandleFunc("/", hello)
	http.HandleFunc("/price", price)
	http.HandleFunc("/notification", notification)
	http.HandleFunc("/stock", stockInfo)
	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func notification(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		msg := MsgStruct{}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		err = json.Unmarshal(body, &msg)
		sendNotification(msg.Content)
	}
}

func price(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		fmt.Fprint(w, "URL Param 'key' is missing")
	}
	price, err := parseStockPrice(keys[0])
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
	fmt.Fprint(w, price)
}

func test(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		fmt.Fprint(w, "URL Param 'key' is missing")
	}
	parseStockDetail(keys[0])
}
