package notifier

import (
	"fmt"
	"log"
)

//Screen ...
type Screen struct{}

//Notify ...
func (s *Screen) Notify(body string, log *log.Logger) error {
	log.Println("Delivering output to Screen")
	fmt.Printf("%v\n", body)
	return nil
}
