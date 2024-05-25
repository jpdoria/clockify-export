package main

import (
	"fmt"

	"github.com/jpdoria/clockify-export/handler"
)

func main() {
	userId, workspaceId := handler.ClockifyGetWorkspace()
	handler.ClockifyGetWorkHoursGroupByDate(userId, workspaceId)
	invoice := handler.ClockifyGetWorkHoursGroupByProject(userId, workspaceId)
	fmt.Printf("Total Earnings: $%.2f (₱%.2f)\n", invoice.SubTotal, handler.GetExchangeRates(invoice.SubTotal))
	handler.CreateSpreadsheet(invoice)
}
