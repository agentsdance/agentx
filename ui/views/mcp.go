package views

import (
	"fmt"
	"sort"
	"strings"

	"github.com/agentsdance/agentx/internal/agent"
	"github.com/agentsdance/agentx/internal/config"
	"github.com/agentsdance/agentx/ui/components"
	"github.com/agentsdance/agentx/ui/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MCPServer represents an MCP server type
type MCPServer struct {
	Name        string
	Description string
}

const (
	serverNameWidth = 14
	cellWidth       = 14
	rowPrefixWidth  = 2
)

// Available MCP servers
var embeddedMCPServers = []MCPServer{
	{Name: "playwright", Description: "Browser automation"},
	{Name: "context7", Description: "Library documentation"},
	{Name: "remix-icon", Description: "Icon library"},
}

// AgentMCPStatus represents an agent's MCP installation status
type AgentMCPStatus struct {
	Agent     agent.Agent
	Exists    bool
	Installed map[string]bool
	Errors    map[string]error
}

// MCPView displays MCP server installation status across code agents
type MCPView struct {
	agents        []AgentMCPStatus
	servers       []MCPServer
	serverConfigs map[string]agent.MCPConfigEntry
	cursorRow     int // MCP server row
	cursorCol     int // Agent column
	width         int
	height        int
	message       string
}

// NewMCPView creates a new MCP view
func NewMCPView() *MCPView {
	agents := agent.GetAllAgents()
	serverConfigs := agent.CollectMCPConfigs(agents)
	servers := buildMCPServerList(serverConfigs)
	statuses := make([]AgentMCPStatus, len(agents))

	for i, a := range agents {
		installed := make(map[string]bool)
		errors := make(map[string]error)
		for _, srv := range servers {
			ok, err := a.HasMCP(srv.Name)
			installed[srv.Name] = ok
			errors[srv.Name] = err
		}
		statuses[i] = AgentMCPStatus{
			Agent:     a,
			Exists:    a.Exists(),
			Installed: installed,
			Errors:    errors,
		}
	}

	return &MCPView{
		agents:        statuses,
		servers:       servers,
		serverConfigs: serverConfigs,
		cursorRow:     0,
		cursorCol:     0,
	}
}

func (v *MCPView) Init() tea.Cmd {
	return nil
}

func (v *MCPView) Update(msg tea.Msg) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if v.cursorRow > 0 {
				v.cursorRow--
			}
		case "down", "j":
			if v.cursorRow < len(v.servers)-1 {
				v.cursorRow++
			}
		case "left", "h":
			if v.cursorCol > 0 {
				v.cursorCol--
			}
		case "right", "l":
			if v.cursorCol < len(v.agents)-1 {
				v.cursorCol++
			}
		case "i", "enter", " ":
			v.installSelected()
		case "I":
			v.installAllForSelectedMCP()
		case "r":
			v.removeSelected()
		case "c":
			v.refreshStatus()
			v.message = "Status refreshed"
		}
	}
	return v, nil
}

func buildMCPServerList(discovered map[string]agent.MCPConfigEntry) []MCPServer {
	servers := make([]MCPServer, 0, len(embeddedMCPServers))
	servers = append(servers, embeddedMCPServers...)

	embedded := map[string]struct{}{}
	for _, srv := range embeddedMCPServers {
		embedded[srv.Name] = struct{}{}
	}

	extraNames := make([]string, 0, len(discovered))
	for name := range discovered {
		if _, ok := embedded[name]; ok {
			continue
		}
		extraNames = append(extraNames, name)
	}
	sort.Strings(extraNames)

	for _, name := range extraNames {
		entry := discovered[name]
		description := "Detected MCP"
		if entry.Source != "" {
			description = fmt.Sprintf("Detected in %s", entry.Source)
		}
		servers = append(servers, MCPServer{
			Name:        name,
			Description: description,
		})
	}

	return servers
}

