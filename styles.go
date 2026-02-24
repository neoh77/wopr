package main

import "github.com/charmbracelet/lipgloss"

// phosphor green palette
const (
	colorGreen     = lipgloss.Color("#00FF41")
	colorDimGreen  = lipgloss.Color("#006B1B")
	colorBg        = lipgloss.Color("#0D0D0D")
	colorHighlight = lipgloss.Color("#39FF14")
)

// rightPaneWidth is fixed so the WOPR banner always fits cleanly.
const rightPaneWidth = 42

// woprBanner is rendered verbatim inside the right pane.
const woprBanner = `
 ██╗    ██╗ ██████╗ ██████╗ ██████╗
 ██║    ██║██╔═══██╗██╔══██╗██╔══██╗
 ██║ █╗ ██║██║   ██║██████╔╝██████╔╝
 ██║███╗██║██║   ██║██╔═══╝ ██╔══██╗
 ╚███╔███╔╝╚██████╔╝██║     ██║  ██║
  ╚══╝╚══╝  ╚═════╝ ╚═╝     ╚═╝  ╚═╝`

var (
	// leftPaneStyle wraps the menu list.
	leftPaneStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(colorDimGreen).
			Padding(0, 1)

	// rightPaneStyle wraps the WOPR banner + blank area.
	rightPaneStyle = lipgloss.NewStyle().
			Width(rightPaneWidth).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(colorGreen)

	// bannerStyle colors the ASCII art.
	bannerStyle = lipgloss.NewStyle().
			Foreground(colorGreen).
			Bold(true)

	// menuTitleStyle is used by the list widget's title.
	menuTitleStyle = lipgloss.NewStyle().
			Foreground(colorHighlight).
			Bold(true).
			MarginLeft(1)

	// dimStyle for secondary text.
	dimStyle = lipgloss.NewStyle().Foreground(colorDimGreen)

	// resultOpenStyle for an open port verdict.
	resultOpenStyle = lipgloss.NewStyle().Foreground(colorHighlight).Bold(true)

	// resultClosedStyle for a closed/filtered port verdict.
	resultClosedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF3333")).Bold(true)

	// resultDividerStyle for the separator line in the results block.
	resultDividerStyle = lipgloss.NewStyle().Foreground(colorDimGreen)
)
