package reader

import (
	"log"
	"os"
)

//InputInline ...
type InputInline struct{}

//GetDomains reads parameters from command line
func (r *InputInline) GetDomains(log *log.Logger) ([]string, error) {
	log.Println("Reading domains from command line")
	return os.Args[1:], nil
}
