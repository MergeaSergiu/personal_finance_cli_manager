package ui

import (
	"peronal_finance_cli_manager/internal/db"
	"peronal_finance_cli_manager/internal/models"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type InputModel struct {
	input       textinput.Model
	inputBudget textinput.Model
	focusIndex  int
	errMsg      string
}

func NewInputModelPtr() *InputModel {
	ti := textinput.New()
	ti.Placeholder = "Category name"
	ti.Focus()
	ti.CharLimit = 64

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

		case "c":
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

// View renders the input box
func (m *InputModel) View() string {

	view := ""

	if m.errMsg != "" {
		return "❌ " + m.errMsg + "\n\n" + m.input.View() +
			"\n" + m.inputBudget.View()
	}
	view += m.input.View() + "\n"
	view += m.inputBudget.View()

	return view
}
