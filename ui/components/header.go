package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// HeaderStats contains statistics to display in the header
type HeaderStats struct {
	MCPInstalled int
	MCPTotal     int
	SkillsCount  int
	AgentsOnline int
	AgentsTotal  int
}

// Header is the top header component
type Header struct {
	title string
	stats HeaderStats
	width int
}

// NewHeader creates a new header with the given title
func NewHeader(title string) Header {
	return Header{title: title}
}

// SetStats sets the header statistics
func (h *Header) SetStats(stats HeaderStats) {
	h.stats = stats
}

// SetWidth sets the header width
func (h *Header) SetWidth(width int) {
	h.width = width
}

// View renders the header
func (h Header) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED")).
		Padding(0, 1)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Padding(0, 1)

	statsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF"))

	title := titleStyle.Render(h.title)
	subtitle := subtitleStyle.Render("Agent Extension: MCP Servers & Agent Skills Manager")

	// Build stats string
	statsStr := fmt.Sprintf(
		"MCP: %d/%d  Skills: %d  Agents: %d/%d",
		h.stats.MCPInstalled, h.stats.MCPTotal,
		h.stats.SkillsCount,
		h.stats.AgentsOnline, h.stats.AgentsTotal,
	)
	stats := statsStyle.Render(statsStr)

	// Calculate spacing
	leftPart := title + subtitle
	leftWidth := lipgloss.Width(leftPart)
	statsWidth := lipgloss.Width(stats)
	spacing := h.width - leftWidth - statsWidth - 2
	if spacing < 0 {
		spacing = 1
	}

	spacer := lipgloss.NewStyle().Width(spacing).Render("")

	row := lipgloss.JoinHorizontal(lipgloss.Center, leftPart, spacer, stats)

	return lipgloss.NewStyle().
		Width(h.width).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(lipgloss.Color("#374151")).
		Render(row)
}
