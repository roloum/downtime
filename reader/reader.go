package reader

//Reader ...
type Reader interface {
	GetDomains() ([]string, error)
}

//Inline ...
const Inline string = "Inline"

//S3 ...
const S3 string = "S3"

//Factory ...
func Factory(rtype string) Reader {
	if rtype == S3 {
		return &s3{}
	}
	return &inline{}
}
