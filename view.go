package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// borderWidth accounts for the two border characters (left + right).
const borderWidth = 4

// leftPaneDimensions returns the usable inner width and height for the left pane,
// given the total terminal dimensions.
func leftPaneDimensions(totalW, totalH int) (w, h int) {
	w = totalW - rightPaneWidth - borderWidth*2
	if w < 1 {
		w = 1
	}
	// subtract 2 for the top/bottom borders, 2 for top/bottom padding
	h = totalH - borderWidth
	if h < 1 {
		h = 1
	}
	return w, h
}

// renderView composes the full terminal screen.
func renderView(m Model) string {
	leftW, _ := leftPaneDimensions(m.width, m.height)
	paneH := m.height - borderWidth

	// ── Left pane: menu ─────────────────────────────────────────────────────
	left := leftPaneStyle.
		Width(leftW).
		Height(paneH).
		Render(m.menu.View())

	// ── Right pane: WOPR banner + blank space ───────────────────────────────
	banner := bannerStyle.Render(woprBanner)
	blankLines := paneH - lipgloss.Height(banner)
	if blankLines < 0 {
		blankLines = 0
	}
	rightContent := lipgloss.JoinVertical(
		lipgloss.Left,
		banner,
		strings.Repeat("\n", blankLines),
	)
	right := rightPaneStyle.
		Height(paneH).
		Render(rightContent)

	// ── Combine side-by-side ────────────────────────────────────────────────
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}
