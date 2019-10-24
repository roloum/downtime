package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const sid = "TWILIO_SID"
const token = "TWILIO_AUTH_TOKEN"
const from = "TWILIO_FROM"
const to = "TWILIO_TO"

func init() {
	if os.Getenv(sid) == "" || os.Getenv(token) == "" || os.Getenv(from) == "" ||
		os.Getenv(to) == "" {
		panic("Twilio auth tokens not defined")
	}
}

func sendMessage(body string) error {

	uri := fmt.Sprintf(
		"https://api.twilio.com/2010-04-01/Accounts/%v/Messages.json",
		os.Getenv(sid))

	data := url.Values{}
	data.Set("To", os.Getenv(to))
	data.Set("From", os.Getenv(from))
	data.Set("Body", body)
	dataReader := *strings.NewReader(data.Encode())

	client := &http.Client{}
	req, err := http.NewRequest("POST", uri, &dataReader)
	if err != nil {
		return err
	}

	req.SetBasicAuth(os.Getenv(sid), os.Getenv(token))
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return errors.New(resp.Status)
	}

	var result map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	if err == nil {
		fmt.Println(result)
	}

	return nil
}
