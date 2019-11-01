package reader

import (
	"os"
)

type inline struct{}

func (r *inline) GetDomains() ([]string, error) {
	return os.Args[1:], nil
}
