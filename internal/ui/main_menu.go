package ui

import (
	_ "encoding/csv"
	"fmt"
	_ "os"
	"peronal_finance_cli_manager/internal/db"
	"peronal_finance_cli_manager/internal/models"
	"peronal_finance_cli_manager/internal/transaction"
	"strconv"
	"time"

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
	StateFilterTransactions state = iota + 100
	StateBudgetOverview
)

type FilterTransactionsModel struct {
	input        textinput.Model
	transactions []models.Transaction
	filtered     []models.Transaction
	mode         string // "date", "beforeDate", "year"
	modes        []string
	errMsg       string
}

type MenuModel struct {
	list                  list.Model
	inputModel            *InputModel
	transactionInputModel *TransactionInputModel
	state                 state

	transactions     []models.Transaction
	selectedCategory *models.Category

	importInput textinput.Model
	importMsg   string

	filterModel *FilterTransactionsModel
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

func NewFilterTransactionsModel(txs []models.Transaction, mode string) *FilterTransactionsModel {
	ti := textinput.New()
	ti.Placeholder = "Enter filter value"
	ti.Focus()

	modes := []string{"date", "beforeDate", "year"}

	return &FilterTransactionsModel{
		input:        ti,
		transactions: txs,
		mode:         mode,
		modes:        modes,
	}
}

// Init satisfies tea.Model interface
func (m *MenuModel) Init() tea.Cmd {
	return nil
}

func (m *FilterTransactionsModel) Update(msg tea.Msg) (*FilterTransactionsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			value := m.input.Value()
			m.errMsg = ""

			switch m.mode {
			case "date":
				date, err := time.Parse("2006-01-02", value)
				if err != nil {
					m.errMsg = "Invalid date format (YYYY-MM-DD)"
					return m, nil
				}
				m.filtered = transaction.FilterByExactDate(m.transactions, date)

			case "beforeDate":
				date, err := time.Parse("2006-01-02", value)
				if err != nil {
					m.errMsg = "Invalid date format (YYYY-MM-DD)"
					return m, nil
				}
				m.filtered = transaction.FilterBeforeDate(m.transactions, date)

			case "year":
				y, err := strconv.Atoi(value)
				if err != nil {
					m.errMsg = "Invalid year"
					return m, nil
				}
				m.filtered = transaction.FilterByYear(m.transactions, y)

			}

			return m, nil

		case tea.KeyRunes:
			if len(msg.Runes) > 0 {
				switch msg.Runes[0] {
				case 'f':
					// cycle filter mode
					idx := 0
					for i, mname := range m.modes {
						if mname == m.mode {
							idx = i
							break
						}
					}
					idx = (idx + 1) % len(m.modes)
					m.mode = m.modes[idx]
					m.input.SetValue("")
					m.filtered = nil
					m.errMsg = ""
					return m, nil

				case 'b':
					// Back to previous menu
					return m, nil
				}
			}
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// Update handles key presses
func (m *MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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
				m.transactionInputModel.reset()
				m.state = StateAddTransaction
				return m, nil

			case "i":
				m.importInput.SetValue("")
				m.importMsg = ""
				m.state = StateImportCSV
				return m, nil

			case "p":
				m.state = StateBudgetOverview
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
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "b": // Back
				m.state = StateView
			case "f": // Open filter menu
				// Here we let user select the filter mode first (hardcoded "date" for example)
				// Later we can add a dynamic selection menu for mode
				m.filterModel = NewFilterTransactionsModel(m.transactions, "date")
				m.state = StateFilterTransactions
			}
		}
		return m, nil

	case StateFilterTransactions:
		var cmd tea.Cmd
		m.filterModel, cmd = m.filterModel.Update(msg)
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "b": // Go back to transactions view
				m.state = StateViewTransactions
				m.filterModel = nil
			case "f": // Change filter mode (optional, if implemented)
				// Logic to switch filter mode
			}
		}
		return m, cmd

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

	case StateBudgetOverview:
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "b" {
			m.state = StateList
		}
		return m, nil
	}

	return m, nil
}

// View renders the UI
func (m *MenuModel) View() string {

	switch m.state {
	case StateList:
		return "[v] View Categories • [p] Budget overview • [a] Add category • [t] Add transaction • [i] Import CSV • [q] Quit"

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

		view += "\n[b] Back • [f] Filter Transactions • [q] Quit"
		return view

	case StateImportCSV:
		view := fmt.Sprintf("📥 Import CSV\n\n%s", m.importInput.View())
		if m.importMsg != "" {
			view += "\n\n" + m.importMsg
		}
		view += "\n\n[Enter] Import • [b] Back"
		return view

	case StateFilterTransactions:
		if m.filterModel != nil {
			return m.filterModel.View()
		}

	case StateBudgetOverview:
		stats, err := db.GetBudgetStats()
		if err != nil {
			return "❌ Failed to load budget stats\n\n[b] Back"
		}

		view := "📊 Budget Overview\n\n"

		view += headerStyle.Render(
			fmt.Sprintf("%-16s %7s / %-7s %6s\n",
				"Category",
				"Spent",
				"Budget",
				"%"),
		)
		view += "\n"

		for _, s := range stats {
			percent := calculatePercentage(s.Spent, float32(s.Budget))

			base := fmt.Sprintf(
				"%-16s ",
				s.CategoryName,
			)

			numbers := fmt.Sprintf(
				"%7.0f / %-7.0f %6.0f%%",
				s.Spent,
				s.Budget,
				percent,
			)

			switch {
			case percent >= 100:
				view += base + redStyle.Render(numbers) + "\n"
			case percent >= 80:
				view += base + orangeStyle.Render(numbers) + "\n"
			default:
				view += base + greenStyle.Render(numbers) + "\n"
			}
		}

		view += "\n[b] Back"
		return view
	}

	return ""
}
