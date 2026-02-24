package main

import (
	"fmt"
	"net"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// portCheckResultMsg is returned by the background TCP check.
type portCheckResultMsg struct {
	host string
	port string
	open bool
	rtt  time.Duration
	err  error
}

// doPortCheck performs a TCP dial in the background and returns the result.
func doPortCheck(host, port string) tea.Cmd {
	return func() tea.Msg {
		target := net.JoinHostPort(host, port)
		start := time.Now()
		conn, err := net.DialTimeout("tcp", target, 5*time.Second)
		rtt := time.Since(start)
		if err != nil {
			return portCheckResultMsg{host: host, port: port, open: false, err: err}
		}
		conn.Close()
		return portCheckResultMsg{host: host, port: port, open: true, rtt: rtt}
	}
}

const (
	inputAddr = 0
	inputPort = 1
)

// portCheckModel holds the state for the CHECK PORT sub-screen.
type portCheckModel struct {
	inputs   [2]textinput.Model
	focused  int
	scanning bool
	width    int
	height   int
}

func newPortCheckModel() portCheckModel {
	addr := textinput.New()
	addr.Placeholder = "192.168.1.1  or  hostname"
	addr.CharLimit = 253
	addr.Width = 38
	addr.Prompt = "▶ "
	addr.PromptStyle = lipgloss.NewStyle().Foreground(colorGreen)
	addr.TextStyle = lipgloss.NewStyle().Foreground(colorHighlight)
	addr.PlaceholderStyle = lipgloss.NewStyle().Foreground(colorDimGreen)

	port := textinput.New()
	port.Placeholder = "80"
	port.CharLimit = 5
	port.Width = 10
	port.Prompt = "▶ "
	port.PromptStyle = lipgloss.NewStyle().Foreground(colorGreen)
	port.TextStyle = lipgloss.NewStyle().Foreground(colorHighlight)
	port.PlaceholderStyle = lipgloss.NewStyle().Foreground(colorDimGreen)

	m := portCheckModel{
		inputs:  [2]textinput.Model{addr, port},
		focused: inputAddr,
	}
	_ = m.inputs[inputAddr].Focus()
	return m
}

func (p portCheckModel) Init() tea.Cmd {
	return textinput.Blink
}

func (p portCheckModel) Update(msg tea.Msg) (portCheckModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "up", "down":
			if p.focused == inputAddr {
				p.focused = inputPort
			} else {
				p.focused = inputAddr
			}
			p.inputs[inputAddr].Blur()
			p.inputs[inputPort].Blur()
			cmd := p.inputs[p.focused].Focus()
			return p, cmd

		case "enter":
			host := p.inputs[inputAddr].Value()
			port := p.inputs[inputPort].Value()
			if host == "" || port == "" || p.scanning {
				return p, nil
			}
			p.scanning = true
			return p, doPortCheck(host, port)
		}
	}

	var cmds [2]tea.Cmd
	for i := range p.inputs {
		p.inputs[i], cmds[i] = p.inputs[i].Update(msg)
	}
	return p, tea.Batch(cmds[0], cmds[1])
}

func (p portCheckModel) View() string {
	focusedStyle := lipgloss.NewStyle().Foreground(colorHighlight).Bold(true)
	dimLabelStyle := lipgloss.NewStyle().Foreground(colorDimGreen)
	helpStyle := lipgloss.NewStyle().Foreground(colorDimGreen)

	addrLabel := dimLabelStyle.Render("TARGET ADDRESS")
	portLabel := dimLabelStyle.Render("PORT")
	if p.focused == inputAddr {
		addrLabel = focusedStyle.Render("TARGET ADDRESS")
	} else {
		portLabel = focusedStyle.Render("PORT")
	}

	title := menuTitleStyle.Render("[ CHECK PORT ]")

	addrBlock := lipgloss.JoinVertical(lipgloss.Left,
		"  "+addrLabel,
		"  "+p.inputs[inputAddr].View(),
	)
	portBlock := lipgloss.JoinVertical(lipgloss.Left,
		"  "+portLabel,
		"  "+p.inputs[inputPort].View(),
	)
	var statusLine string
	if p.scanning {
		statusLine = lipgloss.NewStyle().Foreground(colorHighlight).Bold(true).Render(
			fmt.Sprintf("  SCANNING %s:%s ...",
				p.inputs[inputAddr].Value(),
				p.inputs[inputPort].Value()),
		)
	}

	help := helpStyle.Render("  [TAB] switch field   [ENTER] check   [ESC] back")

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		addrBlock,
		"",
		portBlock,
		"",
		statusLine,
		"",
		help,
	)
}
