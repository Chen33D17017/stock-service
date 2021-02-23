package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type Holiday struct {
	Date string `json:"date" db:"date"`
}

func registHoliday(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		req := Holiday{}
		body, err := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(body, &req)
		if err != nil {
			errorResponse(w, err.Error())
			return
		}

		_, err = time.Parse(DateFormat, req.Date)
		if err != nil {
			errorResponse(w, err.Error())
			return
		}

		// add alert data into db
		_, err = dbm.NamedExec(`INSERT INTO holiday (date) VALUES(:date)`, req)
		if err != nil {
			errorResponse(w, err.Error())
			return
		}
		successResponse(w, req)
	}
}
