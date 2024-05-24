package main

import (
	"fmt"

	"github.com/jpdoria/clockify-export/handler"
)

func main() {
	userId, workspaceId := handler.ClockifyGetWorkspace()
	handler.ClockifyGetWorkHoursGroupByDate(userId, workspaceId)
	decimalHours := handler.ClockifyGetWorkHoursGroupByProject(userId, workspaceId)
	fmt.Printf("Total Earnings: $%.2f (₱%.2f)\n", handler.CalculateEarnings(decimalHours), handler.GetExchangeRates(handler.CalculateEarnings(decimalHours)))
}
