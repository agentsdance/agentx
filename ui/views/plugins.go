package views

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/agentsdance/agentx/internal/agent"
	"github.com/agentsdance/agentx/internal/plugins"
	"github.com/agentsdance/agentx/ui/components"
	"github.com/agentsdance/agentx/ui/theme"
)

// AvailablePlugin represents a plugin that can be installed
type AvailablePlugin struct {
	Name        string
	Description string
	Source      string
	Components  string // Summary like "2 cmd, 1 skill"
}

// Available plugins from registry (fetched on startup)
var AvailablePlugins []AvailablePlugin

// InitAvailablePlugins fetches plugins from the registry
func InitAvailablePlugins() {
	registryPlugins, err := plugins.FetchRegistryWithFallback()
	if err != nil {
		// Use empty list if fetch fails
		AvailablePlugins = []AvailablePlugin{}
		return
	}

	AvailablePlugins = make([]AvailablePlugin, len(registryPlugins))
	for i, p := range registryPlugins {
		components := ""
		if len(p.Components) > 0 {
			components = strings.Join(p.Components, ", ")
		}
		AvailablePlugins[i] = AvailablePlugin{
			Name:        p.Name,
			Description: p.Description,
			Source:      p.Source,
			Components:  components,
		}
	}
}

// AgentPluginStatus represents an agent's plugin installation status
type AgentPluginStatus struct {
	Agent          agent.Agent
	Exists         bool
	SupportsPlugin bool
	PluginStatus   map[string]bool // pluginName -> installed
	PluginError    map[string]error
}

// PluginsView displays plugins installation status across code agents
type PluginsView struct {
	agents    []AgentPluginStatus
	cursorRow int // Plugin row
	cursorCol int // Agent column
	width     int
	height    int
	message   string
}

// NewPluginsView creates a new plugins view
func NewPluginsView() *PluginsView {
	// Initialize available plugins from registry
	InitAvailablePlugins()

	agents := agent.GetAllAgents()
	statuses := make([]AgentPluginStatus, len(agents))

	for i, a := range agents {
		statuses[i] = AgentPluginStatus{
			Agent:          a,
			Exists:         a.Exists(),
			SupportsPlugin: a.SupportsPlugins(),
			PluginStatus:   make(map[string]bool),
			PluginError:    make(map[string]error),
		}

		// Check each plugin
		for _, plugin := range AvailablePlugins {
			installed, err := a.HasPlugin(plugin.Name)
			statuses[i].PluginStatus[plugin.Name] = installed
			if err != nil {
				statuses[i].PluginError[plugin.Name] = err
			}
		}
	}

	return &PluginsView{
		agents:    statuses,
		cursorRow: 0,
		cursorCol: 0,
	}
}

func (v *PluginsView) Init() tea.Cmd {
	return nil
}

func (v *PluginsView) Update(msg tea.Msg) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if v.cursorRow > 0 {
				v.cursorRow--
			}
		case "down", "j":
			maxRows := len(AvailablePlugins)
			if maxRows == 0 {
				maxRows = 1
			}
			if v.cursorRow < maxRows-1 {
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
			v.installAllForSelectedPlugin()
		case "r":
			v.removeSelected()
		case "c":
			v.refreshStatus()
			v.message = "Status refreshed"
		}
	}
	return v, nil
}

func (v *PluginsView) installSelected() {
	if len(AvailablePlugins) == 0 {
		v.message = "No plugins available in registry"
		return
	}

	status := &v.agents[v.cursorCol]
	agentName := status.Agent.Name()
	plugin := AvailablePlugins[v.cursorRow]

	if !status.SupportsPlugin {
		v.message = fmt.Sprintf("%s doesn't support plugins", agentName)
		return
	}

	if status.PluginStatus[plugin.Name] {
		v.message = fmt.Sprintf("%s already has %s", agentName, plugin.Name)
		return
	}

	if err := status.Agent.InstallPlugin(plugin.Name, plugin.Source); err != nil {
		v.message = fmt.Sprintf("Failed to install %s: %v", plugin.Name, err)
	} else {
		v.message = fmt.Sprintf("Installed %s to %s", plugin.Name, agentName)
		v.refreshStatus()
	}
}

func (v *PluginsView) installAllForSelectedPlugin() {
	if len(AvailablePlugins) == 0 {
		v.message = "No plugins available in registry"
		return
	}

	installed := 0
	plugin := AvailablePlugins[v.cursorRow]

	for i := range v.agents {
		if !v.agents[i].SupportsPlugin {
			continue
		}
		if !v.agents[i].PluginStatus[plugin.Name] {
			if err := v.agents[i].Agent.InstallPlugin(plugin.Name, plugin.Source); err == nil {
				installed++
			}
		}
	}
	v.refreshStatus()
	v.message = fmt.Sprintf("Installed %s to %d agent(s)", plugin.Name, installed)
}

