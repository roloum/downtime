package reader

import (
	"bufio"
	"log"
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
func (r *InputS3Bucket) GetDomains(log *log.Logger) ([]string, error) {

	log.Println("Reading domains from S3 bucket")

	log.Println("Creating AWS session")
	sess, err := createAWSSession(r.AwsRegion)
	if err != nil {
		return nil, err
	}

	log.Printf("Creating temporary file: %v\n", tmpFile)
	file, err := os.Create(tmpFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	log.Println("Downloading data from S3 to temporary file")
	downloader := s3manager.NewDownloader(sess)
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(r.Bucket),
			Key:    aws.String(r.Key)})
	if err != nil {
		return nil, err
	}

	if numBytes == 0 {
		log.Println("S3 bucket is empty")
		return []string{}, nil
	}

	log.Println("Reading domains from temporary file")
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
