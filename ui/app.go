package ui

import (
	"context"
	"time"

	"github.com/agentsdance/agentx/internal/version"
	"github.com/agentsdance/agentx/ui/components"
	"github.com/agentsdance/agentx/ui/theme"
	"github.com/agentsdance/agentx/ui/views"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	TabMCP     = 0
	TabSkills  = 1
	TabPlugins = 2
	TabAgents  = 3
)

// AppModel is the main TUI application model
type AppModel struct {
	// Views
	mcpView     *views.MCPView
	skillsView  *views.SkillsView
	pluginsView *views.PluginsView
	agentsView  *views.AgentsView

	// Components
	tabBar  components.TabBar
	header  components.Header
	sidebar components.Sidebar
	footer  components.Footer

	// Layout
	layout    Layout
	activeTab int

	// State
	quitting bool

	// Update state
	updateAvailable bool
	updateNotice    version.UpdateNotice
}

type updateMsg struct {
	notice    version.UpdateNotice
	available bool
	err       error
}

// NewAppModel creates a new TUI application model
func NewAppModel() AppModel {
	// Initialize views
	mcpView := views.NewMCPView()
	skillsView := views.NewSkillsView()
	pluginsView := views.NewPluginsView()
	agentsView := views.NewAgentsView()

	// Initialize tab bar
	tabBar := components.NewTabBar([]components.TabItem{
		{Name: "MCP Servers", Icon: "⬡", Active: true},
		{Name: "Skills", Icon: "🛠", Active: false},
		{Name: "Plugins", Icon: "📦", Active: false},
		{Name: "Code Agents", Icon: ">_", Active: false},
	})

	// Initialize header
	header := components.NewHeader("AgentX")
	header.SetStats(components.HeaderStats{
		MCPInstalled: mcpView.GetInstalledCount(),
		MCPTotal:     mcpView.GetTotalCount(),
		SkillsCount:  skillsView.GetSkillsCount(),
		PluginsCount: pluginsView.GetPluginsCount(),
		AgentsOnline: agentsView.GetOnlineCount(),
		AgentsTotal:  agentsView.GetTotalCount(),
	})

	// Initialize sidebar
	sidebar := components.NewSidebar()
	sidebar.SetSections(mcpView.GetSidebarSections())
	sidebar.SetVersion(version.Version)

	// Initialize footer
	footer := components.NewFooter(mcpView.ShortHelp())

	return AppModel{
		mcpView:     mcpView,
		skillsView:  skillsView,
		pluginsView: pluginsView,
		agentsView:  agentsView,
		tabBar:      tabBar,
		header:      header,
		sidebar:     sidebar,
		footer:      footer,
		activeTab:   TabMCP,
	}
}

// Init implements tea.Model
func (m AppModel) Init() tea.Cmd {
	// Request initial window size
	return tea.Batch(tea.WindowSize(), checkForUpdateCmd())
}

// Update implements tea.Model
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case updateMsg:
		if msg.err == nil && msg.available {
			m.updateAvailable = true
			m.updateNotice = msg.notice
		}

	case tea.WindowSizeMsg:
		m.layout = CalculateLayout(msg.Width, msg.Height)
		m.updateDimensions()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "1":
			m.switchTab(TabMCP)
		case "2":
			m.switchTab(TabSkills)
		case "3":
			m.switchTab(TabPlugins)
		case "4":
			m.switchTab(TabAgents)

		case "tab":
			// Cycle through tabs
			nextTab := (m.activeTab + 1) % 4
			m.switchTab(nextTab)

		case "shift+tab":
			// Cycle backwards
			prevTab := (m.activeTab - 1 + 4) % 4
			m.switchTab(prevTab)

		default:
			// Pass to active view
			m.updateActiveView(msg)
		}
	}

	return m, nil
}

// View implements tea.Model
func (m AppModel) View() string {
	if m.quitting {
		return ""
	}

	mainContent := m.getActiveViewContent()
	sidebarContent := m.sidebar.View()

	// Calculate total height available for the main split area (MainHeight + FooterHeight)
	totalPaneHeight := m.layout.MainHeight + FooterHeight

	// Left Pane (Main View + Keyboard Shortcuts)
	leftActions := m.buildInnerFooter(m.layout.MainWidth - 2)
	leftContentHeight := totalPaneHeight - 1 // 1 line for actions

	leftTop := theme.MainStyle.
		Width(m.layout.MainWidth).
		Height(leftContentHeight).
		MaxHeight(leftContentHeight).
		Padding(1, 2).
		Render(mainContent)

	leftBottom := theme.FooterStyle.
		Width(m.layout.MainWidth).
		Height(1).
		Render(leftActions)

	leftPane := lipgloss.JoinVertical(lipgloss.Left, leftTop, leftBottom)

	// Right Pane (Sidebar + Version)
	rightContentHeight := totalPaneHeight - 1 // 1 line for version
	versionStr := "dev/" + m.sidebar.Version()
	if m.sidebar.Version() == "" {
		versionStr = "dev/v0.0.1"
	}

	statusText := versionStr
	statusColor := lipgloss.Color("#4B5563")
	if m.updateAvailable {
		statusText = "new version " + m.updateNotice.Latest + ", run " + m.updateNotice.Command + " to upgrade"
		statusColor = lipgloss.Color("#F59E0B")
	}

	availableWidth := m.layout.SidebarWidth - 4
	if availableWidth < 0 {
		availableWidth = 0
	}
	statusText = truncateString(statusText, availableWidth)

	rightTop := theme.SidebarStyle.
		Width(m.layout.SidebarWidth).
		Height(rightContentHeight).
		MaxHeight(rightContentHeight).
		Render(sidebarContent)

	rightBottom := lipgloss.NewStyle().
		Background(theme.SidebarBgColor).
		Foreground(statusColor).
		Width(m.layout.SidebarWidth).
		Height(1).
		Padding(0, 2).
		Align(lipgloss.Right).
		Render(statusText)

	rightPane := lipgloss.JoinVertical(lipgloss.Left, rightTop, rightBottom)

	// Join them horizontally
	mainRow := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)

	// Combine everything vertically (Header + TabBar + Split Panes)
	return lipgloss.JoinVertical(lipgloss.Left,
		m.header.View(),
		m.tabBar.View(),
		mainRow,
	)
}

