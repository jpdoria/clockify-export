package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jpdoria/clockify-export/handler"
)

var (
	start, end  = getCurrentDateRange()
	customRange = flag.String("customRange", fmt.Sprintf("%v to %v", start, end), "Custom range using this format: \"1970-01-01T00:00:00.000Z to 1970-01-31T23:59:59.999Z\"")
	version     = flag.Bool("version", false, "Print the current version.")
	ver, build  string
)

// This function checks if a flag is passed.
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

// Function that gets the current date range.
func getCurrentDateRange() (string, string) {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	lastOfMonth = lastOfMonth.Add(time.Hour * 23)
	lastOfMonth = lastOfMonth.Add(time.Minute * 59)
	lastOfMonth = lastOfMonth.Add(time.Second * 59)
	lastOfMonth = lastOfMonth.Add(time.Millisecond * 999)
	layout := "2006-01-02T15:04:05.999Z"

	return firstOfMonth.Format(layout), lastOfMonth.Format(layout)
}

func main() {
	flag.Parse()

	if isFlagPassed("customRange") {
		start = strings.Split(*customRange, " to ")[0]
		end = strings.Split(*customRange, " to ")[1]
	}

	if *version {
		fmt.Printf("Version: %v.%v\n", ver, build)
		os.Exit(0)
	}

	userId, workspaceId := handler.ClockifyGetWorkspace()
	handler.ClockifyGetWorkHoursGroupByDate(userId, workspaceId, start, end)
	invoice := handler.ClockifyGetWorkHoursGroupByProject(userId, workspaceId)
	fmt.Printf("Total Earnings: $%.2f (₱%.2f)\n", invoice.SubTotal, handler.GetExchangeRates(invoice.SubTotal))
	handler.CreateSpreadsheet(invoice)
}
