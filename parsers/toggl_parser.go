package parsers

import (
	"TogglSync/models"
	"TogglSync/utils"
	"encoding/json"
)

func ParseJson(jsonBody []byte) []models.TogglEntry {

	var responseBody = []models.TogglEntry{}

	//fmt.Println(string(jsonBody))

	err := json.Unmarshal(jsonBody, &responseBody)
	utils.CheckError(err)

	return responseBody
}
