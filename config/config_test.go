package config

import (
	"encoding/json"
	"os"
	"testing"
)

func TestLoadConfiguration(t *testing.T) {
	var config Config
	configFile, err := os.Open("./../config.json")
	defer configFile.Close()
	if err != nil {
		t.Error(err.Error())
	}

	configStats, err := configFile.Stat()
	if err != nil {
		t.Error(err.Error())
	}

	if configStats.Size() <= 0 {
		t.Error("Opened empty file")
	}

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return
}
