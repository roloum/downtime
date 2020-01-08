package url

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

var (
	//ErrInvalidURL error returned in case URL does not have valid format
	ErrInvalidURL = errors.New("Invalid URL")
	//ErrCouldNotCheck returned in case the domain does not exist
	ErrCouldNotCheck = errors.New("Could not check URL")
)

//Check checks if the URL has a valid format and if the domain is up
func Check(uri string, domain bool) error {

	if !(strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://")) {
		uri = "http://" + uri
	}

	u, err := url.ParseRequestURI(uri)
	if err != nil || !strings.Contains(u.Host, ".") {
		return ErrInvalidURL
	}

	//Remove path if we're checking the domain
	if domain {
		uri = strings.TrimSuffix(uri, u.Path)
	}

	if resp, err := http.Get(uri); err != nil ||
		resp.StatusCode != http.StatusOK {
		return ErrCouldNotCheck
	}

	return nil

}
