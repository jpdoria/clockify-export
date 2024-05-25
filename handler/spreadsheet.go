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
	f.SetCellValue("Sheet1", "E1", len(invoice.WorkLog)-1)
	f.SetCellValue("Sheet1", "D2", "Invoice No.:")
	f.SetCellValue("Sheet1", "E2", "1")
	f.SetCellValue("Sheet1", "A4", "From:")
	f.SetCellValue("Sheet1", "A5", "Foo Bar\nDeveloper")
	f.SetCellValue("Sheet1", "A7", "To:")
	f.SetCellValue("Sheet1", "A8", "Client\nClient Position")
	f.SetCellValue("Sheet1", "E10", "Rate Per Hour:")
	f.SetCellValue("Sheet1", "E11", fmt.Sprintf("$%.2f", invoice.HourlyRate))

	// Add the column headers.
	for _, columnHeaders := range [][]interface{}{{"ID", "Description", "Date", "Hours", "Amount"}} {
		f.SetSheetRow("Sheet1", "A13", &columnHeaders)
	}

	// Add the work log entries just below the column headers.
	for idx, entry := range invoice.WorkLog {
		cell, err := excelize.CoordinatesToCellName(1, idx+14)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetSheetRow("Sheet1", cell, &[]interface{}{entry.Id, entry.Description, entry.Date, entry.Hours, entry.Amount})
	}

	// Add the subtotal, Payoneer fee, and grand total at the bottom.
	f.SetCellValue("Sheet1", "D33", "Subtotal")
	f.SetCellValue("Sheet1", "E33", fmt.Sprintf("$%.2f", invoice.SubTotal))
	f.SetCellValue("Sheet1", "D34", "Payoneer Fee (3.1%)")
	f.SetCellValue("Sheet1", "E34", fmt.Sprintf("$%.2f", invoice.PayoneerFee))
	f.SetCellValue("Sheet1", "D35", "Grand Total")
	f.SetCellValue("Sheet1", "E35", fmt.Sprintf("$%.2f", invoice.GrandTotal))

	// Save spreadsheet by the given path.
	err := os.Mkdir("out", 0755)
	if err != nil {
		fmt.Printf("Directory already exists: %v\n", err)
	}
	if err := f.SaveAs("out/invoice.xlsx"); err != nil {
		fmt.Println(err)
	}
}
