package main

import (
	"fmt"
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

// trunc shortens s to max n runes, appending "…" if cut.
func trunc(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n-1]) + "…"
}

// renderView composes the full terminal screen.
func renderView(m Model) string {
	leftW, _ := leftPaneDimensions(m.width, m.height)
	paneH := m.height - borderWidth

	// ── Left pane: content switches on active screen ────────────────────────
	var leftContent string
	switch m.screen {
	case screenCheckPort:
		leftContent = m.portCheck.View()
	case screenDNS:
		leftContent = m.dnsLookup.View()
	default:
		leftContent = m.menu.View()
	}
	left := leftPaneStyle.
		Width(leftW).
		Height(paneH).
		Render(leftContent)

	// ── Right pane: WOPR banner + optional result block ────────────────────
	banner := bannerStyle.Render(woprBanner)

	var resultBlock string
	if m.portResult != nil {
		r := m.portResult
		innerW := rightPaneWidth - 4
		divider := resultDividerStyle.Render(strings.Repeat("─", innerW))
		target := dimStyle.Render(fmt.Sprintf(" HOST  %s", r.host))
		port := dimStyle.Render(fmt.Sprintf(" PORT  %s", r.port))
		var verdict string
		if r.open {
			verdict = resultOpenStyle.Render(fmt.Sprintf(" ██ OPEN   (%.0fms)", float64(r.rtt.Milliseconds())))
		} else {
			verdict = resultClosedStyle.Render(" ░░ CLOSED / FILTERED")
		}
		resultBlock = lipgloss.JoinVertical(lipgloss.Left,
			"",
			divider,
			dimStyle.Render(" LAST RESULT"),
			divider,
			target,
			port,
			"",
			verdict,
			divider,
		)
	}

	var dnsBlock string
	if m.dnsResult != nil {
		r := m.dnsResult
		innerW := rightPaneWidth - 4
		divider := resultDividerStyle.Render(strings.Repeat("─", innerW))

		lines := []string{
			"",
			divider,
			dimStyle.Render(" DNS RESULT"),
			divider,
			dimStyle.Render(" QUERY  " + trunc(r.name, 28)),
		}

		for _, res := range r.results {
			if res.cname != "" {
				lines = append(lines,
					resultOpenStyle.Render(" CNAME → "+trunc(res.cname, 27)),
				)
				break
			}
		}

		lines = append(lines, "")
		for _, res := range r.results {
			lines = append(lines, dimStyle.Render(" VIA "+res.server))
			if res.err != nil {
				lines = append(lines, resultClosedStyle.Render("  NO RESPONSE"))
			} else {
				for _, addr := range res.addrs {
					lines = append(lines, dimStyle.Render("  "+trunc(addr, 34)))
				}
			}
			lines = append(lines, "")
		}
		lines = append(lines, divider)
		dnsBlock = lipgloss.JoinVertical(lipgloss.Left, lines...)
	}

	rightContent := lipgloss.JoinVertical(lipgloss.Left, banner, resultBlock, dnsBlock)
	blankLines := paneH - lipgloss.Height(rightContent)
	if blankLines < 0 {
		blankLines = 0
	}
	rightContent = lipgloss.JoinVertical(
		lipgloss.Left,
		rightContent,
		strings.Repeat("\n", blankLines),
	)
	right := rightPaneStyle.
		Height(paneH).
		Render(rightContent)

	// ── Combine side-by-side ────────────────────────────────────────────────
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}
