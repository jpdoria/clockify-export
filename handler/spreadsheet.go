package handler

import (
	"fmt"
	"os"

	"github.com/jpdoria/clockify-export/model"
	"github.com/xuri/excelize/v2"
)

// CreateSpreadsheet creates a new spreadsheet.
func CreateSpreadsheet(invoice *model.Invoice) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Set value of a cell.
	f.SetCellValue("Sheet1", "D1", "Date:")
	f.SetCellValue("Sheet1", "E1", invoice.Date)
	f.SetCellValue("Sheet1", "D2", "Invoice ID:")
	f.SetCellValue("Sheet1", "E2", invoice.Id)
	f.SetCellValue("Sheet1", "A4", "From:")
	f.SetCellValue("Sheet1", "A5", "Foo Bar\nDeveloper")
	f.SetCellValue("Sheet1", "A7", "To:")
	f.SetCellValue("Sheet1", "A8", "Client\nClient Position")
	f.SetCellValue("Sheet1", "E10", "Rate Per Hour:")
	f.SetCellValue("Sheet1", "E11", fmt.Sprintf("$%.2f", invoice.HourlyRate))

	// Add the column headers.
	columnHeaders := []interface{}{"ID", "Description", "Date", "Hours", "Amount"}
	f.SetSheetRow("Sheet1", "A13", &columnHeaders)

	// Add the work log entries just below the column headers.
	const headerRow = 13
	for idx, entry := range invoice.WorkLog {
		cell, err := excelize.CoordinatesToCellName(1, headerRow+idx+1)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetSheetRow("Sheet1", cell, &[]interface{}{entry.Id, entry.Description, entry.Date, entry.Hours, entry.Amount})
	}

	// Add the subtotal, Payoneer fee, and grand total dynamically after work log entries.
	summaryStartRow := headerRow + len(invoice.WorkLog) + 2
	f.SetCellValue("Sheet1", fmt.Sprintf("D%d", summaryStartRow), "Subtotal")
	f.SetCellValue("Sheet1", fmt.Sprintf("E%d", summaryStartRow), fmt.Sprintf("$%.2f", invoice.SubTotal))
	f.SetCellValue("Sheet1", fmt.Sprintf("D%d", summaryStartRow+1), fmt.Sprintf("Payoneer Fee (%.1f%%)", payoneerFeeRate*100))
	f.SetCellValue("Sheet1", fmt.Sprintf("E%d", summaryStartRow+1), fmt.Sprintf("$%.2f", invoice.PayoneerFee))
	f.SetCellValue("Sheet1", fmt.Sprintf("D%d", summaryStartRow+2), "Grand Total")
	f.SetCellValue("Sheet1", fmt.Sprintf("E%d", summaryStartRow+2), fmt.Sprintf("$%.2f", invoice.GrandTotal))

	// Save spreadsheet by the given path.
	if err := os.MkdirAll("out", 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}
	invoiceFileName := fmt.Sprintf("out/invoice-%s.xlsx", invoice.Date)
	if err := f.SaveAs(invoiceFileName); err != nil {
		fmt.Printf("Error saving invoice: %v\n", err)
		return
	}
	fmt.Printf("Successfully created an invoice: %s\n", invoiceFileName)
}
