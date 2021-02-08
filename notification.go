package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type MsgStruct struct {
	Content string `json:"content"`
}

func sendNotification(content string) error {
	webhookURL := os.Getenv("DISCORD_WEBHOOK")
	msgJson, _ := json.Marshal(MsgStruct{content})
	payload := bytes.NewReader(msgJson)
	req, err := http.NewRequest("POST", webhookURL, payload)
	if err != nil {
		return fmt.Errorf("sendNotification err: %s", err)
	}

	req.Header.Add("Content-Type", "application/json")
	err = apiRequest(req, nil)
	if err != nil {
		return fmt.Errorf("sendNotification: %s", err)
	}

	return nil
}
