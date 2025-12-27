package ui

import (
	"fmt"
	"peronal_finance_cli_manager/internal/db"
	"peronal_finance_cli_manager/internal/models"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	StateList state = iota
	StateAdd
	StateView
	StateAddTransaction
	StateViewTransactions
)

type MenuModel struct {
	list                  list.Model
	inputModel            *InputModel
	transactionInputModel *TransactionInputModel
	state                 state

	transactions     []models.Transaction
	selectedCategory *models.Category
}

type CategoryItem models.Category

func (c CategoryItem) Title() string {
	return fmt.Sprintf("%s (ID: %d), (Budget: %2.f)", c.Name, c.ID, c.Budget)
}
func (c CategoryItem) Description() string { return "" }
func (c CategoryItem) FilterValue() string { return c.Name }

// NewMenuModel creates main menu
func NewMenuModel() *MenuModel {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 50, 20)
	l.Title = "📂 Categories"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return &MenuModel{
		list:                  l,
		inputModel:            NewInputModelPtr(),
		transactionInputModel: NewTransactionInputModel(),
		state:                 StateList,
	}
}

// Init satisfies tea.Model interface
func (m *MenuModel) Init() tea.Cmd {
	return nil
}

// Update handles key presses
func (m *MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// 1️⃣ Handle ESC globally first
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
		switch m.state {
		case StateAdd, StateView, StateAddTransaction, StateViewTransactions:
			m.state = StateList
			if m.state == StateAdd {
				m.inputModel.input.SetValue("")
				m.inputModel.inputBudget.SetValue("")
				m.inputModel.errMsg = ""
			}
			return m, nil
		}
	}

	// 2️⃣ Handle window resize
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-4)
		return m, nil
	}

	switch m.state {

	// ====================== CATEGORY LIST ======================
	case StateList:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "a":
				m.state = StateAdd
				m.inputModel.input.SetValue("")
				m.inputModel.inputBudget.SetValue("")
				m.inputModel.focusIndex = 0
				m.inputModel.errMsg = ""
				return m, nil

			case "v":
				// Load categories from DB
				cats, err := db.GetAllCategories()
				if err != nil {
					// handle error
					fmt.Println("Error loading categories:", err)
					return m, nil
				}

				items := make([]list.Item, 0, len(cats))
				for _, c := range cats {
					items = append(items, CategoryItem(c))
				}
				m.list.SetItems(items)
				m.state = StateView
				return m, nil

			case "t": // Add Transaction
				m.transactionInputModel.inputCategory.SetValue("")
				m.transactionInputModel.inputAmount.SetValue("")
				m.transactionInputModel.focusIndex = 0
				m.transactionInputModel.errMsg = ""
				m.state = StateAddTransaction
				return m, nil

			}
		}

		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd

		// ====================== ADD CATEGORY ======================
	case StateAdd:
		var cmd tea.Cmd
		var cat *models.Category

		m.inputModel, cmd, cat, _ = m.inputModel.Update(msg)

		if cat != nil {
			m.list.InsertItem(len(m.list.Items()), CategoryItem(*cat))
			m.state = StateList
		}

		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "b" && cat == nil {
			m.state = StateList
			m.inputModel.input.SetValue("")
			m.inputModel.inputBudget.SetValue("")
			m.inputModel.errMsg = ""
		}

		return m, cmd

		// ====================== VIEW CATEGORY LIST ======================
	case StateView:
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)

		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {

			case "enter":
				item := m.list.SelectedItem()
				if item == nil {
					return m, cmd
				}

				cat := item.(CategoryItem)
				m.selectedCategory = (*models.Category)(&cat)

				txs, err := db.GetTransactionsByCategory(cat.ID)
				if err != nil {
					fmt.Println("Error loading transactions:", err)
					return m, cmd
				}

				m.transactions = txs
				m.state = StateViewTransactions
				return m, cmd

			case "b":
				m.state = StateList
				return m, cmd
			}
		}

		return m, cmd
	// ====================== ADD TRANSACTION ======================
	case StateAddTransaction:
		var cmd tea.Cmd
		var tx *models.Transaction

		m.transactionInputModel, cmd, tx, _ = m.transactionInputModel.Update(msg)

		if tx != nil {
			// append to transaction list

			m.transactions = append(m.transactions, *tx)
			m.state = StateList
		}

		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "b" {
			m.state = StateList
		}

		return m, cmd

	// ====================== VIEW TRANSACTIONS ======================
	case StateViewTransactions:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "b", "esc":
				m.state = StateView
				return m, nil
			}
		}
		return m, nil
	}

	return m, nil
}

// View renders the UI
func (m *MenuModel) View() string {

	switch m.state {
	case StateList:
		return "[v] View Categories • [a] Add category • [t] Add transaction • [q] Quit"

	case StateAdd:
		return fmt.Sprintf(
			"➕ Add Category\n\n%s\n\n[Tab] Switch field • [Enter] Save • [b] Back",
			m.inputModel.View())

	case StateView:
		return fmt.Sprintf("%s\n\n[b] Back", m.list.View())

	case StateAddTransaction:
		return fmt.Sprintf(
			"➕ Add Transaction\n\n%s",
			m.transactionInputModel.View(),
		)

	case StateViewTransactions:
		if m.selectedCategory == nil {
			return "📄 Transactions\n\nNo category selected.\n\n[b] Back"
		}

		if len(m.transactions) == 0 {
			return fmt.Sprintf(
				"📄 Transactions for %s\n\nNo transactions.\n\n[b] Back",
				m.selectedCategory.Name,
			)
		}

		view := fmt.Sprintf(
			"📄 Transactions for %s\n\n",
			m.selectedCategory.Name,
		)

		for _, tx := range m.transactions {
			sign := "+"
			if tx.Amount < 0 {
				sign = "-"
			}

			view += fmt.Sprintf(
				"%s%.2f\n",
				sign,
				tx.Amount,
			)
		}

		view += "\n[b] Back"
		return view
	}

	return ""
}
