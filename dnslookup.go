package main

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	inputDNS1 = 0
	inputDNS2 = 1
	inputName = 2
)

// dnsServerResult holds the response from one DNS server.
type dnsServerResult struct {
	server string
	addrs  []string
	cname  string
	err    error
}

// dnsLookupResultMsg is returned after querying both DNS servers.
type dnsLookupResultMsg struct {
	name    string
	results [2]dnsServerResult
}

// queryDNS queries a single DNS server for the given hostname.
func queryDNS(ctx context.Context, server, name string) dnsServerResult {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: 5 * time.Second}
			return d.DialContext(ctx, "udp", net.JoinHostPort(server, "53"))
		},
	}

	cname, _ := r.LookupCNAME(ctx, name)
	cname = strings.TrimSuffix(cname, ".")
	// LookupCNAME returns the original name when there's no redirect.
	if strings.EqualFold(cname, strings.TrimSuffix(name, ".")) {
		cname = ""
	}

	addrs, err := r.LookupHost(ctx, name)
	return dnsServerResult{server: server, addrs: addrs, cname: cname, err: err}
}

// doDNSLookup runs both DNS server queries concurrently and returns one combined msg.
func doDNSLookup(dns1, dns2, name string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		type indexed struct {
			i   int
			res dnsServerResult
		}
		ch := make(chan indexed, 2)
		go func() { ch <- indexed{0, queryDNS(ctx, dns1, name)} }()
		go func() { ch <- indexed{1, queryDNS(ctx, dns2, name)} }()

		var results [2]dnsServerResult
		for range 2 {
			r := <-ch
			results[r.i] = r.res
		}
		return dnsLookupResultMsg{name: name, results: results}
	}
}

// dnsModel holds the state for the DNS LOOKUP sub-screen.
type dnsModel struct {
	inputs   [3]textinput.Model
	focused  int
	scanning bool
	width    int
	height   int
}

func newDNSModel() dnsModel {
	make := func(placeholder, defaultVal string, width int) textinput.Model {
		t := textinput.New()
		t.Placeholder = placeholder
		t.SetValue(defaultVal)
		t.CharLimit = 253
		t.Width = width
		t.Prompt = "▶ "
		t.PromptStyle = lipgloss.NewStyle().Foreground(colorGreen)
		t.TextStyle = lipgloss.NewStyle().Foreground(colorHighlight)
		t.PlaceholderStyle = lipgloss.NewStyle().Foreground(colorDimGreen)
		return t
	}

	dns1 := make("1.1.1.1", "1.1.1.1", 20)
	dns2 := make("8.8.8.8", "8.8.8.8", 20)
	name := make("example.com", "", 38)

	m := dnsModel{
		inputs:  [3]textinput.Model{dns1, dns2, name},
		focused: inputName,
	}
	_ = m.inputs[inputName].Focus()
	return m
}

func (d dnsModel) Init() tea.Cmd {
	return textinput.Blink
}

func (d dnsModel) Update(msg tea.Msg) (dnsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			d.focused = (d.focused + 1) % 3
			for i := range d.inputs {
				d.inputs[i].Blur()
			}
			return d, d.inputs[d.focused].Focus()

		case "shift+tab", "up":
			d.focused = (d.focused + 2) % 3
			for i := range d.inputs {
				d.inputs[i].Blur()
			}
			return d, d.inputs[d.focused].Focus()

		case "enter":
			dns1 := d.inputs[inputDNS1].Value()
			dns2 := d.inputs[inputDNS2].Value()
			name := d.inputs[inputName].Value()
			if name == "" || dns1 == "" || dns2 == "" || d.scanning {
				return d, nil
			}
			d.scanning = true
			return d, doDNSLookup(dns1, dns2, name)
		}
	}

	var cmds [3]tea.Cmd
	for i := range d.inputs {
		d.inputs[i], cmds[i] = d.inputs[i].Update(msg)
	}
	return d, tea.Batch(cmds[0], cmds[1], cmds[2])
}

func (d dnsModel) View() string {
	focused := lipgloss.NewStyle().Foreground(colorHighlight).Bold(true)
	dim := lipgloss.NewStyle().Foreground(colorDimGreen)
	help := lipgloss.NewStyle().Foreground(colorDimGreen)

	label := func(idx int, text string) string {
		if d.focused == idx {
			return focused.Render(text)
		}
		return dim.Render(text)
	}

	title := menuTitleStyle.Render("[ DNS LOOKUP ]")

	dns1Block := lipgloss.JoinVertical(lipgloss.Left,
		"  "+label(inputDNS1, "DNS SERVER 1"),
		"  "+d.inputs[inputDNS1].View(),
	)
	dns2Block := lipgloss.JoinVertical(lipgloss.Left,
		"  "+label(inputDNS2, "DNS SERVER 2"),
		"  "+d.inputs[inputDNS2].View(),
	)
	nameBlock := lipgloss.JoinVertical(lipgloss.Left,
		"  "+label(inputName, "HOSTNAME"),
		"  "+d.inputs[inputName].View(),
	)

	var statusLine string
	if d.scanning {
		statusLine = lipgloss.NewStyle().Foreground(colorHighlight).Bold(true).Render(
			"  QUERYING " + d.inputs[inputName].Value() + " ...",
		)
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		dns1Block,
		"",
		dns2Block,
		"",
		nameBlock,
		"",
		statusLine,
		"",
		help.Render("  [TAB] switch field   [ENTER] lookup   [ESC] back"),
	)
}
