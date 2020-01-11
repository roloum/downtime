package reader

import "log"

//Reader ...
type Reader interface {
	GetDomains(log *log.Logger) ([]string, error)
}

//S3 ...
const S3 string = "s3"

//Input ...
type Input struct{}

//Read returns list of domains from object implementing Reader interface
func (i *Input) Read(r Reader, log *log.Logger) ([]string, error) {
	return r.GetDomains(log)
}
