package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"os/user"
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
	currentAppVersion     = "v1.2.2"
	togglApiBaseURL       = "https://api.track.toggl.com/api/v9"
	defaultOutputFileName = "time-entries.csv"
	defaultActivity       = "Programmering"
	configDir             = "/.config/ttsync"
)

var startDate string
var endDate string
var outputFileName string
var showVersion bool

func main() {
	// Load CLI arguments
	loadAppArgs()

	// Display current version
	if showVersion {
		fmt.Println(currentAppVersion)
		return
	}

	// Construct the full path to the .env file
	envFilePath := getCurrentUserHomeDir() + configDir + "/.env"

	// Load env file
	err := godotenv.Load(envFilePath)
	if err != nil {
		fmt.Printf("Error %s\n", err)
		fmt.Println("Trying to load local .env instead...")

		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading local .env")
		}
	}

	// Get time entries from Toggl
	timeEntries := getTimeEntries()

	// Format and write time entries to a .csv file
	createCsv(timeEntries)

	fmt.Println("Done!")
}

func getTimeEntries() []TimeEntry {
	url := togglApiBaseURL + "/me/time_entries?end_date=" + endDate + "&start_date=" + startDate
	data := fetchToggl(url)
	fmt.Println("Fetching time entries between " + startDate + " and " + endDate + "...")

	var timeEntries []TimeEntry
	err := json.Unmarshal([]byte(data), &timeEntries)
	if err != nil {
		log.Println("URL: ", url)
		log.Fatal("Failed time entry fetch: ", err)
	}

	time.Sleep(100 * time.Millisecond)

	return timeEntries
}

func getProject(workspaceId int, projectId int) Project {
	url := togglApiBaseURL + "/workspaces/" + strconv.Itoa(workspaceId) + "/projects/" + strconv.Itoa(projectId)
	data := fetchToggl(url)

	var project Project
	err := json.Unmarshal([]byte(data), &project)
	if err != nil {
		log.Println("URL: ", url)
		log.Fatal("Failed project fetch: ", err)
	}

	time.Sleep(100 * time.Millisecond)

	return project
}

func getClient(workspaceId int, clientId int) Client {
	url := togglApiBaseURL + "/workspaces/" + strconv.Itoa(workspaceId) + "/clients/" + strconv.Itoa(clientId)
	data := fetchToggl(url)

	var client Client
	err := json.Unmarshal([]byte(data), &client)
	if err != nil {
		log.Println("URL: ", url)
		log.Fatal("Failed client fetch: ", err)
	}

	time.Sleep(100 * time.Millisecond)

	return client
}

func fetchToggl(url string) string {
	// Get env variables
	user := os.Getenv("USERNAME")
	pass := os.Getenv("PASSWORD")

	if user == "" || pass == "" {
		log.Fatal("Username or password variables are empty or does not exist in the supplied .env file")
	}

	// Create request object
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(user, pass)

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// Get the response body
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
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
	header := []string{
		"Datum", "Artikel", "Tid (timmar)",
		"Kund", "Projekt", "Aktivitet", "Ärendenummer", "Beskrivning"}

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

func getCurrentUserHomeDir() string {
	// Get the current user
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal("Could not find current users home directory")
	}

	return currentUser.HomeDir
}

func loadAppArgs() {
	// Define command-line flags
	flag.StringVar(&startDate, "start", getPrevousMonday(), "The start date from which we should fetch time entries")
	flag.StringVar(&startDate, "s", getPrevousMonday(), "Alias for start date")
	flag.StringVar(&endDate, "end", "2100-01-01", "The end date from which we should fetch time entries")
	flag.StringVar(&endDate, "e", "2100-01-01", "Alias for end date")
	flag.StringVar(&outputFileName, "output", defaultOutputFileName, "The output path where the CSV file should be saved")
	flag.StringVar(&outputFileName, "o", defaultOutputFileName, "Alias for output path")
	flag.BoolVar(&showVersion, "version", false, "If the user wants to display the current version")
	flag.BoolVar(&showVersion, "v", false, "Alias for current version")

	// Parse the command-line arguments
	flag.Parse()
}

func getPrevousMonday() string {
	// Get the current date
	today := time.Now()

	// Calculate the number of days to subtract to reach the previous Monday
	daysToSubtract := int(today.Weekday() - time.Monday)
	if daysToSubtract < 0 {
		daysToSubtract += 7 // Wrap around to the previous week
	}

	// Calculate the date of the previous Monday
	previousMonday := today.AddDate(0, 0, -daysToSubtract)

	return previousMonday.Format("2006-01-02")
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
	hours := float64(duration) / 3600.0
	roundedHours := math.Ceil(hours*2) / 2

	// If roundedHours is an integer, return it as is, otherwise format to one decimal point
	if roundedHours == math.Floor(roundedHours) {
		return fmt.Sprintf("%.0f", roundedHours)
	}
	return fmt.Sprintf("%.1f", roundedHours)
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
