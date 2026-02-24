package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type screen int

const (
	screenMenu screen = iota
	screenCheckPort
	screenDNS
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
	width      int
	height     int
	screen     screen
	menu       list.Model
	portCheck  portCheckModel
	portResult *portCheckResultMsg
	dnsLookup  dnsModel
	dnsResult  *dnsLookupResultMsg
}

func initialModel() Model {
	items := []list.Item{
		menuItem{title: "DNS LOOKUP", desc: "Query the machines that know all names"},
		menuItem{title: "CHECK PORT", desc: "Knock on doors and see who answers"},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "[ SELECT OPERATION ]"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = menuTitleStyle
	l.Styles.NoItems = dimStyle

	return Model{menu: l, portCheck: newPortCheckModel(), dnsLookup: newDNSModel()}
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
		inputW := leftW - 4
		if inputW < 10 {
			inputW = 10
		}
		m.portCheck.inputs[inputAddr].Width = inputW
		m.portCheck.width = leftW
		m.portCheck.height = leftH
		m.dnsLookup.inputs[inputName].Width = inputW
		m.dnsLookup.width = leftW
		m.dnsLookup.height = leftH

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.screen == screenMenu {
				return m, tea.Quit
			}
		case "esc":
			if m.screen != screenMenu {
				m.screen = screenMenu
				return m, nil
			}
		case "enter":
			if m.screen == screenMenu {
				if item, ok := m.menu.SelectedItem().(menuItem); ok {
					switch item.title {
					case "CHECK PORT":
						m.screen = screenCheckPort
						m.portCheck = newPortCheckModel()
						m.portResult = nil
						leftW, _ := leftPaneDimensions(m.width, m.height)
						inputW := leftW - 4
						if inputW < 10 {
							inputW = 10
						}
						m.portCheck.inputs[inputAddr].Width = inputW
						return m, m.portCheck.Init()
					case "DNS LOOKUP":
						m.screen = screenDNS
						m.dnsLookup = newDNSModel()
						m.dnsResult = nil
						leftW, _ := leftPaneDimensions(m.width, m.height)
						inputW := leftW - 4
						if inputW < 10 {
							inputW = 10
						}
						m.dnsLookup.inputs[inputName].Width = inputW
						return m, m.dnsLookup.Init()
					}
				}
				return m, nil
			}
		}
	case portCheckResultMsg:
		m.portCheck.scanning = false
		m.portResult = &msg
		return m, nil
	case dnsLookupResultMsg:
		m.dnsLookup.scanning = false
		m.dnsResult = &msg
		return m, nil
	}

	// Route remaining messages to the active screen.
	var cmd tea.Cmd
	switch m.screen {
	case screenCheckPort:
		m.portCheck, cmd = m.portCheck.Update(msg)
	case screenDNS:
		m.dnsLookup, cmd = m.dnsLookup.Update(msg)
	default:
		m.menu, cmd = m.menu.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	if m.width == 0 {
		return "" // waiting for first WindowSizeMsg
	}
	return renderView(m)
}
