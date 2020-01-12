package notifier

import "log"

//Notifier ...
type Notifier interface {
	Notify(body string, log *log.Logger) error
}

//Output ...
type Output struct{}

//Notify calls the Notify method in the Notifier object to deliver message body
func (o *Output) Notify(n Notifier, body string, log *log.Logger) error {
	return n.Notify(body, log)
}
