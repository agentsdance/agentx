package views

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/agentsdance/agentx/internal/agent"
	"github.com/agentsdance/agentx/ui/components"
	"github.com/agentsdance/agentx/ui/theme"
)

// CodeAgentInfo represents information about a code agent
type CodeAgentInfo struct {
	Agent      agent.Agent
	Exists     bool
	ConfigPath string
}

// AgentsView displays code agent status
type AgentsView struct {
	agents  []CodeAgentInfo
	cursor  int
	width   int
	height  int
	message string
}

// NewAgentsView creates a new agents view
func NewAgentsView() *AgentsView {
	allAgents := agent.GetAllAgents()
	infos := make([]CodeAgentInfo, len(allAgents))

	for i, a := range allAgents {
		infos[i] = CodeAgentInfo{
			Agent:      a,
			Exists:     a.Exists(),
			ConfigPath: a.ConfigPath(),
		}
	}

	return &AgentsView{
		agents: infos,
		cursor: 0,
	}
}

func (v *AgentsView) Init() tea.Cmd {
	return nil
}

func (v *AgentsView) Update(msg tea.Msg) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if v.cursor > 0 {
				v.cursor--
			}
		case "down", "j":
			if v.cursor < len(v.agents)-1 {
				v.cursor++
			}
		case "c":
			v.refreshStatus()
			v.message = "Status refreshed"
		case "o":
			// Open config in editor (placeholder for future)
			if len(v.agents) > 0 {
				v.message = fmt.Sprintf("Config: %s", v.agents[v.cursor].ConfigPath)
			}
		}
	}
	return v, nil
}

func (v *AgentsView) View() string {
	var b strings.Builder

	// Styles
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF"))

	borderStyle := lipgloss.NewStyle().
		Foreground(theme.SidebarBgColor)

	activeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#10B981"))

	inactiveStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280"))

	mutedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280"))

	selectedStyle := lipgloss.NewStyle().
		Background(theme.SelectionBgColor)

	// Header
	b.WriteString(headerStyle.Render("  Code Agents"))
	b.WriteString("\n")
	b.WriteString(borderStyle.Render("  " + strings.Repeat("─", 60)))
	b.WriteString("\n")

	// Agent rows
	for i, info := range v.agents {
		var row strings.Builder

		// Cursor
		if i == v.cursor {
			row.WriteString("▸ ")
		} else {
			row.WriteString("  ")
		}

		// Agent name
		name := fmt.Sprintf("%-14s", info.Agent.Name())
		row.WriteString(name)
		row.WriteString("  ")

		// Status
		var statusStr string
		if info.Exists {
			statusStr = activeStyle.Render("● configured")
		} else {
			statusStr = inactiveStyle.Render("○ not found")
		}
		row.WriteString(fmt.Sprintf("%-16s", statusStr))
		row.WriteString("  ")

		// Config path
		row.WriteString(mutedStyle.Render(info.ConfigPath))

		// Apply row style
		if i == v.cursor {
			b.WriteString(selectedStyle.Render(row.String()))
		} else {
			b.WriteString(row.String())
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (v *AgentsView) SetDimensions(width, height int) {
	v.width = width
	v.height = height
}

func (v *AgentsView) Title() string {
	return "Code Agents"
}

func (v *AgentsView) ShortHelp() []components.FooterAction {
	return []components.FooterAction{
		{Key: "c", Label: "check"},
		{Key: "o", Label: "show config"},
		{Key: "q", Label: "quit"},
	}
}

func (v *AgentsView) GetSidebarSections() []components.SidebarSection {
	var configured, notFound []string
	for _, info := range v.agents {
		if info.Exists {
			configured = append(configured, info.Agent.Name())
		} else {
			notFound = append(notFound, info.Agent.Name())
		}
	}

	return []components.SidebarSection{
		{Title: "Configured", Items: configured},
		{Title: "Not Found", Items: notFound},
	}
}

func (v *AgentsView) Message() string {
	return v.message
}

func (v *AgentsView) refreshStatus() {
	for i := range v.agents {
		v.agents[i].Exists = v.agents[i].Agent.Exists()
	}
}

// GetOnlineCount returns number of configured agents
func (v *AgentsView) GetOnlineCount() int {
	count := 0
	for _, info := range v.agents {
		if info.Exists {
			count++
		}
	}
	return count
}

// GetTotalCount returns total number of agents
func (v *AgentsView) GetTotalCount() int {
	return len(v.agents)
}
