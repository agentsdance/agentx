package components

import (
	"fmt"

	"github.com/agentsdance/agentx/internal/version"
	"github.com/agentsdance/agentx/ui/theme"
	"github.com/charmbracelet/lipgloss"
)

// HeaderStats contains statistics to display in the header
type HeaderStats struct {
	MCPInstalled int
	MCPTotal     int
	SkillsCount  int
	PluginsCount int
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

	versionStr := ""
	if version.Version != "" && version.Version != "dev" {
		versionStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Padding(0, 1)
		versionStr = versionStyle.Render(version.Version)
	}

	// Build stats string
	statsStr := fmt.Sprintf(
		"MCP: %d/%d  Skills: %d  Plugins: %d  Agents: %d/%d",
		h.stats.MCPInstalled, h.stats.MCPTotal,
		h.stats.SkillsCount,
		h.stats.PluginsCount,
		h.stats.AgentsOnline, h.stats.AgentsTotal,
	)
	stats := statsStyle.Render(statsStr)

	// Calculate spacing
	leftPart := title + subtitle + versionStr
	leftWidth := lipgloss.Width(leftPart)
	statsWidth := lipgloss.Width(stats)

	// Ensure we fill the entire width with the background color
	return lipgloss.NewStyle().
		Width(h.width).
		Background(theme.HeaderBgColor).
		Render(lipgloss.JoinHorizontal(lipgloss.Center,
			leftPart,
			lipgloss.NewStyle().Width(h.width-leftWidth-statsWidth).Render(""),
			stats,
		))
}
