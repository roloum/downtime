package main

import "testing"

func TestDowntime200(t *testing.T) {
	url := "httpstat.us/200"
	errors := checkDownTime([]string{url})
	if len(errors) != 0 {
		t.Errorf(
			"URL %v that should have returned 200 returned the following error: %v",
			url, errors[0])
	}
}

func TestDowntimeInvalidUrl(t *testing.T) {
	url := "http://"
	errors := checkDownTime([]string{url})
	if len(errors) != 1 {
		t.Errorf(
			"URL %v that should have been invalid passed the test", url)
	}
}

func TestDowntimeNonExistent(t *testing.T) {
	url := "http://desperateprogrammers.com"
	errors := checkDownTime([]string{url})
	if len(errors) != 1 {
		t.Errorf(
			"URL %v that should have been invalid passed the test", url)
	}
}

func TestDowntime400(t *testing.T) {
	url := "httpstat.us/400"
	errors := checkDownTime([]string{url})
	if len(errors) != 1 {
		t.Errorf(
			"URL %v that should have been invalid passed the test", url)
	}
}
