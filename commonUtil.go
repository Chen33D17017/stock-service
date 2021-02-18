package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type ErrorMsg struct {
	Msg string `json:"msg"`
}

type ResponseMsg struct{
	Success bool `json:"success"`
	Data interface{} `json:"data"`
}

func apiRequest(req *http.Request, response interface{}) error {
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("apiRequest err %s", err.Error())
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("apiRequest err: %s", err.Error())
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return fmt.Errorf("apiRequest err %s", err.Error())
	}
	return nil
}

func successResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := ResponseMsg{true, data}
	json.NewEncoder(w).Encode(response)
}

func errorResponse(w http.ResponseWriter, errMsg string) {
	log.Println(errMsg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	response := ResponseMsg{false, ErrorMsg{errMsg}}
	json.NewEncoder(w).Encode(response)
}
