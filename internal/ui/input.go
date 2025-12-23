package ui

import (
	"peronal_finance_cli_manager/internal/db"
	"peronal_finance_cli_manager/internal/models"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type InputModel struct {
	input  textinput.Model
	errMsg string
}

func NewInputModelPtr() *InputModel {
	ti := textinput.New()
	ti.Placeholder = "Category name"
	ti.Focus()
	ti.CharLimit = 64

	return &InputModel{
		input: ti,
	}
}

// Update handles key presses in the input form
func (m *InputModel) Update(msg tea.Msg) (*InputModel, tea.Cmd, *models.Category, error) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			cat, err := db.CreateCategory(m.input.Value())
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

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd, nil, nil
}

// View renders the input box
func (m *InputModel) View() string {
	if m.errMsg != "" {
		return "❌ " + m.errMsg + "\n\n" + m.input.View()
	}
	return m.input.View()
}
