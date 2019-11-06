package reader

import "errors"

type s3 struct{}

func (r *s3) GetDomains() ([]string, error) {
	return nil, errors.New("Implement")
}
