package reader

import (
	"strings"
)

//Reader ...
type Reader interface {
	GetDomains() ([]string, error)
}

//Inline ...
const Inline string = "inline"

//S3 ...
const S3 string = "s3"

//Factory ...
func Factory(rtype string) Reader {
	if strings.ToLower(rtype) == S3 {
		return &s3Bucket{}
	}
	return &inline{}
}
