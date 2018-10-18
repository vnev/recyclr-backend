// Package config reads our configuration for our database and AWS.
package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"
)

// Config is the struct which will hold our various configuration settings such as
// the database username and password, and our Stripe secret.
type Config struct {
	DBHost       string `json:"dbhost"`
	DBPass       string `json:"dbpass"`
	DBUser       string `json:"dbuser"`
	DBName       string `json:"dbname"`
	StripeSecret string `json:"stripe_secret"`
}

// AWSCredentials is the struct which holds our access information for AWS
type AWSCredentials struct {
	AWSAccessKeyID  string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

// LoadConfiguration loads the info in from a JSON file and returns the resulting struct.
func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

// LoadAWSConfiguration loads AWS Configuration and returns the credentials.
func LoadAWSConfiguration(file string) (*credentials.Credentials, error) {
	//var creds credentials.Credentials
	var awsCreds AWSCredentials
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&awsCreds)
	creds := credentials.NewStaticCredentials(awsCreds.AWSAccessKeyID, awsCreds.SecretAccessKey, "")
	_, err = creds.Get()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return creds, err
}
