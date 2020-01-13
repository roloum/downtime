package conf

import (
	"fmt"
	"log"
	"strings"

	"github.com/ardanlabs/conf"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pkg/errors"
)

//This package implements a Sourcer for the ArdanLabs conf package
//reading the parameters stored in AWS Parameter Store

//Awsps AWS Paramter Store
type Awsps struct {
	m map[string]string
}

//NewPs Returns an instance of the parameter store struct after it creates the AWS
//Session and loads all the parameters for the given namespace in the map
//@namespace makes this package generic, if we ever want to be able to
//Use this configuration class for any other project
func NewPs(awsRegion, namespace string, l *log.Logger) (*Awsps, error) {

	l.Println("Verifying AWS Region is set")
	//Validate region
	if awsRegion == "" {
		return &Awsps{}, errors.New("AWS Region is not defined")
	}

	//Create AWS session
	l.Println("Creating AWS session")
	ses, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)})
	if err != nil {
		return &Awsps{}, err
	}
	sm := ssm.New(ses)

	//Map for configuration
	m := make(map[string]string)

	//Get parameters from AWS parameter store for @namespace
	l.Printf("Loading parameters for namespace: %v\n", namespace)

	path := fmt.Sprintf("/%v/", namespace)
	recursive := true
	decryption := false
	maxResults := int64(10)
	nextToken := ""

	parametersInput := &ssm.GetParametersByPathInput{
		MaxResults:     &maxResults,
		Path:           &path,
		Recursive:      &recursive,
		WithDecryption: &decryption}

	//Download parameters from Parameter Store using pagination
	for {

		if nextToken != "" {
			l.Println("Loading parameters for next token")
			parametersInput.NextToken = &nextToken
		}

		params, err := sm.GetParametersByPath(parametersInput)
		if err != nil {
			return &Awsps{}, err
		}

		if len(params.Parameters) > 0 {

			//Store parameters in struct's map
			l.Println("Storing parameters in map")
			for _, param := range params.Parameters {
				//Use key in Environment variable form from the ArdanLabs conf package
				k := strings.ToUpper(strings.Replace(strings.TrimPrefix(*param.Name, path),
					"/", "_", -1))
				m[k] = *param.Value
			}

			//Are there more parameters to load from AWS
			if params.NextToken != nil {
				nextToken = *params.NextToken
				l.Printf("Next Token: %v\n", nextToken)
			} else {
				//We have reached the last page of parameters. Exit loop
				l.Println("Reached the last page of parameters")
				break
			}

		} else {
			//Parameters array is empty. Exit loop
			l.Println("Parameters array is empty")
			break
		}
	}

	return &Awsps{m: m}, err
}

//Source implements the conf.Sourcer interface
func (ps *Awsps) Source(field conf.Field) (string, bool) {
	k := strings.ToUpper(strings.Join(field.EnvKey, `_`))
	v, ok := ps.m[k]

	return v, ok
}
