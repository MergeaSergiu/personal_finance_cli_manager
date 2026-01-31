package main

import (
	"log"
	"peronal_finance_cli_manager/internal/db"
	"peronal_finance_cli_manager/internal/models"
	"peronal_finance_cli_manager/internal/ui"

	_ "github.com/charmbracelet/bubbletea"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	db.Connect()

	err := db.DB.AutoMigrate(&models.Category{}, &models.Transaction{})
	if err != nil {
		log.Fatal(err)
	}

	menu := ui.NewMenuModel()

	p := tea.NewProgram(menu)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}

}
