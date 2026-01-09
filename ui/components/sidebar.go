package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// SidebarSection represents a section in the sidebar
type SidebarSection struct {
	Title string
	Items []string
}

// Sidebar is a side panel component
type Sidebar struct {
	sections []SidebarSection
	version  string
	width    int
	height   int
	focused  bool
}

// NewSidebar creates a new sidebar
func NewSidebar() Sidebar {
	return Sidebar{}
}

// SetDimensions sets the sidebar dimensions
func (s *Sidebar) SetDimensions(width, height int) {
	s.width = width
	s.height = height
}

// SetFocused sets whether the sidebar is focused
func (s *Sidebar) SetFocused(focused bool) {
	s.focused = focused
}

// SetSections sets the sidebar sections
func (s *Sidebar) SetSections(sections []SidebarSection) {
	s.sections = sections
}

// SetVersion sets the version to display at the bottom
func (s *Sidebar) SetVersion(version string) {
	s.version = version
}

// Version returns the version string
func (s Sidebar) Version() string {
	return s.version
}

// View renders the sidebar content only
func (s Sidebar) View() string {
	if s.width == 0 {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#E5E7EB")).
		MarginBottom(1)

	itemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF")).
		PaddingLeft(1)

	var sections []string
	for _, section := range s.sections {
		var items []string
		items = append(items, titleStyle.Render(section.Title))
		for _, item := range section.Items {
			items = append(items, itemStyle.Render(item))
		}
		sections = append(sections, strings.Join(items, "\n"))
	}

	return strings.Join(sections, "\n\n")
}
