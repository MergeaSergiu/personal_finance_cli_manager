package db

import (
	"peronal_finance_cli_manager/internal/models"
	"regexp"
)

var categoryRules = []models.CategoryRule{
	{regexp.MustCompile(`(?i)gas|electric`), "Bills"},
	{regexp.MustCompile(`(?i)uber|bolt|taxi`), "Transport"},
	{regexp.MustCompile(`(?i)netflix|spotify`), "Subscriptions"},
}

// RecommendCategory Returns empty string if no match
func RecommendCategory(description string) string {
	for _, rule := range categoryRules {
		if rule.Pattern.MatchString(description) {
			return rule.Category
		}
	}
	return ""
}
