package ui

// Layout contains calculated dimensions for UI components
type Layout struct {
	// Total dimensions
	Width  int
	Height int

	// Header
	HeaderHeight int

	// Tab bar
	TabBarHeight int

	// Main content area
	MainWidth  int
	MainHeight int

	// Sidebar
	SidebarWidth  int
	SidebarHeight int

	// Footer
	FooterHeight int
}

const (
	MinWidth        = 80
	MinHeight       = 24
	HeaderHeight    = 1
	TabBarHeight    = 1
	FooterHeight    = 1
	SidebarMinWidth = 28
	SidebarRatio    = 0.28 // 28% of width for sidebar
)

// CalculateLayout computes the layout based on terminal size
func CalculateLayout(width, height int) Layout {
	layout := Layout{
		Width:        width,
		Height:       height,
		HeaderHeight: HeaderHeight,
		TabBarHeight: TabBarHeight,
		FooterHeight: FooterHeight,
	}

	// Calculate sidebar width (28% of total, min 28 chars)
	layout.SidebarWidth = int(float64(width) * SidebarRatio)
	if layout.SidebarWidth < SidebarMinWidth {
		layout.SidebarWidth = SidebarMinWidth
	}
	// Hide sidebar on very small screens
	if width < MinWidth {
		layout.SidebarWidth = 0
	}

	// Main content dimensions
	layout.MainWidth = width - layout.SidebarWidth
	layout.MainHeight = height - HeaderHeight - TabBarHeight - FooterHeight

	// Sidebar dimensions
	layout.SidebarHeight = layout.MainHeight

	return layout
}
