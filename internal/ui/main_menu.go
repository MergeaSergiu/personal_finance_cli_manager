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
)

type MenuModel struct {
	list       list.Model
	inputModel *InputModel
	state      state
}

type CategoryItem models.Category

func (c CategoryItem) Title() string       { return fmt.Sprintf("%s (ID: %d)", c.Name, c.ID) }
func (c CategoryItem) Description() string { return "" }
func (c CategoryItem) FilterValue() string { return c.Name }

// NewMenuModel creates main menu
func NewMenuModel() *MenuModel {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 50, 20)
	l.Title = "📂 Categories"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return &MenuModel{
		list:       l,
		inputModel: NewInputModelPtr(),
		state:      StateList,
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
		case StateAdd, StateView:
			m.state = StateList
			// Clear input if coming from Add
			m.inputModel.input.SetValue("")
			m.inputModel.errMsg = ""
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
	case StateList:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "a":
				m.state = StateAdd
				m.inputModel.input.SetValue("")
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
			}
		}

		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd

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
			m.inputModel.errMsg = ""
		}

		return m, cmd

	case StateView:
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)

		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "b" {
			m.state = StateList
		}

		return m, cmd
	}

	return m, nil
}

// View renders the UI
func (m *MenuModel) View() string {
	switch m.state {
	case StateList:
		return "[v] View Categories • [a] Add category • [q] Quit"
	case StateAdd:
		return fmt.Sprintf("➕ Add Category\n\n%s\n\n[Enter] Save • [b] Back", m.inputModel.View())
	case StateView:
		return fmt.Sprintf("%s\n\n[b] Back", m.list.View())
	}
	return ""
}
