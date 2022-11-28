package main

import (
	"TogglSync/models"
	"TogglSync/parsers"
	"TogglSync/utils"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"time"
)

var config models.Config
var httpClient *http.Client
var help = flag.Bool("help", false, "Show help")
var startFrom = flag.String("start", "", "Date from")
var endTo = flag.String("stop", "", "Date to")

func main() {
	log.Println("Start toggl-sync")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// fmt.Println("start ", *startFrom)
	// fmt.Println("stop ", *endTo)

	models.InitDb()
	config = utils.ReadConfig()

	httpClient = &http.Client{}

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.track.toggl.com/api/v9/me/time_entries?start_date=%s&end_date=%s", *startFrom, *endTo), nil)
	req.SetBasicAuth(config.TogglToken, "api_token")
	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	utils.CheckError(err)
	log.Println("Toggl: ", resp.Status)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	utils.CheckError(err)

	// fmt.Println(string(body))
	togglEntries := parsers.ParseJson(body)

	for _, entry := range togglEntries {
		fmt.Println("Start at", entry.Start, "Desc: ", entry.Description, "\t\t\tDuration", entry.Duration, "\tHuman:", transformToHumanView(entry.Duration))
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
			Started:   parsedDate.Format("2006-01-02T15:04:05.000+0000"),
			TimeSpent: transformToHumanView(entry.Duration),
		}
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(postBody)
		url := fmt.Sprintf("https://%s.atlassian.net/rest/api/3/issue/%s/worklog", config.JiraHost, ticketId)
		fmt.Println("url: ", url)
		req, _ := http.NewRequest(http.MethodPost, url, buf)
		req.SetBasicAuth(config.JiraUsername, config.JiraToken)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// fmt.Scanln()
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
	re, _ := regexp.Compile(`^(TIC|INT|TCR|TPI|TPFA|TPC|TPFB|MSV|CUST|TCRD)-\d{1,6}`)
	res := re.FindAllString(entry.Description, -1)
	if len(res) > 0 {
		return res[0]
	} else {
		return ""
	}
}

func transformToHumanView(duration int64) string {
	hours := math.Floor(float64(duration) / 60 / 60)
	minutes := math.Ceil((float64(duration) - hours*60*60) / 60)
	return fmt.Sprintf("%dh %dm", int(hours), int(minutes))
}
