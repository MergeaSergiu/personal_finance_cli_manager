package ui

import (
	"peronal_finance_cli_manager/internal/db"
	"peronal_finance_cli_manager/internal/models"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TransactionInputModel struct {
	inputCategory textinput.Model
	inputAmount   textinput.Model
	inputDate     textinput.Model
	focusIndex    int
	errMsg        string
}

func NewTransactionInputModel() *TransactionInputModel {
	catInput := textinput.New()
	catInput.Placeholder = "Category Name"
	catInput.CharLimit = 64
	catInput.Focus()

	amountInput := textinput.New()
	amountInput.Placeholder = "Amount"
	amountInput.CharLimit = 10
	amountInput.Blur()

	dateInput := textinput.New()
	dateInput.Placeholder = "Date"
	dateInput.CharLimit = 10
	dateInput.Blur()

	return &TransactionInputModel{
		inputCategory: catInput,
		inputAmount:   amountInput,
		inputDate:     dateInput,
		focusIndex:    0,
	}
}

type InputModel struct {
	input       textinput.Model
	inputBudget textinput.Model
	focusIndex  int
	errMsg      string
}

func NewInputModelPtr() *InputModel {
	ti := textinput.New()
	ti.Placeholder = "Category name"
	ti.CharLimit = 64
	ti.Focus()

	budget := textinput.New()
	budget.Placeholder = "Budget"
	budget.CharLimit = 10
	budget.Blur()

	return &InputModel{
		input:       ti,
		inputBudget: budget,
		focusIndex:  0,
	}
}

// Update handles key presses in the input form
func (m *InputModel) Update(msg tea.Msg) (*InputModel, tea.Cmd, *models.Category, error) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "tab":
			m.focusIndex = (m.focusIndex + 1) % 2
			if m.focusIndex == 0 {
				m.input.Focus()
				m.inputBudget.Blur()
			} else {
				m.input.Blur()
				m.inputBudget.Focus()
			}

			return m, nil, nil, nil

		case "enter":

			budget := 0.0
			if m.inputBudget.Value() != "" {
				parsed, err := strconv.ParseFloat(m.inputBudget.Value(), 32)
				if err != nil {
					m.errMsg = "Invalid budget"
					return m, nil, nil, nil
				}
				budget = parsed
			}

			if m.input.Value() == "" {
				m.errMsg = "Invalid Name"
				return m, nil, nil, nil
			}
			cat, err := db.CreateCategory(m.input.Value(), float32(budget))
			if err != nil {

				m.errMsg = err.Error()
				return m, nil, nil, err
			}
			m.errMsg = ""
			return m, nil, cat, nil

		case "b":
			return m, nil, nil, nil
		}
	}

	var cmd1, cmd2 tea.Cmd
	m.input, cmd1 = m.input.Update(msg)
	m.inputBudget, cmd2 = m.inputBudget.Update(msg)

	return m, tea.Batch(cmd1, cmd2), nil, nil
}

func (m *TransactionInputModel) Update(msg tea.Msg) (*TransactionInputModel, tea.Cmd, *models.Transaction, error) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "tab":
			m.focusIndex = (m.focusIndex + 1) % 3
			if m.focusIndex == 0 {
				m.inputCategory.Focus()
				m.inputAmount.Blur()
				m.inputDate.Blur()
			} else if m.focusIndex == 1 {
				m.inputCategory.Blur()
				m.inputAmount.Focus()
				m.inputDate.Blur()
			} else {
				m.inputCategory.Blur()
				m.inputAmount.Blur()
				m.inputDate.Focus()
			}
			return m, nil, nil, nil

		case "enter":
			// parse amount
			amount := float32(0)
			if m.inputAmount.Value() != "" {
				val, err := strconv.ParseFloat(m.inputAmount.Value(), 32)
				if err != nil {
					m.errMsg = "Invalid amount"
					return m, nil, nil, nil
				}
				amount = float32(val)
			} else {
				m.errMsg = "Amount cannot be empty"
				return m, nil, nil, nil
			}

			categoryName := m.inputCategory.Value()
			if categoryName == "" {
				m.errMsg = "Category cannot be empty"
				return m, nil, nil, nil
			}

			dateStr := m.inputDate.Value()
			if dateStr == "" {
				m.errMsg = "Date cannot be empty (YYYY-MM-DD)"
				return m, nil, nil, nil
			}

			tx, err := db.CreateTransaction(categoryName, amount, dateStr)
			if err != nil {
				m.errMsg = err.Error()
				m.inputCategory.SetValue("") // reset category if not found
				m.inputAmount.SetValue("")
				m.inputDate.SetValue("")
				return m, nil, nil, nil
			}

			m.errMsg = ""
			return m, nil, tx, nil

		case "b":
			return m, nil, nil, nil
		}
	}

	// update all 3 inputs
	var cmd1, cmd2, cmd3 tea.Cmd
	m.inputCategory, cmd1 = m.inputCategory.Update(msg)
	m.inputAmount, cmd2 = m.inputAmount.Update(msg)
	m.inputDate, cmd3 = m.inputDate.Update(msg)

	return m, tea.Batch(cmd1, cmd2, cmd3), nil, nil
}

// View renders the input box
func (m *InputModel) View() string {

	view := ""

	if m.errMsg != "" {
		view += "❌ " + m.errMsg + "\n\n"
	}
	view += m.input.View() + "\n"
	view += m.inputBudget.View()

	view += "\n\n[Tab] Switch field • [Enter] Save • [b] Back"

	return view
}

func (m *TransactionInputModel) View() string {
	view := ""

	if m.errMsg != "" {
		view += "❌ " + m.errMsg + "\n\n"
	}

	view += m.inputCategory.View() + "\n"
	view += m.inputAmount.View() + "\n"
	view += m.inputDate.View()

	view += "\n\n[Tab] Switch field • [Enter] Save • [b] Back"

	return view
}
