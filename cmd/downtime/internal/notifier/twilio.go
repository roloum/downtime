package notifier

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

//Twilio struct implements the Notifier interface
type Twilio struct {
	//Sid Twilio Account Id
	Sid string
	//Token Twilio token
	Token string
	//From Sender
	From string
	//To Phone number where message is sent
	To string
}

//Notify sends a text message
func (t *Twilio) Notify(body string, log *log.Logger) error {
	log.Println("Checking Twilio configuration parameters")
	if !t.canSendSMS() {
		return errors.New("Missing Twilio Configuration Parameters")
	}

	uri := fmt.Sprintf(
		"https://api.twilio.com/2010-04-01/Accounts/%v/Messages.json",
		t.Sid)

	log.Println("Encoding Twilio parameters")
	data := url.Values{}
	data.Set("To", t.To)
	data.Set("From", t.From)
	data.Set("Body", body)
	dataReader := *strings.NewReader(data.Encode())

	client := &http.Client{}

	log.Printf("Creating POST request to %v\n", uri)
	req, err := http.NewRequest("POST", uri, &dataReader)
	if err != nil {
		return err
	}

	log.Println("Setting authentication")
	req.SetBasicAuth(t.Sid, t.Token)
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	log.Println("Executing POST request")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	log.Println("Checking HTTP status code")
	if resp.StatusCode != http.StatusCreated {
		return errors.New(resp.Status)
	}

	log.Println("Creating decoding object for Twilio's response")
	var result map[string]interface{}
	decoder := json.NewDecoder(resp.Body)

	log.Println("Decoding response")
	return decoder.Decode(&result)
}

func (t *Twilio) canSendSMS() bool {
	if t.Sid == "" || t.Token == "" || t.From == "" || t.To == "" {
		return false
	}
	return true
}
