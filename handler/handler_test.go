package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func init() {
	clockifyApiKey = "Str0ng3stKey3v3r"
	hourlyRateUsd = "10.0"
}

// Helper function to set up environment variables for testing.
func setUpEnv() {
	_ = os.Setenv("CLOCKIFY_API_KEY", "Str0ng3stKey3v3r")
	_ = os.Setenv("HOURLY_RATE_USD", "10.0")
}

// Helper function to tear down environment variables after testing.
func tearDownEnv() {
	_ = os.Unsetenv("CLOCKIFY_API_KEY")
	_ = os.Unsetenv("HOURLY_RATE_USD")
}

func TestCallSummaryReportAPI(t *testing.T) {
	setUpEnv()
	defer tearDownEnv()

	// Mock server for Clockify Reports API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("\nExpected method POST\nGot %s", r.Method)
		}
		if r.Header.Get("X-Api-Key") != clockifyApiKey {
			t.Errorf("\nExpected API key %s\nGot %s", clockifyApiKey, r.Header.Get("X-Api-Key"))
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
	"totals": [
		{
			"totalTime": 28800,
			"totalBillableTime": 0,
			"entriesCount": 1,
			"amounts": [],
			"numOfCurrencies": 1,
			"_id": "",
			"totalAmount": null,
			"totalAmountByCurrency": {}
		}
	],
	"groupOne": [
		{
			"duration": 28800,
			"_id": "2024-05-01",
			"name": "2024-05-01",
			"children": [
				{
					"duration": 28800,
					"hasMembership": true,
					"_id": "e6205drm558ps2y2pscu6hh1",
					"name": "Foo Bar Baz",
					"nameLowerCase": "foo bar baz",
					"sortingField": "1-name",
					"children": [
						{
							"duration": 28800,
							"_id": {
								"DATE": "2024-05-01",
								"USER": "e6205drm558ps2y2pscu6hh1",
								"TIMEENTRY": "[CARD-441]: May Overhead - FBB"
							},
							"name": "[CARD-441]: May Overhead - FBB",
							"amounts": []
						}
					],
					"amounts": []
				}
			],
			"amounts": []
		}
	]
}`))
	}))
	defer mockServer.Close()

	// Override the Clockify Reports API URL with the mock server URL
	clockifyReportsApiUrl = mockServer.URL

	// Create a payload buffer
	payloadBuffer := new(bytes.Buffer)
	_ = json.NewEncoder(payloadBuffer).Encode(payload)

	// Call the function
	response := callSummaryReportAPI("workspace-id", payloadBuffer)

	// Verify the response
	expectedResponse := `{
	"totals": [
		{
			"totalTime": 28800,
			"totalBillableTime": 0,
			"entriesCount": 1,
			"amounts": [],
			"numOfCurrencies": 1,
			"_id": "",
			"totalAmount": null,
			"totalAmountByCurrency": {}
		}
	],
	"groupOne": [
		{
			"duration": 28800,
			"_id": "2024-05-01",
			"name": "2024-05-01",
			"children": [
				{
					"duration": 28800,
					"hasMembership": true,
					"_id": "e6205drm558ps2y2pscu6hh1",
					"name": "Foo Bar Baz",
					"nameLowerCase": "foo bar baz",
					"sortingField": "1-name",
					"children": [
						{
							"duration": 28800,
							"_id": {
								"DATE": "2024-05-01",
								"USER": "e6205drm558ps2y2pscu6hh1",
								"TIMEENTRY": "[CARD-441]: May Overhead - FBB"
							},
							"name": "[CARD-441]: May Overhead - FBB",
							"amounts": []
						}
					],
					"amounts": []
				}
			],
			"amounts": []
		}
	]
}`
	if string(response) != expectedResponse {
		t.Errorf("\nExpected %s\nGot %s", expectedResponse, string(response))
	}
}

func TestGetExchangeRates(t *testing.T) {
	setUpEnv()
	defer tearDownEnv()

	// Mock server for FX Rates API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("\nExpected method GET\nGot %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"rates": {
				"PHP": 50.0
			}
		}`))
	}))
	defer mockServer.Close()

	// Override the FX Rates API URL with the mock server URL
	fxRatesApiUrl = mockServer.URL

	// Call the function
	exchangeRate := GetExchangeRates(1.0)

	// Verify the response
	expectedRate := 50.0
	if exchangeRate != expectedRate {
		t.Errorf("\nExpected %.2f\nGot %.2f", expectedRate, exchangeRate)
	}
}

