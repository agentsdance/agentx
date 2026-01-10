package theme

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	primaryColor   = lipgloss.Color("#7C3AED") // Purple
	successColor   = lipgloss.Color("#10B981") // Green
	warningColor   = lipgloss.Color("#F59E0B") // Amber
	errorColor     = lipgloss.Color("#EF4444") // Red
	mutedColor     = lipgloss.Color("#6B7280") // Gray
	borderColor    = lipgloss.Color("#374151") // Dark gray
	highlightColor = lipgloss.Color("#8B5CF6") // Light purple

	// Backgrounds
	MainBgColor    = lipgloss.Color("#000000") // Black
	SidebarBgColor = lipgloss.Color("#121212") // Dark Gray
	FooterBgColor  = lipgloss.Color("#000000")
	HeaderBgColor    = lipgloss.Color("#1a1b26") // Deep Tokyo Night Blue
	TabBarBgColor    = lipgloss.Color("#1a1b26")
	SelectionBgColor = lipgloss.Color("#3b4261") // Slate Blue for selection

	// Title style
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor)

	// Sidebar styles
	SidebarTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#E5E7EB")).
				MarginBottom(1)

	SidebarItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#9CA3AF")).
				PaddingLeft(1)

	VersionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4B5563"))

	// Layout Styles
	MainStyle = lipgloss.NewStyle().
			Background(MainBgColor)

	SidebarStyle = lipgloss.NewStyle().
			Background(SidebarBgColor).
			Padding(1, 2)

	FooterStyle = lipgloss.NewStyle().
			Background(FooterBgColor).
			Padding(0, 1)
)
