package config

import (
	"encoding/json"
	"fmt"
	"os"
)

//Config : struct for reading config.json file
type Config struct {
	DBHost string `json:"dbhost"`
	DBPass string `json:"dbpass"`
	DBUser string `json:"dbuser"`
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
