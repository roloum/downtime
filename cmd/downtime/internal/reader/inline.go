package reader

import (
	"os"
)

//InputInline ...
type InputInline struct{}

//GetDomains reads parameters from command line
func (r *InputInline) GetDomains() ([]string, error) {
	return os.Args[1:], nil
}
