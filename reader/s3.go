package reader

import (
	"bufio"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type s3Bucket struct{}

const tmpFile string = "/tmp/downtime_domains"

func (r *s3Bucket) GetDomains() ([]string, error) {

	sess, err := createAWSSession()
	if err != nil {
		return nil, err
	}

	file, err := os.Create(tmpFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(sess)
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(os.Getenv("DOWNTIME_S3_BUCKET")),
			Key:    aws.String(os.Getenv("DOWNTIME_S3_KEY"))})
	if err != nil {
		return nil, err
	}

	if numBytes == 0 {
		return []string{}, nil
	}

	urls := []string{}
	data := bufio.NewScanner(file)
	for data.Scan() {
		urls = append(urls, data.Text())
	}

	return urls, nil
}

func createAWSSession() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")})
	if err != nil {
		return nil, err
	}
	return sess, nil
}
