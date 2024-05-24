package model

// FxRates struct to hold the exchange rates.
type FxRates struct {
	Rates Rates `json:"rates"`
}

// Rates struct to hold the PHP exchange rate.
type Rates struct {
	PHP float64 `json:"PHP"`
}

// User struct to hold the default workspace.
type User struct {
	Id               string `json:"id"`
	DefaultWorkspace string `json:"defaultWorkspace"`
}

// SummaryReport struct to hold the summary report payload.
type SummaryReport struct {
	DateRangeStart string        `json:"dateRangeStart"`
	DateRangeEnd   string        `json:"dateRangeEnd"`
	SortOrder      string        `json:"sortOrder"`
	ExportType     string        `json:"exportType"`
	AmountShown    string        `json:"amountShown"`
	SummaryFilter  SummaryFilter `json:"summaryFilter"`
	Users          Users         `json:"users"`
}

// SummaryFilter struct to hold the summary filter.
type SummaryFilter struct {
	Groups []string `json:"groups"`
}

// Users struct to hold the user id.
type Users struct {
	Contains string   `json:"contains"`
	Ids      []string `json:"ids"`
	Status   string   `json:"status"`
}

// OutputSummary struct to hold the response from the summary report API.
type OutputSummary struct {
	Total    []Total    `json:"totals"`
	GroupOne []GroupOne `json:"groupOne"`
}

// Totals struct for the OutputSummary struct which holds the total work hours.
type Total struct {
	TotalTime int `json:"totalTime"`
}

// GroupOne struct for the OutputSummary struct which holds the date and total work hours.
type GroupOne struct {
	Duration int    `json:"duration"`
	Name     string `json:"name"`
}
