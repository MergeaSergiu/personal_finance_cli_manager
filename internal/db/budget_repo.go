package db

import (
	"fmt"
	"peronal_finance_cli_manager/internal/models"

	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

func CheckBudget(DB *gorm.DB, category models.Category, amount float32, date string) error {
	var total float32
	err := DB.Model(&models.Transaction{}).
		Where("category_id = ?", category.ID).
		Select("SUM(amount)").
		Row().Scan(&total)
	if err != nil {
		return err
	}

	if total > category.Budget {
		return PublishBudgetAlert(category.Name, total, category.Budget, amount, date)
	}

	return nil
}

func PublishBudgetAlert(category string, total, budget float32, latestAmount float32, latestDate string) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {

		}
	}(ch)

	q, err := ch.QueueDeclare(
		"budget_alerts",
		true, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf(
		"⚠️ Budget exceeded for %s\nTotal spent: %.2f / %.2f\nLatest transaction: %.2f on %s\nPlease review your spending!",
		category,
		total,
		budget,
		latestAmount,
		latestDate,
	)
	return ch.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msg),
	})
}
