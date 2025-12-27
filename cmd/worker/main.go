package main

import (
	"fmt"

	"github.com/streadway/amqp"
	"gopkg.in/gomail.v2"
)

func main() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {

		}
	}(ch)

	q, err := ch.QueueDeclare("budget_alerts", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Worker listening for budget alerts...")

	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			fmt.Println("Received alert:", string(msg.Body))
			sendEmail("user@example.com", "Budget Alert", string(msg.Body))
		}
	}()

	<-forever
}

func sendEmail(to, subject, body string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "no-reply@example.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("localhost", 1025, "", "")
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Failed to send email:", err)
	} else {
		fmt.Println("Email sent successfully!")
	}
}
