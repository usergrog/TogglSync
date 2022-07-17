package models

type Config struct {
	JiraToken    string
	JiraHost     string
	JiraUsername string
}

type TogglEntry struct {
	Id          int     `json:"id"`
	Price       float64 `json:"price"`
	Start       string  `json:"start"`
	Stop        string  `json:"stop"`
	Description string  `json:"description"`
	Duration    int64   `json:"duration"`
}

type Worklog struct {
	Started          string `json:"started"`
	TimeSpentSeconds int64  `json:"timeSpentSeconds"`
}
