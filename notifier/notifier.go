package notifier

import "errors"

//Notifier ...
type Notifier interface {
	Notify(body string) error
}

//Twilio ...
const Twilio string = "Twilio"

//Screen ...
const Screen string = "Screen"

//Factory ...
func Factory(ntype string) (Notifier, error) {
	if ntype == Twilio {
		return &twilio{}, nil
	} else if ntype == Screen {
		return &screen{}, nil
	}
	return nil, errors.New("Unknown notifier type")
}
