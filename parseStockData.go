package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
)

const INSTANT_STOCK_ADDR = "https://stocks.finance.yahoo.co.jp/stocks/detail/?code=%s.T"
const SOURCE_ADDR = "https://info.finance.yahoo.co.jp/history/?code=%s.T&sy=2000&sm=1&sd=1&ey=%d&em=%d&ed=%d&tm=d"

type stockDetail struct {
	Date  string
	Open  float64
	High  float64
	Low   float64
	Close float64
	Vol   float64
}

func parseStockPrice(sc string) (string, error) {
	addr := fmt.Sprintf(INSTANT_STOCK_ADDR, sc)
	resp, err := soup.Get(addr)
	if err != nil {
		return "", fmt.Errorf("parseStockData err: %s", err)
	}
	doc := soup.HTMLParse(resp)
	price := doc.FindAll("td", "class", "stoksPrice")
	return price[1].Text(), nil
}

func parseStockName(sc string) (string, error) {
	addr := fmt.Sprintf(INSTANT_STOCK_ADDR, sc)
	resp, err := soup.Get(addr)
	if err != nil {
		return "", fmt.Errorf("parsStockName err: %s", err)
	}
	doc := soup.HTMLParse(resp)
	return doc.Find("th", "class", "symbol").Find("h1").Text(), nil
}

func parseStockDetail(sc string) {
	location, _ := time.LoadLocation("Asia/Tokyo")
	today := time.Now().In(location)
	addr := fmt.Sprintf(SOURCE_ADDR,
		sc, today.Year(), today.Month(), today.Day())
	
	var id int
	ok := dbm.Get(&id, `SELECT id FROM stock WHERE code=?`, sc)
	if ok != nil{
		log.Printf("parseStockDetail: fail to get stock %s", sc)
	}
	worklist := make(chan []int, 1)
	unseenPage := make(chan int)
	

	worklist <- []int{1}
	n := 1

	for i := 0; i < 10; i++ {
		go func() {
			for page := range unseenPage {
				targetPage := fmt.Sprintf("%s&p=%d", addr, page)
				resp, err := soup.Get(targetPage)
				if err != nil {
					log.Fatalf("parseStockDetail err %s@%d : %s", sc, page, err.Error())
				}
				doc := soup.HTMLParse(resp)
				pages := parsePageNumber(doc)
				parseStockPage(doc, id)
				go func() {
					worklist <- pages
				}()
			}
		}()
	}

	pageSeen := make(map[int]bool)
	for ; n > 0; n-- {
		list := <-worklist
		for _, page := range list {
			if !pageSeen[page] {
				pageSeen[page] = true
				n++
				unseenPage <- page
			}
		}
	}
}

func parseStockPage(doc soup.Root, id int) {
	elems := doc.Find("div", "class", "padT12").Find("tbody").FindAll("tr")
	//var rst stockDetail
	tx := dbm.MustBegin()
	for i, elem := range elems {
		if i == 0 {
			// pass the column name
			continue
		}
		data := elem.FindAll("td")
		// need to consider the case like 分割: 1株 -> 2株
		if len(data) == 7 {
			rst := stockDetail{getDate(data[0].Text()), getStockVal(data[1].Text()), getStockVal(data[2].Text()),
				getStockVal(data[3].Text()), getStockVal(data[4].Text()), getStockVal(data[5].Text())}
			tx.MustExec("INSERT INTO stock_data(stock_id, price_at, open, high, low, close, vol) VALUES (?, ?, ?, ?, ?, ?, ?)",
				id, rst.Date, rst.Open, rst.High, rst.Low, rst.Close, rst.Vol)
		}
	}
	tx.Commit()
}

func parsePageNumber(doc soup.Root) []int {
	var rst []int
	target := doc.Find("ul", "class", "ymuiPagingBottom")

	for _, data := range target.FindAll("a") {
		content := data.Text()
		page, err := strconv.Atoi(content)
		if err == nil {
			rst = append(rst, page)
		}
	}
	return rst
}

func getDate(s string) string {
	re := regexp.MustCompile(`\d+`)
	reRst := re.FindAllString(s, -1)
	year, yearErr := strconv.Atoi(reRst[0])
	month, monthErr := strconv.Atoi(reRst[1])
	day, dayErr := strconv.Atoi(reRst[2])
	if dayErr != nil || monthErr != nil || yearErr != nil {
		log.Println("parseDate: fail to convert date")
		return ""
	}
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC).Format(DateFormat)
}

func getStockVal(s string) float64 {
	s = strings.ReplaceAll(s, ",", "")
	rst, err := strconv.ParseFloat(s, 64)

	if err != nil {
		log.Printf("Convert error with float: %v\n", err)
		return rst
	}
	return rst
}
