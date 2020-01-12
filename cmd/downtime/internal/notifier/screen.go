package notifier

import (
	"fmt"
	"log"
)

//Screen struct will implement the Notifier interface
type Screen struct{}

//Notify sends output to screen
func (s *Screen) Notify(body string, log *log.Logger) error {
	log.Println("Delivering output to Screen")
	fmt.Printf("%v\n", body)
	return nil
}