func (v *MCPView) installSelected() {
	status := &v.agents[v.cursorCol]
	agentName := status.Agent.Name()
	serverName := v.servers[v.cursorRow].Name
	if status.Installed[serverName] {
		v.message = fmt.Sprintf("%s already has %s", agentName, serverName)
		return
	}

	cfg := v.configForServer(serverName)
	if cfg == nil {
		v.message = fmt.Sprintf("No config found for %s", serverName)
		return
	}

	if err := status.Agent.InstallMCP(serverName, cfg); err != nil {
		v.message = fmt.Sprintf("Failed to install %s: %v", serverName, err)
		return
	}
	v.message = fmt.Sprintf("Installed %s to %s", serverName, agentName)
	v.refreshStatus()
}

func (v *MCPView) installAllForSelectedMCP() {
	installed := 0
	mcpName := v.servers[v.cursorRow].Name
	cfg := v.configForServer(mcpName)
	if cfg == nil {
		v.message = fmt.Sprintf("No config found for %s", mcpName)
		return
	}

	for i := range v.agents {
		if v.agents[i].Installed[mcpName] {
			continue
		}
		if err := v.agents[i].Agent.InstallMCP(mcpName, cfg); err == nil {
			installed++
		}
	}
	v.refreshStatus()
	v.message = fmt.Sprintf("Installed %s to %d agent(s)", mcpName, installed)
}

func (v *MCPView) removeSelected() {
	status := &v.agents[v.cursorCol]
	agentName := status.Agent.Name()
	serverName := v.servers[v.cursorRow].Name
	if !status.Installed[serverName] {
		v.message = fmt.Sprintf("%s doesn't have %s", agentName, serverName)
		return
	}
	if err := status.Agent.RemoveMCP(serverName); err != nil {
		v.message = fmt.Sprintf("Failed to remove %s: %v", serverName, err)
		return
	}
	v.message = fmt.Sprintf("Removed %s from %s", serverName, agentName)
	v.refreshStatus()
}

func (v *MCPView) configForServer(name string) map[string]interface{} {
	if entry, ok := v.serverConfigs[name]; ok && entry.Config != nil {
		return entry.Config
	}

	switch name {
	case "playwright":
		return config.PlaywrightMCPConfig
	case "context7":
		return config.Context7MCPConfig
	case "remix-icon":
		return config.RemixIconMCPConfig
	default:
		return nil
	}
}

func displayMCPName(name string) string {
	switch name {
	case "playwright":
		return "Playwright"
	case "context7":
		return "Context7"
	case "remix-icon":
		return "Remix Icon"
	default:
		return name
	}
}

func truncateMCPName(name string, max int) string {
	if max <= 0 {
		return ""
	}
	if lipgloss.Width(name) <= max {
		return name
	}
	runes := []rune(name)
	if len(runes) <= max {
		return name
	}
	return string(runes[:max])
}

