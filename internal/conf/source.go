package conf

import (
	"fmt"
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
func NewPs(awsRegion, namespace string) (*Awsps, error) {

	//Validate region
	if awsRegion == "" {
		return &Awsps{}, errors.New("AWS Region is not defined")
	}

	//Create AWS session
	ses, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)})
	if err != nil {
		return &Awsps{}, err
	}
	sm := ssm.New(ses)

	//Get parameters from AWS parameter store for @namespace
	path := fmt.Sprintf("/%v/", namespace)
	recursive := true
	decryption := false
	params, err := sm.GetParametersByPath(&ssm.GetParametersByPathInput{
		Path:           &path,
		Recursive:      &recursive,
		WithDecryption: &decryption})
	if err != nil {
		return &Awsps{}, err
	}

	//Store parameters in struct's map
	m := make(map[string]string)
	for _, param := range params.Parameters {
		//Use key in Environment variable form from the ArdanLabs conf package
		k := strings.ToUpper(strings.Replace(strings.TrimPrefix(*param.Name, path),
			"/", "_", -1))
		m[k] = *param.Value
	}

	return &Awsps{m: m}, err
}

//Source implements the conf.Sourcer interface
func (ps *Awsps) Source(field conf.Field) (string, bool) {
	k := strings.ToUpper(strings.Join(field.EnvKey, `_`))
	v, ok := ps.m[k]
	return v, ok
}
