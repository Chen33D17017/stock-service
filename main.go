package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/robfig/cron/v3"
)

var cronManager *cron.Cron
var dbm *sqlx.DB
var cronJobManager []cron.EntryID
var beginDate time.Time

const DateFormat = "2006-01-02"

var wg *sync.WaitGroup

func init() {
	cronManager = cron.New()
	dbm = NewDBManager()
	beginDate , _ = time.Parse(DateFormat, "2000-01-01")
}

func main() {
	http.HandleFunc("/alertCheck", alertCheck)
	http.HandleFunc("/price", getPrice)
	http.HandleFunc("/notification", notification)
	http.HandleFunc("/regist", registerAlert)
	http.HandleFunc("/holiday", registHoliday)
	http.ListenAndServe(":8080", nil)
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

func getPrice(w http.ResponseWriter, r *http.Request) {
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

func alertCheck(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["stock"]
	if !ok || len(keys[0]) < 1 {
		fmt.Fprint(w, "URL Param 'key' is missing")
		return
	}
	price, _ := r.URL.Query()["price"]
	if !ok || len(price[0]) < 1 {
		fmt.Fprintf(w, "URL param 'price' is missing'")
	}
	p, _ := strconv.ParseFloat(price[0], 64)
	rst, err := checkStockAlert(keys[0], p)
	if err != nil {
		errorResponse(w, err.Error())
	}

	successResponse(w, rst)
}

func updateStock(w http.ResponseWriter, r *http.Request){
	keys, ok := r.URL.Query()["stock"]
	if !ok || len(keys[0]) < 1 {
		fmt.Fprint(w, "URL Param 'key' is missing")
		return
	}

	updateStockData(keys[0])
	successResponse(w, fmt.Sprintf("Update stock %s", keys[0]))
}