func (m *AppModel) buildInnerFooter(width int) string {
	keyStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED"))

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF"))

	separatorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4B5563"))

	// Get actions from active view
	var actions []components.FooterAction
	switch m.activeTab {
	case TabMCP:
		actions = m.mcpView.ShortHelp()
	case TabSkills:
		actions = m.skillsView.ShortHelp()
	case TabPlugins:
		actions = m.pluginsView.ShortHelp()
	case TabAgents:
		actions = m.agentsView.ShortHelp()
	}

	var parts []string
	for i, action := range actions {
		part := keyStyle.Render(action.Key) + labelStyle.Render(" "+action.Label)
		parts = append(parts, part)
		if i < len(actions)-1 {
			parts = append(parts, separatorStyle.Render("  |  "))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
}

func (m *AppModel) switchTab(tab int) {
	m.activeTab = tab
	m.tabBar.SetActive(tab)

	// Update footer actions and sidebar based on active view
	switch tab {
	case TabMCP:
		m.footer.SetActions(m.mcpView.ShortHelp())
		m.footer.SetMessage(m.mcpView.Message())
		m.sidebar.SetSections(m.mcpView.GetSidebarSections())
	case TabSkills:
		m.footer.SetActions(m.skillsView.ShortHelp())
		m.footer.SetMessage(m.skillsView.Message())
		m.sidebar.SetSections(m.skillsView.GetSidebarSections())
	case TabPlugins:
		m.footer.SetActions(m.pluginsView.ShortHelp())
		m.footer.SetMessage(m.pluginsView.Message())
		m.sidebar.SetSections(m.pluginsView.GetSidebarSections())
	case TabAgents:
		m.footer.SetActions(m.agentsView.ShortHelp())
		m.footer.SetMessage(m.agentsView.Message())
		m.sidebar.SetSections(m.agentsView.GetSidebarSections())
	}
}

func (m *AppModel) updateActiveView(msg tea.Msg) {
	switch m.activeTab {
	case TabMCP:
		m.mcpView.Update(msg)
		m.footer.SetMessage(m.mcpView.Message())
		m.sidebar.SetSections(m.mcpView.GetSidebarSections())
		m.updateHeaderStats()
	case TabSkills:
		m.skillsView.Update(msg)
		m.footer.SetMessage(m.skillsView.Message())
		m.sidebar.SetSections(m.skillsView.GetSidebarSections())
		m.updateHeaderStats()
	case TabPlugins:
		m.pluginsView.Update(msg)
		m.footer.SetMessage(m.pluginsView.Message())
		m.sidebar.SetSections(m.pluginsView.GetSidebarSections())
		m.updateHeaderStats()
	case TabAgents:
		m.agentsView.Update(msg)
		m.footer.SetMessage(m.agentsView.Message())
		m.sidebar.SetSections(m.agentsView.GetSidebarSections())
		m.updateHeaderStats()
	}
}

func (m *AppModel) getActiveViewContent() string {
	switch m.activeTab {
	case TabMCP:
		return m.mcpView.View()
	case TabSkills:
		return m.skillsView.View()
	case TabPlugins:
		return m.pluginsView.View()
	case TabAgents:
		return m.agentsView.View()
	default:
		return ""
	}
}

func (m *AppModel) updateDimensions() {
	m.header.SetWidth(m.layout.Width)
	m.tabBar.SetWidth(m.layout.Width)
	m.sidebar.SetDimensions(m.layout.SidebarWidth, m.layout.SidebarHeight)
	m.footer.SetWidth(m.layout.Width)

	m.mcpView.SetDimensions(m.layout.MainWidth, m.layout.MainHeight)
	m.skillsView.SetDimensions(m.layout.MainWidth, m.layout.MainHeight)
	m.pluginsView.SetDimensions(m.layout.MainWidth, m.layout.MainHeight)
	m.agentsView.SetDimensions(m.layout.MainWidth, m.layout.MainHeight)
}

func (m *AppModel) updateHeaderStats() {
	m.header.SetStats(components.HeaderStats{
		MCPInstalled: m.mcpView.GetInstalledCount(),
		MCPTotal:     m.mcpView.GetTotalCount(),
		SkillsCount:  m.skillsView.GetSkillsCount(),
		PluginsCount: m.pluginsView.GetPluginsCount(),
		AgentsOnline: m.agentsView.GetOnlineCount(),
		AgentsTotal:  m.agentsView.GetTotalCount(),
	})
}

func checkForUpdateCmd() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		notice, available, err := version.CheckForUpdate(ctx)
		return updateMsg{notice: notice, available: available, err: err}
	}
}

func truncateString(value string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	if lipgloss.Width(value) <= maxWidth {
		return value
	}
	runes := []rune(value)
	if len(runes) <= maxWidth {
		return value
	}
	return string(runes[:maxWidth])
}

// Run starts the TUI application
func Run() error {
	p := tea.NewProgram(NewAppModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
