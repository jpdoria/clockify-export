package handler

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/jpdoria/clockify-export/model"
	"github.com/xuri/excelize/v2"
)

func TestCreateSpreadsheet(t *testing.T) {
	// Set up test data
	invoice := &model.Invoice{
		Date:        "2024-05-01",
		Id:          "foobarbaz",
		HourlyRate:  10.0,
		SubTotal:    80.0,
		PayoneerFee: 2.4,
		GrandTotal:  82.4,
		WorkLog: []model.WorkLog{
			{
				Id:          1,
				Description: "Independent Contractor",
				Date:        "2024-05-01",
				Hours:       "8",
				Amount:      "80.0",
			},
		},
	}

	// Call the function to create the spreadsheet
	CreateSpreadsheet(invoice)

	// Verify the file is created
	invoiceFileName := fmt.Sprintf("out/invoice-%s.xlsx", invoice.Date)
	if _, err := os.Stat(invoiceFileName); os.IsNotExist(err) {
		t.Fatalf("Expected file %s to be created", invoiceFileName)
	}
	// defer os.Remove(invoiceFileName) // Clean up

	// Open the spreadsheet
	f, err := excelize.OpenFile(invoiceFileName)
	if err != nil {
		t.Fatalf("Failed to open created spreadsheet: %v", err)
	}
	defer f.Close()

	// Verify the contents of the spreadsheet
	assertCellValue := func(t *testing.T, sheet, axis, expected string) {
		t.Helper()
		value, err := f.GetCellValue(sheet, axis)
		if err != nil {
			t.Fatalf("Failed to get cell value at %s: %v", axis, err)
		}
		if value != expected {
			t.Fatalf("Expected cell value at %s to be %s, got %s", axis, expected, value)
		}
	}

	assertCellValue(t, "Sheet1", "E1", invoice.Date)
	assertCellValue(t, "Sheet1", "E2", invoice.Id)
	assertCellValue(t, "Sheet1", "E11", fmt.Sprintf("$%.2f", invoice.HourlyRate))
	assertCellValue(t, "Sheet1", "D38", "Subtotal")
	assertCellValue(t, "Sheet1", "E38", fmt.Sprintf("$%.2f", invoice.SubTotal))
	assertCellValue(t, "Sheet1", "D39", "Payoneer Fee (3.1%)")
	assertCellValue(t, "Sheet1", "E39", fmt.Sprintf("$%.2f", invoice.PayoneerFee))
	assertCellValue(t, "Sheet1", "D40", "Grand Total")
	assertCellValue(t, "Sheet1", "E40", fmt.Sprintf("$%.2f", invoice.GrandTotal))

	// Verify the work log entry
	assertCellValue(t, "Sheet1", "A14", strconv.Itoa(invoice.WorkLog[0].Id))
	assertCellValue(t, "Sheet1", "B14", invoice.WorkLog[0].Description)
	assertCellValue(t, "Sheet1", "C14", invoice.WorkLog[0].Date)
	assertCellValue(t, "Sheet1", "D14", invoice.WorkLog[0].Hours)
	assertCellValue(t, "Sheet1", "E14", invoice.WorkLog[0].Amount)
}
