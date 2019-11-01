package reader

type s3 struct{}

func (r *s3) GetDomains() ([]string, error) {
	return []string{}, nil
}
