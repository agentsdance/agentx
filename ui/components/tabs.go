package components

import (
	"github.com/charmbracelet/lipgloss"
)

// TabItem represents a single tab
type TabItem struct {
	Name   string
	Icon   string
	Active bool
}

// TabBar is a horizontal tab bar component
type TabBar struct {
	tabs   []TabItem
	active int
	width  int
}

// NewTabBar creates a new tab bar with the given tabs
func NewTabBar(tabs []TabItem) TabBar {
	return TabBar{tabs: tabs, active: 0}
}

// SetActive sets the active tab by index
func (t *TabBar) SetActive(index int) {
	if index >= 0 && index < len(t.tabs) {
		t.active = index
		for i := range t.tabs {
			t.tabs[i].Active = (i == index)
		}
	}
}

// SetWidth sets the tab bar width
func (t *TabBar) SetWidth(width int) {
	t.width = width
}

// Active returns the active tab index
func (t *TabBar) Active() int {
	return t.active
}

// View renders the tab bar
func (t TabBar) View() string {
	activeStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#7C3AED")).
		Padding(0, 2)

	inactiveStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF")).
		Padding(0, 2)

	var tabs []string
	for _, tab := range t.tabs {
		label := tab.Icon + " " + tab.Name
		if tab.Active {
			tabs = append(tabs, activeStyle.Render(label))
		} else {
			tabs = append(tabs, inactiveStyle.Render(label))
		}
	}

	content := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)

	return lipgloss.NewStyle().
		Width(t.width).
		Background(lipgloss.Color("#1F2937")).
		Render(content)
}
