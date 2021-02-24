package main

import (
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

type PriceCompare struct{
	UserId int64 `json:"user_id" db:"user_id"`
	Value float64 `json:"price" db:"value"`
	BuySell bool `json:"buy_sell" db:"buy_sell"`
}

func cronMain() {
	cronManager.AddFunc("0 59 8 * * 1-5", runEveryDay)
	cronManager.AddFunc("0 30 11 * * 1-5", removeDailycron)
	cronManager.AddFunc("0 30 12 * * 1-5", runEveryDay)
	cronManager.AddFunc("0 1 15 * * 1-5", removeDailycron)
}

// add the corn job on every moring AM 8:50
func runEveryDay() {
	var tmp int
	dbm.Get(&tmp, "SELECT COUNT(*) FROM holiday WHERE date=?", time.Now().Format(DateFormat))
	if tmp > 0{
		return
	}
	var rst []string
	// Find DISTINCT stock code from stock_alert
	dbm.Select(&rst, "SELECT DISCTINCT s.code FROM stock_alert AS sa JOIN stock AS s ON sa.stock_id=s.id")
	for _, code := range rst {
		code := code
		entryId, err := cronManager.AddFunc("0 * 9-15 * * *", func() {
			price, err := addStockLog(code)
			
			// Check whether hit alert on specific stock
			checkRst, err := checkStockAlert(code, price)
			if err != nil{
				log.Printf("Daily parse check stock %s err : %s\n", code, err.Error())
				return 
			}
			for _, target := range checkRst {
				if target.BuySell == true{
					// TODO: Find which user to notify or operate
					sendNotification(fmt.Sprintf("%v -> Buy %s now at price: %v",target.UserId, code, price))
				} else if target.BuySell == false {
					sendNotification(fmt.Sprintf("%v -> Sell %s now at price: %v",target.UserId, code, price))
				}
			}
		})

		if err != nil {
			log.Printf("Error on adding cron job %s : %s\n", code, err.Error())
		}
		cronJobManager = append(cronJobManager, entryId)
	}
}

// remove the cron job at PM3:00
func removeDailycron() {
	var rst []string
	for _, jobEntity := range cronJobManager {
		cronManager.Remove(jobEntity)
	}
	cronJobManager = make([]cron.EntryID, 0)
	dbm.Select(&rst, "SELECT DISCTINCT s.code FROM stock_alert AS sa JOIN stock AS s ON sa.stock_id=s.id")
	for _, stock := range rst {
		updateStockData(stock)
	}
}

func addStockLog(code string) (float64, error) {
	_, err := checkStockInfo(code)
	if err != nil {
		return 0, err
	}
	price, err := parseStockPrice(code)
	log.Println(price)
	if err != nil {
		return 0, err
	}
	dbm.MustExec("INSERT INTO `stock_log`(stock_id, price, time) VALUES(?, ?, ?)", code, price, time.Now().UTC().Format("2006-01-02 03:04:05"))
	return price, nil
}

func checkStockAlert(code string, price float64) ([]PriceCompare, error) {
	var minQueryRst []PriceCompare
	var maxQueryRst []PriceCompare
	var rst []PriceCompare
	ok := dbm.Select(&minQueryRst,
		"SELECT MIN(price) AS value, buy_sell, user_id FROM stock_alert WHERE stock_id=? AND cross_direction=false AND alert_on=true GROUP BY buy_sell, user_id", code)
	if ok != nil {
		return rst, fmt.Errorf("check stock price lower bound err: %s", ok)
	} 

	for _, target := range minQueryRst{
		if price <= target.Value {
			rst = append(rst, target)
		}
	}
	
	ok = dbm.Select(&maxQueryRst,
		"select MAX(price) as value, buy_sell, user_id from stock_alert where stock_id=? and cross_direction=true and alert_on=true group by buy_sell, user_id", code)
	if ok != nil{
		return rst,	fmt.Errorf("check stock price higher bound err: %s", ok)
	}

	for _, target := range maxQueryRst{
		if price >= target.Value {
			rst = append(rst, target)
		}
	}
	return rst, nil
}

func closeAlert(target PriceCompare){
	
}