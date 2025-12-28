package ui

import (
	_ "encoding/csv"
	"fmt"
	_ "os"
	"peronal_finance_cli_manager/internal/db"
	"peronal_finance_cli_manager/internal/models"
	"peronal_finance_cli_manager/internal/transaction"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	StateList state = iota
	StateAdd
	StateView
	StateAddTransaction
	StateViewTransactions
	StateImportCSV
)

type MenuModel struct {
	list                  list.Model
	inputModel            *InputModel
	transactionInputModel *TransactionInputModel
	state                 state

	transactions     []models.Transaction
	selectedCategory *models.Category

	importInput textinput.Model
	importMsg   string
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

	ti := textinput.New()
	ti.Placeholder = "Enter CSV file path..."
	ti.CharLimit = 256
	ti.Focus()

	return &MenuModel{
		list:                  l,
		inputModel:            NewInputModelPtr(),
		transactionInputModel: NewTransactionInputModel(),
		importInput:           ti,
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
	//if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEsc {
	//	switch m.state {
	//	case StateAdd, StateAddTransaction, StateImportCSV:
	//		m.state = StateList
	//		m.inputModel.reset()
	//		m.transactionInputModel.reset()
	//	case StateView:
	//		m.state = StateList
	//	case StateViewTransactions:
	//		m.state = StateView
	//	}
	//	return m, nil
	//}

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
				//m.inputModel.input.SetValue("")
				//m.inputModel.inputBudget.SetValue("")
				//m.inputModel.focusIndex = 0
				//m.inputModel.errMsg = ""
				m.inputModel.reset()
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
				//m.transactionInputModel.inputCategory.SetValue("")
				//m.transactionInputModel.inputAmount.SetValue("")
				//m.transactionInputModel.inputDate.SetValue("")
				//m.transactionInputModel.focusIndex = 0
				//m.transactionInputModel.errMsg = ""
				m.transactionInputModel.reset()
				m.state = StateAddTransaction
				return m, nil

			case "i":
				//fmt.Print("Enter CSV file path: ")
				//var path string
				//fmt.Scanln(&path)
				//
				//txs, err := transaction.ParseCSV(path)
				//if err != nil {
				//	fmt.Println("Error parsing CSV:", err)
				//	return m, nil
				//}
				//
				//for _, tx := range txs {
				//	_, err := db.CreateTransaction(tx.Category.Name, tx.Amount, tx.Date.Format("2006-01-02"))
				//	if err != nil {
				//		fmt.Println("Error saving transaction:", err)
				//	}
				//}
				//
				//fmt.Printf("✅ Imported %d transactions\n", len(txs))
				m.importInput.SetValue("")
				m.importMsg = ""
				m.state = StateImportCSV
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

		//if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "b" && cat == nil {
		//	m.state = StateList
		//	m.inputModel.input.SetValue("")
		//	m.inputModel.inputBudget.SetValue("")
		//	m.inputModel.errMsg = ""
		//}

		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "b" {
			m.state = StateList
			m.inputModel.reset()
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
			case "b":
				m.state = StateList
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
			m.transactionInputModel.reset()
		}
		return m, cmd

	// ====================== VIEW TRANSACTIONS ======================
	case StateViewTransactions:
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "b" {
			m.state = StateView
		}
		return m, nil

	case StateImportCSV:
		var cmd tea.Cmd
		m.importInput, cmd = m.importInput.Update(msg)

		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "enter":

				filePath := m.importInput.Value()
				format := transaction.DetectFormat(filePath)
				if format != "csv" {
					m.importMsg = "❌ Unsupported format"
					return m, nil
				}

				transactions, err := transaction.ParseCSV(filePath)
				if err != nil {
					m.importMsg = "❌ Error parsing CSV: " + err.Error()
					return m, nil
				}

				count := 0
				for _, tx := range transactions {
					// Check if category exists
					cat, err := db.GetCategoryByName(tx.Category.Name)
					if err != nil {
						// If not exists, create
						cat, err = db.CreateCategory(tx.Category.Name, 10000)
						if err != nil {
							continue // skip this transaction if category creation fails
						}
					}

					// Create transaction
					_, _ = db.CreateTransaction(cat.Name, tx.Amount, tx.Date.Format("2006-01-02"))
					count++
				}

				m.importMsg = fmt.Sprintf("✅ Imported %d transactions", count)
			case "b":
				m.state = StateList
			}
		}

		return m, cmd

	}

	return m, nil
}

// View renders the UI
func (m *MenuModel) View() string {

	switch m.state {
	case StateList:
		return "[v] View Categories • [a] Add category • [t] Add transaction • [i] Import CSV • [q] Quit"

	case StateAdd:
		return fmt.Sprintf(
			"➕ Add Category\n\n%s",
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
				"%s%.2f  |  %s\n",
				sign,
				tx.Amount,
				tx.Date.Format("2006-01-02"),
			)
		}

		view += "\n[b] Back"
		return view

	case StateImportCSV:
		view := fmt.Sprintf("📥 Import CSV\n\n%s", m.importInput.View())
		if m.importMsg != "" {
			view += "\n\n" + m.importMsg
		}
		view += "\n\n[Enter] Import • [b] Back"
		return view
	}

	return ""
}
