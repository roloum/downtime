package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/roloum/downtime/notifier"
)

var wg sync.WaitGroup

func main() {

	uris := os.Args[1:]

	errors := checkDownTime(uris)

	if len(errors) == 0 {
		return
	}

	sep := "\n\t- "
	body := fmt.Sprintf("The following domains returned errors:%v%v", sep,
		strings.Join(errors, sep))

	types := []string{}
	//types = append(types, notifier.Screen)

	if os.Getenv("TWILIO") != "" {
		types = append(types, notifier.Twilio)
	}

	for _, ntype := range types {
		n, err := notifier.Factory(ntype)
		if err != nil {
			fmt.Printf("Unsupported notifier type: %v", ntype)
		}
		n.Notify(body)
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
