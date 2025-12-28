package transaction

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"peronal_finance_cli_manager/internal/models"
)

// ParseCSV parses a CSV file into a slice of Transactions.
// CSV format: Category,Amount,Date
func ParseCSV(filePath string) ([]models.Transaction, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}
	if len(header) < 3 {
		return nil, errors.New("CSV must have at least 3 columns: Category, Amount, Date")
	}

	var transactions []models.Transaction

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(record) < 3 {
			continue
		}

		category := record[0]

		amount, err := strconv.ParseFloat(record[1], 32)
		if err != nil {
			return nil, errors.New("invalid amount in CSV: " + record[1])
		}

		date, err := time.Parse("2006-01-02 15:04:05-07:00", record[2])
		if err != nil {
			return nil, errors.New("invalid date in CSV: " + record[2])
		}

		tx := models.Transaction{
			Category: models.Category{Name: category},
			Amount:   float32(amount),
			Date:     date,
		}

		transactions = append(transactions, tx)
	}

	return transactions, nil
}

// DetectFormat detects the file format based on extension
func DetectFormat(filePath string) string {
	if strings.HasSuffix(strings.ToLower(filePath), ".csv") {
		return "csv"
	}
	return ""
}
