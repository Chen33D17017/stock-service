package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ErrorMsg struct {
	Msg string `json:"msg"`
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

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func errorResponse(w http.ResponseWriter, errMsg string){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(ErrorMsg{errMsg})
}