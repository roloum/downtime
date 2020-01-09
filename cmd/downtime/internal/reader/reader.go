package reader

//Reader ...
type Reader interface {
	GetDomains() ([]string, error)
}

//Inline ...
const Inline string = "inline"

//S3 ...
const S3 string = "s3"

//Input ...
type Input struct{}

//Read returns list of domains from object implementing Reader interface
func (i *Input) Read(r Reader) ([]string, error) {
	return r.GetDomains()
}
