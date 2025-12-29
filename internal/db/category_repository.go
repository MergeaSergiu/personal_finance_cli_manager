package db

import (
	"peronal_finance_cli_manager/internal/models"
)

func GetBudgetStats() ([]models.BudgetStats, error) {
	rows, err := DB.Raw(`
		SELECT 
			c.name,
			c.budget,
			COALESCE(SUM(t.amount), 0) as spent
		FROM categories c
		LEFT JOIN transactions t ON t.category_id = c.id
		GROUP BY c.id
	`).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := []models.BudgetStats{}

	for rows.Next() {
		var s models.BudgetStats
		if err := rows.Scan(&s.CategoryName, &s.Budget, &s.Spent); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}

	return stats, nil
}
