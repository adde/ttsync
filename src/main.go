package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type TimeEntry struct {
	ID              int        `json:"id"`
	WorkspaceID     int        `json:"workspace_id"`
	ProjectID       int        `json:"project_id"`
	TaskID          *int       `json:"task_id"`
	Billable        bool       `json:"billable"`
	Start           time.Time  `json:"start"`
	Stop            time.Time  `json:"stop"`
	Duration        int        `json:"duration"`
	Description     string     `json:"description"`
	Tags            []string   `json:"tags"`
	TagIDs          []int      `json:"tag_ids"`
	Duronly         bool       `json:"duronly"`
	At              time.Time  `json:"at"`
	ServerDeletedAt *time.Time `json:"server_deleted_at"`
	UserID          int        `json:"user_id"`
	UID             int        `json:"uid"`
	Wid             int        `json:"wid"`
	Pid             int        `json:"pid"`
}

type Project struct {
	ID                  int       `json:"id"`
	WorkspaceID         int       `json:"workspace_id"`
	ClientID            int       `json:"client_id"`
	Name                string    `json:"name"`
	IsPrivate           bool      `json:"is_private"`
	Active              bool      `json:"active"`
	At                  time.Time `json:"at"`
	CreatedAt           time.Time `json:"created_at"`
	ServerDeletedAt     time.Time `json:"server_deleted_at"`
	Color               string    `json:"color"`
	Billable            bool      `json:"billable"`
	Template            bool      `json:"template"`
	AutoEstimates       bool      `json:"auto_estimates"`
	EstimatedHours      float64   `json:"estimated_hours"`
	Rate                float64   `json:"rate"`
	RateLastUpdated     time.Time `json:"rate_last_updated"`
	Currency            string    `json:"currency"`
	Recurring           bool      `json:"recurring"`
	RecurringParameters struct {
		// Define fields for recurring parameters if needed
	} `json:"recurring_parameters"`
	FixedFee      float64 `json:"fixed_fee"`
	ActualHours   int     `json:"actual_hours"`
	ActualSeconds int     `json:"actual_seconds"`
	Wid           int     `json:"wid"`
	Cid           int     `json:"cid"`
}

type Client struct {
	ID       int       `json:"id"`
	Wid      int       `json:"wid"`
	Archived bool      `json:"archived"`
	Name     string    `json:"name"`
	At       time.Time `json:"at"`
}

const (
	togglApiBaseURL = "https://api.track.toggl.com/api/v9"
	outputFileName  = "time-entries.csv"
	defaultActivity = "Programmering"
)

func main() {
	// Load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	timeEntries := getTimeEntries()

	createCsv(timeEntries)

	fmt.Println("Done!")
}

func getTimeEntries() []TimeEntry {
	fromStr := os.Getenv("FROM")

	data := fetchToggl(togglApiBaseURL + "/me/time_entries?end_date=2030-01-01&start_date=" + fromStr)
	fmt.Println("Fetching time entries since " + fromStr + "...")

	var timeEntries []TimeEntry
	err := json.Unmarshal([]byte(data), &timeEntries)
	if err != nil {
		log.Fatal(err)
	}

	return timeEntries
}

func getProject(workspaceId int, projectId int) Project {
	url := togglApiBaseURL + "/workspaces/" + strconv.Itoa(workspaceId) + "/projects/" + strconv.Itoa(projectId)
	data := fetchToggl(url)

	var project Project
	err := json.Unmarshal([]byte(data), &project)
	if err != nil {
		log.Fatal(err)
	}

	return project
}

func getClient(workspaceId int, clientId int) Client {
	url := togglApiBaseURL + "/workspaces/" + strconv.Itoa(workspaceId) + "/clients/" + strconv.Itoa(clientId)
	data := fetchToggl(url)

	var client Client
	err := json.Unmarshal([]byte(data), &client)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func fetchToggl(url string) string {
	// Get env variables
	user := os.Getenv("USERNAME")
	pass := os.Getenv("PASSWORD")

	// Fetch time entries from Toggl
	req, err := http.NewRequest(http.MethodGet,
		url, nil)
	if err != nil {
		print(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(user, pass)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		print(err)
	}

	// Get the response body
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}

	return string(body)
}

func createCsv(timeEntries []TimeEntry) {
	// Create a new CSV file
	file, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Define the header row
	header := []string{"Datum", "Artikel", "Tid (timmar)", "Kund", "Projekt", "Aktivitet", "Ärendenummer", "Beskrivning"}

	// Write the header row to the CSV file
	if err := writer.Write(header); err != nil {
		panic(err)
	}

	// Write the data to the CSV file
	for _, record := range timeEntries {
		project := getProject(record.WorkspaceID, record.ProjectID)
		date := getDate(record.Start)
		duration := getDuration(record.Duration)
		activity := getActivity(record.Tags)
		ticketNumber := getTicketNumber(record.Description)
		description := getDescription(record.Description)

		if project.ClientID != 0 {
			client := getClient(record.WorkspaceID, project.ClientID)
			data := []string{date, "Normal", duration, client.Name, project.Name, activity, ticketNumber, description}
			if err := writer.Write(data); err != nil {
				panic(err)
			}
		}
	}

	fmt.Println("Writing to file " + outputFileName + "...")
}

func convertDatetoUnix(dateStr string) string {
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Fatal(err)
	}

	unixTimestamp := parsedDate.Unix()
	return strconv.FormatInt(unixTimestamp, 10)
}

func getDate(date time.Time) string {
	return date.Format("2006-01-02")
}

func getDuration(duration int) string {
	rounded := math.Round(float64(duration) / 60 / 60)

	return strconv.Itoa(int(rounded))
}

func getActivity(tags []string) string {
	if len(tags) > 0 {
		tag := tags[0]
		return tag
	}
	return defaultActivity
}

func getTicketNumber(description string) string {
	substrings := strings.Split(description, " | ")

	if len(substrings) > 1 {
		return substrings[0]
	}

	return ""
}

func getDescription(description string) string {
	substrings := strings.Split(description, " | ")

	if len(substrings) > 1 {
		return substrings[1]
	}

	return description
}