func (v *PluginsView) removeSelected() {
	if len(AvailablePlugins) == 0 {
		v.message = "No plugins available in registry"
		return
	}

	status := &v.agents[v.cursorCol]
	agentName := status.Agent.Name()
	plugin := AvailablePlugins[v.cursorRow]

	if !status.SupportsPlugin {
		v.message = fmt.Sprintf("%s doesn't support plugins", agentName)
		return
	}

	if !status.PluginStatus[plugin.Name] {
		v.message = fmt.Sprintf("%s doesn't have %s", agentName, plugin.Name)
		return
	}

	if err := status.Agent.RemovePlugin(plugin.Name); err != nil {
		v.message = fmt.Sprintf("Failed to remove %s: %v", plugin.Name, err)
	} else {
		v.message = fmt.Sprintf("Removed %s from %s", plugin.Name, agentName)
		v.refreshStatus()
	}
}

func (v *PluginsView) View() string {
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
	b.WriteString(headerStyle.Render("  Plugins Status"))
	b.WriteString("\n")
	b.WriteString(borderStyle.Render("  " + strings.Repeat("-", 70)))
	b.WriteString("\n")

	if len(AvailablePlugins) == 0 {
		b.WriteString("\n")
		b.WriteString("  No plugins available in registry.\n")
		b.WriteString("  Install plugins manually using: agentx plugins install <source>\n")
		b.WriteString("\n")
		return b.String()
	}

	// Column headers (Agent names)
	b.WriteString("                  ")
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

	// Plugin rows
	for pluginIdx, plugin := range AvailablePlugins {
		var row strings.Builder

		// Row cursor
		if pluginIdx == v.cursorRow {
			row.WriteString("> ")
		} else {
			row.WriteString("  ")
		}

		// Plugin name
		name := plugin.Name
		if len(name) > 14 {
			name = name[:14]
		}
		row.WriteString(fmt.Sprintf("%-14s", name))

		// Status for each agent
		for agentIdx, status := range v.agents {
			var cellContent string
			var style lipgloss.Style

			if !status.SupportsPlugin {
				cellContent = "o n/a"
				style = notInstalledStyle
			} else if err := status.PluginError[plugin.Name]; err != nil {
				cellContent = "x error"
				style = errorStyle
			} else if status.PluginStatus[plugin.Name] {
				cellContent = "* installed"
				style = installedStyle
			} else if !status.Exists {
				cellContent = "o n/a"
				style = notInstalledStyle
			} else {
				cellContent = "o ---"
				style = notInstalledStyle
			}

			if pluginIdx == v.cursorRow && agentIdx == v.cursorCol {
				row.WriteString(cursorCellStyle.Render(cellContent))
			} else {
				row.WriteString(style.Render(cellContent))
			}
		}

		// Apply row style
		if pluginIdx == v.cursorRow {
			b.WriteString(selectedRowStyle.Render(row.String()))
		} else {
			b.WriteString(row.String())
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (v *PluginsView) SetDimensions(width, height int) {
	v.width = width
	v.height = height
}

func (v *PluginsView) Title() string {
	return "Plugins"
}

func (v *PluginsView) ShortHelp() []components.FooterAction {
	return []components.FooterAction{
		{Key: "i/Enter", Label: "install"},
		{Key: "I", Label: "install all"},
		{Key: "r", Label: "remove"},
		{Key: "arrows", Label: "navigate"},
		{Key: "c", Label: "refresh"},
		{Key: "q", Label: "quit"},
	}
}

func (v *PluginsView) GetSidebarSections() []components.SidebarSection {
	var installedPlugins []string

	for _, plugin := range AvailablePlugins {
		for _, status := range v.agents {
			if status.PluginStatus[plugin.Name] {
				installedPlugins = append(installedPlugins, plugin.Name)
				break
			}
		}
	}

	return []components.SidebarSection{
		{Title: "Installed Plugins", Items: installedPlugins},
	}
}

func (v *PluginsView) Message() string {
	return v.message
}

func (v *PluginsView) refreshStatus() {
	for i := range v.agents {
		v.agents[i].Exists = v.agents[i].Agent.Exists()
		v.agents[i].SupportsPlugin = v.agents[i].Agent.SupportsPlugins()

		for _, plugin := range AvailablePlugins {
			installed, err := v.agents[i].Agent.HasPlugin(plugin.Name)
			v.agents[i].PluginStatus[plugin.Name] = installed
			if err != nil {
				v.agents[i].PluginError[plugin.Name] = err
			} else {
				delete(v.agents[i].PluginError, plugin.Name)
			}
		}
	}
}

// GetInstalledCount returns total plugin installations across all agents
func (v *PluginsView) GetInstalledCount() int {
	count := 0
	for _, status := range v.agents {
		for _, installed := range status.PluginStatus {
			if installed {
				count++
			}
		}
	}
	return count
}

// GetTotalCount returns total possible installations (agents that support plugins x plugins)
func (v *PluginsView) GetTotalCount() int {
	supportingAgents := 0
	for _, status := range v.agents {
		if status.SupportsPlugin {
			supportingAgents++
		}
	}
	return supportingAgents * len(AvailablePlugins)
}

// GetPluginsCount returns the total number of available plugins
func (v *PluginsView) GetPluginsCount() int {
	return len(AvailablePlugins)
}
