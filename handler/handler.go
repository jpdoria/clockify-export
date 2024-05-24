package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/jpdoria/clockify-export/model"
)

// Constants for the API URLs.
const (
	clockifyApiUrl        = "https://api.clockify.me/api"
	clockifyReportsApiUrl = "https://reports.api.clockify.me"
	fxRatesApiUrl         = "https://api.fxratesapi.com/latest"
)

// Variables for the http client, clockify API key, hourly rate in USD, and the payload for the summary report.
var (
	client         = &http.Client{}
	clockifyApiKey = os.Getenv("CLOCKIFY_API_KEY")
	hourlyRateUsd  = os.Getenv("HOURLY_RATE_USD")
	payload        = &model.SummaryReport{
		DateRangeStart: "2024-05-01T00:00:00.000Z",
		DateRangeEnd:   "2024-05-30T23:59:59.999Z",
		SortOrder:      "ASCENDING",
		ExportType:     "JSON",
		AmountShown:    "HIDE_AMOUNT",
		SummaryFilter: model.SummaryFilter{
			Groups: []string{"placeholder", "USER", "TIMEENTRY"},
		},
		Users: model.Users{
			Contains: "CONTAINS",
			Ids:      []string{"placeholder"},
			Status:   "ALL",
		},
	}
)

// Check if the required environment variables are set.
func init() {
	if clockifyApiKey == "" {
		fmt.Println("CLOCKIFY_API_KEY is not set")
		os.Exit(1)
	}

	if hourlyRateUsd == "" {
		fmt.Println("HOURLY_RATE_USD is not set")
		os.Exit(1)
	}
}

// callSummaryReportAPI calls the summary report API of Clockify.
func callSummaryReportAPI(workspaceId string, payloadBuffer *bytes.Buffer) []byte {
	req, err := http.NewRequest("POST", clockifyReportsApiUrl+"/v1/workspaces/"+workspaceId+"/reports/summary", payloadBuffer)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	req.Header.Set("X-Api-Key", clockifyApiKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer res.Body.Close()

	r, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v", err)
	}

	return r
}

// getExchangeRates fetches the latest exchange rates from the fxRatesApiUrl.
func GetExchangeRates(usdCurrency float64) (phpCurrency float64) {
	resp, err := http.Get(fxRatesApiUrl)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v", err)
	}

	var fxRates model.FxRates
	json.Unmarshal(r, &fxRates)
	fmt.Printf("\nExchange rate right now for 1 USD to PHP: %.2f\n", fxRates.Rates.PHP)
	return usdCurrency * fxRates.Rates.PHP
}

// CalculateEarnings calculates the earnings of the user based on the hours worked.
func CalculateEarnings(hours float64) float64 {
	hr, err := strconv.ParseFloat(hourlyRateUsd, 64)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return hours * hr
}

// func convertTimetoHHMMSS output to decimal format.
func convertTimeToDecimal(seconds int) float64 {
	return float64(seconds) / 3600
}

// convertTimetoHHMMSS converts the seconds to HH:MM:SS format.
func convertTimetoHHMMSS(seconds int) string {
	hours := seconds / 3600
	seconds %= 3600
	minutes := seconds / 60
	seconds %= 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

// ClockifyGetWorkHoursGroupByDate fetches the work hours of the user grouped by date.
func ClockifyGetWorkHoursGroupByDate(userId, workspaceId string) {
	payload.SummaryFilter.Groups[0] = "DATE"
	payload.Users.Ids[0] = userId
	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(payload)
	res := callSummaryReportAPI(workspaceId, payloadBuffer)

	var outputSummary model.OutputSummary
	json.Unmarshal(res, &outputSummary)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	day := 0
	fmt.Println("Work Log:")
	fmt.Fprintln(w, "ID\tDATE\tHOURS\tEARNINGS")
	for _, groupOne := range outputSummary.GroupOne {
		day += 1
		msg := fmt.Sprintf("%v\t%v\t%v (%.2f)\t$%.2f", day, groupOne.Name, convertTimetoHHMMSS(groupOne.Duration), convertTimeToDecimal(groupOne.Duration), CalculateEarnings(convertTimeToDecimal(groupOne.Duration)))
		fmt.Fprintln(w, msg)
	}
	w.Flush()
}

// ClockifyGetWorkHoursGroupByProject fetches the work hours of the user grouped by project.
func ClockifyGetWorkHoursGroupByProject(userId, workspaceId string) float64 {
	payload.SummaryFilter.Groups[0] = "PROJECT"
	payload.Users.Ids[0] = userId
	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(payload)
	res := callSummaryReportAPI(workspaceId, payloadBuffer)

	var outputSummary model.OutputSummary
	json.Unmarshal(res, &outputSummary)

	fmt.Printf("Total Hours: %v (%.2f)\n", convertTimetoHHMMSS(outputSummary.Total[0].TotalTime), convertTimeToDecimal(outputSummary.Total[0].TotalTime))
	return convertTimeToDecimal(outputSummary.Total[0].TotalTime)
}

// clockifyGetWorkspace fetches the default workspace and user id of the user.
func ClockifyGetWorkspace() (string, string) {
	req, err := http.NewRequest("GET", clockifyApiUrl+"/v1/user", nil)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	req.Header.Set("X-Api-Key", clockifyApiKey)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer res.Body.Close()

	var user model.User
	r, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v", err)
	}
	json.Unmarshal(r, &user)
	return user.Id, user.DefaultWorkspace
}
