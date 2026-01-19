package views

import (
	"fmt"
	"strings"

	"github.com/agentsdance/agentx/internal/agent"
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

// Available MCP servers
var MCPServers = []MCPServer{
	{Name: "playwright", Description: "Browser automation"},
	{Name: "context7", Description: "Library documentation"},
	{Name: "remix-icon", Description: "Icon library"},
}

// AgentMCPStatus represents an agent's MCP installation status
type AgentMCPStatus struct {
	Agent           agent.Agent
	Exists          bool
	PlaywrightOK    bool
	Context7OK      bool
	RemixIconOK     bool
	PlaywrightError error
	Context7Error   error
	RemixIconError  error
}

// MCPView displays MCP server installation status across code agents
type MCPView struct {
	agents    []AgentMCPStatus
	cursorRow int // MCP server row
	cursorCol int // Agent column
	width     int
	height    int
	message   string
}

// NewMCPView creates a new MCP view
func NewMCPView() *MCPView {
	agents := agent.GetAllAgents()
	statuses := make([]AgentMCPStatus, len(agents))

	for i, a := range agents {
		pwOK, pwErr := a.HasPlaywright()
		c7OK, c7Err := a.HasContext7()
		riOK, riErr := a.HasRemixIcon()
		statuses[i] = AgentMCPStatus{
			Agent:           a,
			Exists:          a.Exists(),
			PlaywrightOK:    pwOK,
			Context7OK:      c7OK,
			RemixIconOK:     riOK,
			PlaywrightError: pwErr,
			Context7Error:   c7Err,
			RemixIconError:  riErr,
		}
	}

	return &MCPView{
		agents:    statuses,
		cursorRow: 0,
		cursorCol: 0,
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
			if v.cursorRow < len(MCPServers)-1 {
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

func (v *MCPView) installSelected() {
	status := &v.agents[v.cursorCol]
	agentName := status.Agent.Name()

	switch v.cursorRow {
	case 0: // Playwright
		if !status.PlaywrightOK {
			if err := status.Agent.InstallPlaywright(); err != nil {
				v.message = fmt.Sprintf("Failed to install Playwright: %v", err)
			} else {
				v.message = fmt.Sprintf("Installed Playwright to %s", agentName)
				v.refreshStatus()
			}
		} else {
			v.message = fmt.Sprintf("%s already has Playwright", agentName)
		}
	case 1: // Context7
		if !status.Context7OK {
			if err := status.Agent.InstallContext7(); err != nil {
				v.message = fmt.Sprintf("Failed to install Context7: %v", err)
			} else {
				v.message = fmt.Sprintf("Installed Context7 to %s", agentName)
				v.refreshStatus()
			}
		} else {
			v.message = fmt.Sprintf("%s already has Context7", agentName)
		}
	case 2: // Remix Icon
		if !status.RemixIconOK {
			if err := status.Agent.InstallRemixIcon(); err != nil {
				v.message = fmt.Sprintf("Failed to install Remix Icon: %v", err)
			} else {
				v.message = fmt.Sprintf("Installed Remix Icon to %s", agentName)
				v.refreshStatus()
			}
		} else {
			v.message = fmt.Sprintf("%s already has Remix Icon", agentName)
		}
	}
}

func (v *MCPView) installAllForSelectedMCP() {
	installed := 0
	mcpName := MCPServers[v.cursorRow].Name

	for i := range v.agents {
		switch v.cursorRow {
		case 0: // Playwright
			if !v.agents[i].PlaywrightOK {
				if err := v.agents[i].Agent.InstallPlaywright(); err == nil {
					installed++
				}
			}
		case 1: // Context7
			if !v.agents[i].Context7OK {
				if err := v.agents[i].Agent.InstallContext7(); err == nil {
					installed++
				}
			}
		case 2: // Remix Icon
			if !v.agents[i].RemixIconOK {
				if err := v.agents[i].Agent.InstallRemixIcon(); err == nil {
					installed++
				}
			}
		}
	}
	v.refreshStatus()
	v.message = fmt.Sprintf("Installed %s to %d agent(s)", mcpName, installed)
}

func (v *MCPView) removeSelected() {
	status := &v.agents[v.cursorCol]
	agentName := status.Agent.Name()

	switch v.cursorRow {
	case 0: // Playwright
		if status.PlaywrightOK {
			if err := status.Agent.RemovePlaywright(); err != nil {
				v.message = fmt.Sprintf("Failed to remove Playwright: %v", err)
			} else {
				v.message = fmt.Sprintf("Removed Playwright from %s", agentName)
				v.refreshStatus()
			}
		} else {
			v.message = fmt.Sprintf("%s doesn't have Playwright", agentName)
		}
	case 1: // Context7
		if status.Context7OK {
			if err := status.Agent.RemoveContext7(); err != nil {
				v.message = fmt.Sprintf("Failed to remove Context7: %v", err)
			} else {
				v.message = fmt.Sprintf("Removed Context7 from %s", agentName)
				v.refreshStatus()
			}
		} else {
			v.message = fmt.Sprintf("%s doesn't have Context7", agentName)
		}
	case 2: // Remix Icon
		if status.RemixIconOK {
			if err := status.Agent.RemoveRemixIcon(); err != nil {
				v.message = fmt.Sprintf("Failed to remove Remix Icon: %v", err)
			} else {
				v.message = fmt.Sprintf("Removed Remix Icon from %s", agentName)
				v.refreshStatus()
			}
		} else {
			v.message = fmt.Sprintf("%s doesn't have Remix Icon", agentName)
		}
	}
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
		Width(14)

	notInstalledStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Width(14)

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#EF4444")).
		Width(14)

	selectedRowStyle := lipgloss.NewStyle().
		Background(theme.SelectionBgColor)

	cursorCellStyle := lipgloss.NewStyle().
		Background(theme.SelectionBgColor).
		Bold(true).
		Width(14)

	// Header
	b.WriteString(headerStyle.Render("  MCP Server Status"))
	b.WriteString("\n")
	b.WriteString(borderStyle.Render("  " + strings.Repeat("─", 70)))
	b.WriteString("\n")

	// Column headers (Agent names)
	b.WriteString("                ")
	for i, status := range v.agents {
		style := colHeaderStyle
		if i == v.cursorCol {
			style = colHeaderSelectedStyle
		}
		name := status.Agent.Name()
		if len(name) > 12 {
			name = name[:12]
		}
		b.WriteString(style.Width(14).Render(name))
	}
	b.WriteString("\n")

	// MCP server rows
	for mcpIdx, mcp := range MCPServers {
		var row strings.Builder

		// Row cursor
		if mcpIdx == v.cursorRow {
			row.WriteString("▸ ")
		} else {
			row.WriteString("  ")
		}

		// MCP server name
		row.WriteString(fmt.Sprintf("%-12s", mcp.Name))

		// Status for each agent
		for agentIdx, status := range v.agents {
			var installed bool
			var err error
			var notFound bool = !status.Exists

			switch mcpIdx {
			case 0: // Playwright
				installed = status.PlaywrightOK
				err = status.PlaywrightError
			case 1: // Context7
				installed = status.Context7OK
				err = status.Context7Error
			case 2: // Remix Icon
				installed = status.RemixIconOK
				err = status.RemixIconError
			}

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
	var playwrightAgents, context7Agents, remixIconAgents []string

	for _, s := range v.agents {
		if s.PlaywrightOK {
			playwrightAgents = append(playwrightAgents, s.Agent.Name())
		}
		if s.Context7OK {
			context7Agents = append(context7Agents, s.Agent.Name())
		}
		if s.RemixIconOK {
			remixIconAgents = append(remixIconAgents, s.Agent.Name())
		}
	}

	return []components.SidebarSection{
		{Title: "Playwright", Items: playwrightAgents},
		{Title: "Context7", Items: context7Agents},
		{Title: "Remix Icon", Items: remixIconAgents},
	}
}

func (v *MCPView) Message() string {
	return v.message
}

func (v *MCPView) refreshStatus() {
	for i := range v.agents {
		pwOK, pwErr := v.agents[i].Agent.HasPlaywright()
		c7OK, c7Err := v.agents[i].Agent.HasContext7()
		riOK, riErr := v.agents[i].Agent.HasRemixIcon()
		v.agents[i].PlaywrightOK = pwOK
		v.agents[i].Context7OK = c7OK
		v.agents[i].RemixIconOK = riOK
		v.agents[i].PlaywrightError = pwErr
		v.agents[i].Context7Error = c7Err
		v.agents[i].RemixIconError = riErr
		v.agents[i].Exists = v.agents[i].Agent.Exists()
	}
}

// GetInstalledCount returns total MCP installations across all agents
func (v *MCPView) GetInstalledCount() int {
	count := 0
	for _, s := range v.agents {
		if s.PlaywrightOK {
			count++
		}
		if s.Context7OK {
			count++
		}
		if s.RemixIconOK {
			count++
		}
	}
	return count
}

// GetTotalCount returns total possible installations (agents × MCP servers)
func (v *MCPView) GetTotalCount() int {
	return len(v.agents) * len(MCPServers)
}
