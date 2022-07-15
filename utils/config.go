package utils

import (
	"TogglSync/models"
	"encoding/json"
	"os"
)

func ReadConfig() models.Config {
	file, err := os.Open("config.json")
	CheckError(err)
	defer func(file *os.File) {
		err := file.Close()
		CheckError(err)
	}(file)
	decoder := json.NewDecoder(file)
	configuration := models.Config{}
	err = decoder.Decode(&configuration)
	CheckError(err)
	return configuration
}
