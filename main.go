package main

import (
	"TogglSync/utils"
	"log"
	"net/http"
)

func main() {
	log.Println("Start toggl-sync")

	config := utils.ReadConfig()

	client := &http.Client{}

	req, _ := http.NewRequest("GET", "https://api.localeapp.com/v1/projects/"+config.ApiKey+"/translations/all.yml", nil)

	resp, _ := client.Do(req)

	log.Println(resp.Status)

}
