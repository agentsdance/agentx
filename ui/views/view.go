package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/agentsdance/agentx/ui/components"
)

// View represents a tab view interface
type View interface {
	// Init initializes the view
	Init() tea.Cmd

	// Update handles messages and returns updated view
	Update(msg tea.Msg) (View, tea.Cmd)

	// View renders the view content
	View() string

	// SetDimensions sets the view dimensions
	SetDimensions(width, height int)

	// Title returns the view title
	Title() string

	// ShortHelp returns keyboard shortcuts for footer
	ShortHelp() []components.FooterAction

	// GetSidebarSections returns sidebar content
	GetSidebarSections() []components.SidebarSection

	// Message returns any status message to display
	Message() string
}
