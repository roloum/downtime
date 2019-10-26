package notifier

import "fmt"

//Screen ...
type screen struct {
}

//Notify ...
func (s *screen) Notify(body string) error {
	fmt.Printf("%v\n", body)
	return nil
}