func (v *MCPView) View() string {
	var b strings.Builder

	// Styles
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF"))

	borderStyle := lipgloss.NewStyle().
		Foreground(theme.SidebarBgColor)

	colHeaderStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#9CA3AF"))

	colHeaderSelectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.SelectionBgColor)

	installedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#10B981")).
		Width(cellWidth)

	notInstalledStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Width(cellWidth)

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#EF4444")).
		Width(cellWidth)

	selectedRowStyle := lipgloss.NewStyle().
		Background(theme.SelectionBgColor)

	cursorCellStyle := lipgloss.NewStyle().
		Background(theme.SelectionBgColor).
		Bold(true).
		Width(cellWidth)

	// Header
	b.WriteString(headerStyle.Render("  MCP Server Status"))
	b.WriteString("\n")
	b.WriteString(borderStyle.Render("  " + strings.Repeat("─", 70)))
	b.WriteString("\n")

	// Column headers (Agent names)
	b.WriteString(strings.Repeat(" ", serverNameWidth+rowPrefixWidth))
	for i, status := range v.agents {
		style := colHeaderStyle
		if i == v.cursorCol {
			style = colHeaderSelectedStyle
		}
		name := status.Agent.Name()
		if len(name) > 12 {
			name = name[:12]
		}
		b.WriteString(style.Width(cellWidth).Render(name))
	}
	b.WriteString("\n")

	// MCP server rows
	for mcpIdx, mcp := range v.servers {
		var row strings.Builder

		// Row cursor
		if mcpIdx == v.cursorRow {
			row.WriteString("▸ ")
		} else {
			row.WriteString("  ")
		}

		// MCP server name
		name := truncateMCPName(mcp.Name, serverNameWidth)
		row.WriteString(fmt.Sprintf("%-*s", serverNameWidth, name))

		// Status for each agent
		for agentIdx, status := range v.agents {
			var installed bool
			var err error
			var notFound bool = !status.Exists

			installed = status.Installed[mcp.Name]
			err = status.Errors[mcp.Name]

			var cellContent string
			var style lipgloss.Style

			if err != nil {
				cellContent = "✗ error"
				style = errorStyle
			} else if installed {
				cellContent = "✓ installed"
				style = installedStyle
			} else if notFound {
				cellContent = "○ n/a"
				style = notInstalledStyle
			} else {
				cellContent = "○ ---"
				style = notInstalledStyle
			}

			if mcpIdx == v.cursorRow && agentIdx == v.cursorCol {
				row.WriteString(cursorCellStyle.Render(cellContent))
			} else {
				row.WriteString(style.Render(cellContent))
			}
		}

		// Apply row style
		if mcpIdx == v.cursorRow {
			b.WriteString(selectedRowStyle.Render(row.String()))
		} else {
			b.WriteString(row.String())
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (v *MCPView) formatStatus(installed bool, err error, notFound bool,
	installedStyle, notInstalledStyle, errorStyle lipgloss.Style) string {
	if err != nil {
		return errorStyle.Render("✗ error")
	} else if installed {
		return installedStyle.Render("✓ installed")
	} else if notFound {
		return notInstalledStyle.Render("○ n/a")
	}
	return notInstalledStyle.Render("○ ---")
}

func (v *MCPView) SetDimensions(width, height int) {
	v.width = width
	v.height = height
}

func (v *MCPView) Title() string {
	return "MCP Servers"
}

func (v *MCPView) ShortHelp() []components.FooterAction {
	return []components.FooterAction{
		{Key: "i/↵", Label: "install"},
		{Key: "I", Label: "install all"},
		{Key: "r", Label: "remove"},
		{Key: "←→", Label: "select agent"},
		{Key: "↑↓", Label: "select MCP"},
		{Key: "c", Label: "check"},
		{Key: "q", Label: "quit"},
	}
}

func (v *MCPView) GetSidebarSections() []components.SidebarSection {
	sections := make([]components.SidebarSection, 0, len(v.servers))
	for _, server := range v.servers {
		var agents []string
		for _, s := range v.agents {
			if s.Installed[server.Name] {
				agents = append(agents, s.Agent.Name())
			}
		}
		sections = append(sections, components.SidebarSection{
			Title: displayMCPName(server.Name),
			Items: agents,
		})
	}
	return sections
}

func (v *MCPView) Message() string {
	return v.message
}

func (v *MCPView) refreshStatus() {
	agents := make([]agent.Agent, 0, len(v.agents))
	for _, status := range v.agents {
		agents = append(agents, status.Agent)
	}
	v.serverConfigs = agent.CollectMCPConfigs(agents)
	v.servers = buildMCPServerList(v.serverConfigs)
	if v.cursorRow >= len(v.servers) {
		v.cursorRow = len(v.servers) - 1
		if v.cursorRow < 0 {
			v.cursorRow = 0
		}
	}

	for i := range v.agents {
		v.agents[i].Installed = make(map[string]bool)
		v.agents[i].Errors = make(map[string]error)
		for _, server := range v.servers {
			ok, err := v.agents[i].Agent.HasMCP(server.Name)
			v.agents[i].Installed[server.Name] = ok
			v.agents[i].Errors[server.Name] = err
		}
		v.agents[i].Exists = v.agents[i].Agent.Exists()
	}
}

// GetInstalledCount returns total MCP installations across all agents
func (v *MCPView) GetInstalledCount() int {
	count := 0
	for _, s := range v.agents {
		for _, server := range v.servers {
			if s.Installed[server.Name] {
				count++
			}
		}
	}
	return count
}

// GetTotalCount returns total possible installations (agents × MCP servers)
func (v *MCPView) GetTotalCount() int {
	return len(v.agents) * len(v.servers)
}
