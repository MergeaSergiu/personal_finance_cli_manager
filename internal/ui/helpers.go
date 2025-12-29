package ui

func calculatePercentage(spent, budget float32) float32 {
	if budget == 0 {
		return 0
	}
	return (spent / budget) * 100
}
