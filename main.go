package main

import (
	"TogglSync/models"
	"TogglSync/parsers"
	"TogglSync/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

var config models.Config
var httpClient *http.Client

func main() {
	log.Println("Start toggl-sync")

	models.InitDb()
	config = utils.ReadConfig()

	httpClient = &http.Client{}

	req, _ := http.NewRequest(http.MethodGet, "https://api.track.toggl.com/api/v9/me/time_entries?start_date=2022-07-11&end_date=2022-07-17", nil)
	req.SetBasicAuth("a524d4522cf319244dd3625e6ec03ff0", "api_token")
	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	utils.CheckError(err)
	log.Println("Toggl: ", resp.Status)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	utils.CheckError(err)

	togglEntries := parsers.ParseJson(body)

	for _, entry := range togglEntries {
		db_entry := models.GetEntry(entry.Id)
		if (db_entry == models.TogglEntry{} && entry.Stop != "") {
			handleNewEntry(entry)
		}
	}

}

func handleNewEntry(entry models.TogglEntry) {
	ticketId := cutTicketId(entry)
	if len(ticketId) > 0 {
		fmt.Printf("%s - %s\n", ticketId, entry.Start)
		layout := "2006-01-02T15:04:05-07:00"
		parsedDate, err := time.Parse(layout, entry.Start)
		utils.CheckError(err)
		//postBody := fmt.Sprintf("{\"started\":\"%s\",\"timeSpentSeconds\":%d}", parsedDate.Format("2006-01-02T15:04:05.000+0000"), entry.Duration)
		postBody := &models.Worklog{
			Started:          parsedDate.Format("2006-01-02T15:04:05.000+0000"),
			TimeSpentSeconds: entry.Duration,
		}
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(postBody)
		url := fmt.Sprintf("https://%s.atlassian.net/rest/api/3/issue/%s/worklog", config.JiraHost, ticketId)
		fmt.Println("url: ", url)
		req, _ := http.NewRequest(http.MethodPost, url, buf)
		req.SetBasicAuth(config.JiraUsername, config.JiraToken)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		resp, err := httpClient.Do(req)
		utils.CheckError(err)
		log.Println("Jira: ", resp.Status)
		if resp.StatusCode == http.StatusCreated {
			models.InsertLocation(&entry)
		}
		//defer resp.Body.Close()
		//body, err := ioutil.ReadAll(resp.Body)
		//fmt.Println(string(body))
	}

}

func cutTicketId(entry models.TogglEntry) string {
	re, _ := regexp.Compile(`^TIC-\d{5}`)
	res := re.FindAllString(entry.Description, -1)
	if len(res) > 0 {
		return res[0]
	} else {
		return ""
	}
}
