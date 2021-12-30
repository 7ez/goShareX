package main

import (
	"encoding/json"

	"os"
)

var Config Configuration

type Configuration struct {
	Port      string `json:"Port"`
	Domain    string `json:"Domain"`
	UploadKey string `json:"UploadKey"`
	AppName   string `json:"AppName"`
}

func LoadConfig() {
	var config Configuration
	configFile, err := os.Open("./config.json")
	defer configFile.Close()
	if err != nil {
		panic(err)
	}

	json := json.NewDecoder(configFile)
	json.Decode(&config)

	Config = config
}
