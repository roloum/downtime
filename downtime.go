package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/roloum/downtime/notifier"
	"github.com/roloum/downtime/reader"
)

var wg sync.WaitGroup

func main() {

	uris, err := reader.Factory(os.Args[1]).GetDomains()
	if err != nil {
		//// TODO: Print error on screen for now, but it should log or email
		fmt.Printf("Error reading domains: %v\n", err)
		return
	}

	errors := checkDownTime(uris)

	if len(errors) == 0 {
		return
	}

	sep := "\n\t- "
	body := fmt.Sprintf("The following domains returned errors:%v%v", sep,
		strings.Join(errors, sep))

	notify(body)
}

func notify(body string) {

	types := []string{}

	if os.Getenv("DOWNTIME_SCREEN") != "" {
		types = append(types, notifier.Screen)
	}
	if os.Getenv("DOWNTIME_TWILIO") != "" {
		types = append(types, notifier.Twilio)
	}

	for _, ntype := range types {
		n, err := notifier.Factory(ntype)
		if err != nil {
			fmt.Printf("Unsupported notifier type: %v", ntype)
		}

		//// TODO: Print error on screen for now, but it should log or email
		if err = n.Notify(body); err != nil {
			fmt.Println(err)
		}
	}

}

func checkDownTime(uris []string) (errors []string) {

	count := len(uris)
	var ch = make(chan string, count)

	wg.Add(count)

	for _, uri := range uris {
		go checkurl(uri, ch)
	}

	wg.Wait()
	close(ch)

	for error := range ch {
		errors = append(errors, error)
	}

	return
}

func checkurl(uri string, ch chan string) {

	defer wg.Done()

	var msg string

	if !strings.HasPrefix(uri, "http") {
		uri = "http://" + uri
	}

	if _, error := url.ParseRequestURI(uri); error != nil {
		msg = fmt.Sprintf("Invalid URL: %v", uri)
	} else if resp, error := http.Get(uri); error != nil {
		msg = fmt.Sprintf("Could not check URL: %v. %v.", uri, resp)
	} else if resp.StatusCode != http.StatusOK {
		msg = fmt.Sprintf("Error from %v. Code: %v. Message: %v", uri, resp.StatusCode,
			resp.Status)
	}

	if msg != "" {
		ch <- msg
	}
}
