package components

import (
	"github.com/charmbracelet/lipgloss"
)

// FooterAction represents a keyboard shortcut in the footer
type FooterAction struct {
	Key   string
	Label string
}

// Footer is the bottom status bar component
type Footer struct {
	actions []FooterAction
	message string
	version string
	width   int
}

// NewFooter creates a new footer with the given actions
func NewFooter(actions []FooterAction) Footer {
	return Footer{actions: actions}
}

// SetWidth sets the footer width
func (f *Footer) SetWidth(width int) {
	f.width = width
}

// SetActions sets the keyboard shortcuts to display
func (f *Footer) SetActions(actions []FooterAction) {
	f.actions = actions
}

// SetMessage sets a temporary message to display
func (f *Footer) SetMessage(message string) {
	f.message = message
}

// SetVersion sets the version to display
func (f *Footer) SetVersion(version string) {
	f.version = version
}

// View renders the footer
func (f Footer) View() string {
	keyStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED"))

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF"))

	separatorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4B5563"))

	versionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280"))

	var parts []string
	for i, action := range f.actions {
		part := keyStyle.Render(action.Key) + labelStyle.Render(" "+action.Label)
		parts = append(parts, part)
		if i < len(f.actions)-1 {
			parts = append(parts, separatorStyle.Render("  |  "))
		}
	}

	actionsContent := lipgloss.JoinHorizontal(lipgloss.Top, parts...)

	// Version on the right
	var content string
	if f.version != "" {
		versionContent := versionStyle.Render(f.version)
		versionWidth := lipgloss.Width(versionContent)
		actionsWidth := lipgloss.Width(actionsContent)
		spacing := f.width - actionsWidth - versionWidth - 4
		if spacing < 0 {
			spacing = 1
		}
		spacer := lipgloss.NewStyle().Width(spacing).Render("")
		content = lipgloss.JoinHorizontal(lipgloss.Top,
			actionsContent, spacer, versionContent)
	} else {
		content = actionsContent
	}

	return lipgloss.NewStyle().
		Width(f.width).
		Background(lipgloss.Color("#111827")).
		Padding(0, 1).
		Render(content)
}
