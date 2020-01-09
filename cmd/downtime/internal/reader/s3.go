package reader

import (
	"bufio"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//InputS3Bucket holds configuration for s3
type InputS3Bucket struct {
	//AwsRegion ...
	AwsRegion string
	//Bucket ...
	Bucket string
	//Key ...
	Key string
}

const tmpFile string = "/tmp/downtime_domains"

//GetDomains reads the domain list from S3
func (r *InputS3Bucket) GetDomains() ([]string, error) {

	sess, err := createAWSSession(r.AwsRegion)
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
			Bucket: aws.String(r.Bucket),
			Key:    aws.String(r.Key)})
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

func createAWSSession(awsRegion string) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)})
	if err != nil {
		return nil, err
	}

	return sess, nil
}
