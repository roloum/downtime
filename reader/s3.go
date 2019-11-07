package reader

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type s3Bucket struct{}

func (r *s3Bucket) GetDomains() ([]string, error) {
	sess, err := createAWSSession()
	if err != nil {
		return nil, err
	}

	file, err := os.Create("/tmp/downtime_domains")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(sess)
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(os.Getenv("DOWNTIME_S3_BUCKET")),
			Key:    aws.String("DOWNTIME_S3_KEY")})
	if err != nil {
		return nil, err
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")

	return []string{}, nil
}

func createAWSSession() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")})
	if err != nil {
		return nil, err
	}
	return sess, nil
}
