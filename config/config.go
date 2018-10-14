package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"
)

//Config : struct for reading config.json file
type Config struct {
	DBHost       string `json:"dbhost"`
	DBPass       string `json:"dbpass"`
	DBUser       string `json:"dbuser"`
	DBName       string `json:"dbname"`
	StripeSecret string `json:"stripe_secret"`
}

type AWSCredentials struct {
	AWSAccessKeyID  string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

// LoadConfiguration : loads the info in from the config.json file
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

// LoadAWSConfiguration : loads AWS Configuration and returns the credentials
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
