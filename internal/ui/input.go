package ui

import (
	"fmt"
	"peronal_finance_cli_manager/internal/db"
	"peronal_finance_cli_manager/internal/models"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	_ "github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	focusedStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("1"))
)

type FileInputModel struct {
	input  textinput.Model
	errMsg string
	focus  bool
}

func NewFileInputModel() *FileInputModel {
	ti := textinput.New()
	ti.Placeholder = "Enter CSV/OFX file path"
	ti.Focus()
	return &FileInputModel{
		input: ti,
		focus: true,
	}
}

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
		switch msg.Type {

		case tea.KeyTab:
			m.focusIndex = (m.focusIndex + 1) % 2
			m.updateFocus()
			return m, nil, nil, nil

		case tea.KeyEnter:

			return m.submit()

		case tea.KeyEsc:
			return m, nil, nil, nil
		}
	}

	var cmd1, cmd2 tea.Cmd
	m.input, cmd1 = m.input.Update(msg)
	m.inputBudget, cmd2 = m.inputBudget.Update(msg)

	return m, tea.Batch(cmd1, cmd2), nil, nil
}

func (m *InputModel) updateFocus() {
	m.input.Blur()
	m.inputBudget.Blur()

	if m.focusIndex == 0 {
		m.input.Focus()
	} else {
		m.inputBudget.Focus()
	}
}

func (m *InputModel) submit() (
	*InputModel,
	tea.Cmd,
	*models.Category,
	error,
) {

	if m.input.Value() == "" {
		m.errMsg = "Invalid name"
		return m, nil, nil, nil
	}

	var budget float32 = 0
	if m.inputBudget.Value() != "" {
		parsed, err := strconv.ParseFloat(m.inputBudget.Value(), 32)
		if err != nil {
			m.errMsg = "Invalid budget"
			return m, nil, nil, nil
		}
		budget = float32(parsed)
	}

	cat, err := db.CreateCategory(m.input.Value(), budget)
	if err != nil {
		m.errMsg = err.Error()
		return m, nil, nil, err
	}

	m.errMsg = ""
	m.reset()

	return m, nil, cat, nil
}

func (m *InputModel) reset() {
	m.input.SetValue("")
	m.inputBudget.SetValue("")
	m.focusIndex = 0
	m.updateFocus()
}

func (m *TransactionInputModel) Update(msg tea.Msg) (*TransactionInputModel, tea.Cmd, *models.Transaction, error) {

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyTab:
			m.focusIndex = (m.focusIndex + 1) % 3
			m.updateFocus()
			return m, nil, nil, nil

		case tea.KeyEnter:
			return m.submit()

		case tea.KeyEsc:
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

func (m *FileInputModel) Update(msg tea.Msg) (*FileInputModel, tea.Cmd, string, error) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			path := m.input.Value()
			if path == "" {
				m.errMsg = "File path cannot be empty"
				return m, nil, "", nil
			}

			// Import transactions via db package
			imported, err := db.ImportTransactionsFromFile(path)
			if err != nil {
				m.errMsg = fmt.Sprintf("Import failed: %v", err)
				return m, nil, "", err
			}

			fmt.Printf("Imported %d transactions successfully\n", len(imported))
			return m, nil, path, nil

		case tea.KeyEsc:
			return m, nil, "", nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd, "", nil
}

func (m *TransactionInputModel) updateFocus() {
	m.inputCategory.Blur()
	m.inputAmount.Blur()
	m.inputDate.Blur()

	switch m.focusIndex {
	case 0:
		m.inputCategory.Focus()
	case 1:
		m.inputAmount.Focus()
	case 2:
		m.inputDate.Focus()
	}
}

func (m *TransactionInputModel) submit() (
	*TransactionInputModel,
	tea.Cmd,
	*models.Transaction,
	error,
) {

	category := m.inputCategory.Value()
	if category == "" {
		m.errMsg = "Category cannot be empty"
		return m, nil, nil, nil
	}

	amount, err := strconv.ParseFloat(m.inputAmount.Value(), 32)
	if err != nil {
		m.errMsg = "Invalid amount"
		return m, nil, nil, nil
	}

	dateStr := m.inputDate.Value()
	if dateStr == "" {
		m.errMsg = "Date cannot be empty (YYYY-MM-DD)"
		return m, nil, nil, nil
	}

	tx, err := db.CreateTransaction(
		category,
		float32(amount),
		dateStr,
	)
	if err != nil {
		m.errMsg = err.Error()
		return m, nil, nil, nil
	}

	m.errMsg = ""
	m.reset()

	return m, nil, tx, nil
}

func (m *TransactionInputModel) reset() {
	m.inputCategory.SetValue("")
	m.inputAmount.SetValue("")
	m.inputDate.SetValue("")
	m.focusIndex = 0
	m.updateFocus()
}

// View renders the input box
func (m *InputModel) View() string {

	view := ""

	if m.errMsg != "" {
		view += errorStyle.Render("❌ " + m.errMsg)
		view += "\n\n"
	}

	view += renderInput(m.input, m.focusIndex == 0) + "\n"
	view += renderInput(m.inputBudget, m.focusIndex == 1)

	view += "\n\n[Tab] Switch • [Enter] Save • [b] Back"
	return view
}

func (m *TransactionInputModel) View() string {
	view := ""

	if m.errMsg != "" {
		view += errorStyle.Render("❌ " + m.errMsg)
		view += "\n\n"
	}

	view += renderInput(m.inputCategory, m.focusIndex == 0) + "\n"
	view += renderInput(m.inputAmount, m.focusIndex == 1) + "\n"
	view += renderInput(m.inputDate, m.focusIndex == 2)

	view += "\n\n[Tab] Next • [Enter] Save • [b] Back"
	return view
}

func (m *FileInputModel) View() string {
	view := ""
	if m.errMsg != "" {
		view += errorStyle.Render("❌ " + m.errMsg)
		view += "\n\n"
	}
	view += renderInput(m.input, m.focus)
	view += "\n\n[Enter] Import • [Esc] Back"
	return view
}

func renderInput(input textinput.Model, focused bool) string {
	if focused {
		return focusedStyle.Render(input.View())
	}
	return input.View()
}
