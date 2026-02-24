package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// menuItem represents a single entry in the left-pane menu.
type menuItem struct {
	title string
	desc  string
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.desc }
func (i menuItem) FilterValue() string { return i.title }

// Model is the root application model.
type Model struct {
	width  int
	height int
	menu   list.Model
}

func initialModel() Model {
	items := []list.Item{
		menuItem{title: "PING", desc: "Probe a host for signs of life"},
		menuItem{title: "DNS LOOKUP", desc: "Query the machines that know all names"},
		menuItem{title: "PORT SCAN", desc: "Knock on doors and see who answers"},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "[ SELECT OPERATION ]"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = menuTitleStyle
	l.Styles.NoItems = dimStyle

	return Model{menu: l}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		leftW, leftH := leftPaneDimensions(m.width, m.height)
		m.menu.SetWidth(leftW)
		m.menu.SetHeight(leftH)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.menu, cmd = m.menu.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.width == 0 {
		return "" // waiting for first WindowSizeMsg
	}
	return renderView(m)
}
