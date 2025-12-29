package ui

import (
	"fmt"
	"strings"
)

func calculatePercentage(spent, budget float32) float32 {
	if budget == 0 {
		return 0
	}
	return (spent / budget) * 100
}

func generateMonthlyExpenseChart(categoryTotals map[string]float32) string {
	report := "📄 Monthly Expense Report\n\n"
	if len(categoryTotals) == 0 {
		return report + "No expenses found for this month.\n"
	}

	// find max for scaling bars
	var max float32
	for _, amt := range categoryTotals {
		if amt > max {
			max = amt
		}
	}

	barWidth := 40
	for cat, amt := range categoryTotals {
		length := int((amt / max) * float32(barWidth))
		if length < 1 {
			length = 1
		}
		bar := strings.Repeat("▇", length)
		report += fmt.Sprintf("%-15s %s (%.2f)\n", cat, bar, amt)
	}

	return report
}