func TestClockifyGetWorkspace(t *testing.T) {
	setUpEnv()
	defer tearDownEnv()

	// Mock server for Clockify API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("\nExpected method GET\nGot %s", r.Method)
		}
		if r.Header.Get("X-Api-Key") != clockifyApiKey {
			t.Errorf("\nExpected API key %s\nGot %s", clockifyApiKey, r.Header.Get("X-Api-Key"))
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": "user-id",
			"defaultWorkspace": "workspace-id"
		}`))
	}))
	defer mockServer.Close()

	// Override the Clockify API URL with the mock server URL
	clockifyApiUrl = mockServer.URL

	// Call the function
	userId, workspaceId := ClockifyGetWorkspace()

	// Verify the response
	expectedUserId := "user-id"
	expectedWorkspaceId := "workspace-id"
	if userId != expectedUserId || workspaceId != expectedWorkspaceId {
		t.Errorf("\nExpected userId %s and workspaceId %s\nGot userId %s and workspaceId %s", expectedUserId, expectedWorkspaceId, userId, workspaceId)
	}
}

func TestClockifyGetWorkHoursGroupByDate(t *testing.T) {
	setUpEnv()
	defer tearDownEnv()

	// Mock server for Clockify Reports API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("\nExpected method POST\nGot %s", r.Method)
		}
		if r.Header.Get("X-Api-Key") != clockifyApiKey {
			t.Errorf("\nExpected API key %s\nGot %s", clockifyApiKey, r.Header.Get("X-Api-Key"))
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"totals": [{
				"totalTime": 1032750,
				"totalBillableTime": 0,
				"entriesCount": 23,
				"amounts": [],
				"numOfCurrencies": 1,
				"_id": "",
				"totalAmount": null,
				"totalAmountByCurrency": {}
			}],
			"groupOne": [{
				"duration": 28800,
				"_id": "2024-05-01",
				"name": "2024-05-01",
				"children": [{
					"duration": 28800,
					"hasMembership": true,
					"_id": "e6205drm558ps2y2pscu6hh1",
					"name": "Foo Bar Baz",
					"nameLowerCase": "foo bar baz",
					"sortingField": "1-name",
					"children": [{
						"duration": 28800,
						"_id": {
							"DATE": "2024-05-01",
							"USER": "e6205drm558ps2y2pscu6hh1",
							"TIMEENTRY": "[CARD-441]: May Overhead - FBB"
						},
						"name": "[CARD-441]: May Overhead - FBB",
						"amounts": []
					}],
					"amounts": []
				}],
				"amounts": []
			}]
		}`))
	}))
	defer mockServer.Close()

	// Override the Clockify Reports API URL with the mock server URL
	clockifyReportsApiUrl = mockServer.URL

	// Mock the standard output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the function
	ClockifyGetWorkHoursGroupByDate("user-id", "workspace-id", "2024-05-01T00:00:00.000Z", "2024-05-31T23:59:59.999Z")

	// Capture the output
	_ = w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)

	// Verify the output
	expectedOutput := "Work Log:\nID     DATE           HOURS               EARNINGS\n1      2024-05-01     08:00:00 (8.00)     $80.00\n"
	if buf.String() != expectedOutput {
		t.Errorf("\nExpected output %s\nGot %s", expectedOutput, buf.String())
	}
}

func TestClockifyGetWorkHoursGroupByProject(t *testing.T) {
	setUpEnv()
	defer tearDownEnv()

	// Mock server for Clockify Reports API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("\nExpected method POST\nGot %s", r.Method)
		}
		if r.Header.Get("X-Api-Key") != clockifyApiKey {
			t.Errorf("\nExpected API key %s\nGot %s", clockifyApiKey, r.Header.Get("X-Api-Key"))
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"totals": [{
				"totalTime": 28800,
				"totalBillableTime": 0,
				"entriesCount": 1,
				"amounts": [],
				"numOfCurrencies": 1,
				"_id": "",
				"totalAmount": null,
				"totalAmountByCurrency": {}
			}],
			"groupOne": [{
				"duration": 28800,
				"_id": "2024-05-01",
				"name": "Project A",
				"children": [{
					"duration": 28800,
					"hasMembership": true,
					"_id": "user-id",
					"name": "Task A",
					"nameLowerCase": "task a",
					"sortingField": "1-name",
					"children": [{
						"duration": 28800,
						"_id": {
							"DATE": "2024-05-01",
							"USER": "user-id",
							"TIMEENTRY": "[CARD-441]: May Overhead - Task A"
						},
						"name": "[CARD-441]: May Overhead - Task A",
						"amounts": []
					}],
					"amounts": []
				}],
				"amounts": []
			}]
		}`))
	}))
	defer mockServer.Close()

	// Override the Clockify Reports API URL with the mock server URL
	clockifyReportsApiUrl = mockServer.URL

	// Call the function
	invoice := ClockifyGetWorkHoursGroupByProject("user-id", "workspace-id")

	// Verify the invoice details
	expectedSubTotal := 80.0
	expectedPayoneerFee := 2.48
	expectedGrandTotal := 82.48
	if invoice.SubTotal != expectedSubTotal || invoice.PayoneerFee != expectedPayoneerFee || invoice.GrandTotal != expectedGrandTotal {
		t.Errorf("\nExpected SubTotal %.2f, PayoneerFee %.2f, GrandTotal %.2f\nGot SubTotal %.2f, PayoneerFee %.2f, GrandTotal %.2f",
			expectedSubTotal, expectedPayoneerFee, expectedGrandTotal, invoice.SubTotal, invoice.PayoneerFee, invoice.GrandTotal)
	}
}
