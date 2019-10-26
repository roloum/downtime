package notifier

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const sid string = "TWILIO_SID"
const token string = "TWILIO_AUTH_TOKEN"
const from string = "TWILIO_FROM"
const to string = "TWILIO_TO"

type twilio struct {
}

func canSendSMS() bool {
	if os.Getenv(sid) == "" || os.Getenv(token) == "" || os.Getenv(from) == "" ||
		os.Getenv(to) == "" {
		return false
	}
	return true
}

func (t *twilio) Notify(body string) error {

	if !canSendSMS() {
		return errors.New("Missing Twilio Configuration Parameters")
	}

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

	return decoder.Decode(&result)
}
