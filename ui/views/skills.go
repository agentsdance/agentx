package views

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/agentsdance/agentx/internal/agent"
	"github.com/agentsdance/agentx/internal/skills"
	"github.com/agentsdance/agentx/ui/components"
	"github.com/agentsdance/agentx/ui/theme"
)

// AvailableSkill represents a skill that can be installed
type AvailableSkill struct {
	Name        string
	Description string
	Source      string // GitHub tree URL or repo#fragment
}

// AvailableSkills is the list of skills from the registry
var AvailableSkills []AvailableSkill

// InitAvailableSkills fetches skills from the registry
func InitAvailableSkills() {
	registrySkills, err := skills.FetchSkillsRegistryWithFallback()
	if err != nil {
		// Use empty list if fetch fails
		AvailableSkills = []AvailableSkill{}
		return
	}

	AvailableSkills = make([]AvailableSkill, len(registrySkills))
	for i, s := range registrySkills {
		// Convert source to installation URL
		source := s.Source
		if source == "local" {
			// For local skills in agentsdance/agentskills repo
			source = fmt.Sprintf("https://github.com/agentsdance/agentskills/tree/master/skills/%s", s.Name)
		}

		// Truncate description for display
		desc := s.Description
		if len(desc) > 40 {
			desc = desc[:37] + "..."
		}

		AvailableSkills[i] = AvailableSkill{
			Name:        s.Name,
			Description: desc,
			Source:      source,
		}
	}
}

// AgentSkillStatus represents an agent's skill installation status
type AgentSkillStatus struct {
	Agent         agent.Agent
	Exists        bool
	SupportsSkill bool
	SkillStatus   map[string]bool // skillName -> installed
	SkillError    map[string]error
}

// SkillsView displays skills installation status across code agents
type SkillsView struct {
	agents    []AgentSkillStatus
	cursorRow int // Skill row
	cursorCol int // Agent column
	width     int
	height    int
	message   string
}

// NewSkillsView creates a new skills view
func NewSkillsView() *SkillsView {
	// Initialize available skills from registry
	InitAvailableSkills()

	agents := agent.GetAllAgents()
	statuses := make([]AgentSkillStatus, len(agents))

	for i, a := range agents {
		statuses[i] = AgentSkillStatus{
			Agent:         a,
			Exists:        a.Exists(),
			SupportsSkill: a.SupportsSkills(),
			SkillStatus:   make(map[string]bool),
			SkillError:    make(map[string]error),
		}

		// Check each skill
		for _, skill := range AvailableSkills {
			installed, err := a.HasSkill(skill.Name)
			statuses[i].SkillStatus[skill.Name] = installed
			if err != nil {
				statuses[i].SkillError[skill.Name] = err
			}
		}
	}

	return &SkillsView{
		agents:    statuses,
		cursorRow: 0,
		cursorCol: 0,
	}
}

func (v *SkillsView) Init() tea.Cmd {
	return nil
}

func (v *SkillsView) Update(msg tea.Msg) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if v.cursorRow > 0 {
				v.cursorRow--
			}
		case "down", "j":
			if v.cursorRow < len(AvailableSkills)-1 {
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
			v.installAllForSelectedSkill()
		case "r":
			v.removeSelected()
		case "c":
			v.refreshStatus()
			v.message = "Status refreshed"
		case "R":
			// Force refresh skills from registry
			InitAvailableSkills()
			v.refreshStatus()
			v.message = fmt.Sprintf("Refreshed %d skills from registry", len(AvailableSkills))
		}
	}
	return v, nil
}

func (v *SkillsView) installSelected() {
	if len(AvailableSkills) == 0 {
		v.message = "No skills available"
		return
	}

	status := &v.agents[v.cursorCol]
	agentName := status.Agent.Name()
	skill := AvailableSkills[v.cursorRow]

	if !status.SupportsSkill {
		v.message = fmt.Sprintf("%s doesn't support skills", agentName)
		return
	}

	if status.SkillStatus[skill.Name] {
		v.message = fmt.Sprintf("%s already has %s", agentName, skill.Name)
		return
	}

	if err := status.Agent.InstallSkill(skill.Name, skill.Source); err != nil {
		v.message = fmt.Sprintf("Failed to install %s: %v", skill.Name, err)
	} else {
		v.message = fmt.Sprintf("Installed %s to %s", skill.Name, agentName)
		v.refreshStatus()
	}
}

func (v *SkillsView) installAllForSelectedSkill() {
	if len(AvailableSkills) == 0 {
		v.message = "No skills available"
		return
	}

	installed := 0
	skill := AvailableSkills[v.cursorRow]

	for i := range v.agents {
		if !v.agents[i].SupportsSkill {
			continue
		}
		if !v.agents[i].SkillStatus[skill.Name] {
			if err := v.agents[i].Agent.InstallSkill(skill.Name, skill.Source); err == nil {
				installed++
			}
		}
	}
	v.refreshStatus()
	v.message = fmt.Sprintf("Installed %s to %d agent(s)", skill.Name, installed)
}

func (v *SkillsView) removeSelected() {
	if len(AvailableSkills) == 0 {
		v.message = "No skills available"
		return
	}

	status := &v.agents[v.cursorCol]
	agentName := status.Agent.Name()
	skill := AvailableSkills[v.cursorRow]

	if !status.SupportsSkill {
		v.message = fmt.Sprintf("%s doesn't support skills", agentName)
		return
	}

	if !status.SkillStatus[skill.Name] {
		v.message = fmt.Sprintf("%s doesn't have %s", agentName, skill.Name)
		return
	}

	if err := status.Agent.RemoveSkill(skill.Name); err != nil {
		v.message = fmt.Sprintf("Failed to remove %s: %v", skill.Name, err)
	} else {
		v.message = fmt.Sprintf("Removed %s from %s", skill.Name, agentName)
		v.refreshStatus()
	}
}

func (v *SkillsView) View() string {
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
	b.WriteString(headerStyle.Render("  Skills Status"))
	b.WriteString("\n")
	b.WriteString(borderStyle.Render("  " + strings.Repeat("─", 70)))
	b.WriteString("\n")

	// Check if skills are available
	if len(AvailableSkills) == 0 {
		b.WriteString("\n  No skills available. Press 'R' to refresh from registry.\n")
		return b.String()
	}

	// Column headers (Agent names)
	b.WriteString("                              ")
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

	// Skill rows
	for skillIdx, skill := range AvailableSkills {
		var row strings.Builder

		// Row cursor
		if skillIdx == v.cursorRow {
			row.WriteString("▸ ")
		} else {
			row.WriteString("  ")
		}

		// Skill name - width 28 characters for long names
		name := skill.Name
		if len(name) > 28 {
			name = name[:28]
		}
		row.WriteString(fmt.Sprintf("%-28s", name))

		// Status for each agent
		for agentIdx, status := range v.agents {
			var cellContent string
			var style lipgloss.Style

			if !status.SupportsSkill {
				cellContent = "○ n/a"
				style = notInstalledStyle
			} else if err := status.SkillError[skill.Name]; err != nil {
				cellContent = "✗ error"
				style = errorStyle
			} else if status.SkillStatus[skill.Name] {
				cellContent = "✓ installed"
				style = installedStyle
			} else if !status.Exists {
				cellContent = "○ n/a"
				style = notInstalledStyle
			} else {
				cellContent = "○ ---"
				style = notInstalledStyle
			}

			if skillIdx == v.cursorRow && agentIdx == v.cursorCol {
				row.WriteString(cursorCellStyle.Render(cellContent))
			} else {
				row.WriteString(style.Render(cellContent))
			}
		}

		// Apply row style
		if skillIdx == v.cursorRow {
			b.WriteString(selectedRowStyle.Render(row.String()))
		} else {
			b.WriteString(row.String())
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (v *SkillsView) SetDimensions(width, height int) {
	v.width = width
	v.height = height
}

func (v *SkillsView) Title() string {
	return "Skills"
}

func (v *SkillsView) ShortHelp() []components.FooterAction {
	return []components.FooterAction{
		{Key: "i/↵", Label: "install"},
		{Key: "I", Label: "install all"},
		{Key: "r", Label: "remove"},
		{Key: "←→", Label: "select agent"},
		{Key: "↑↓", Label: "select skill"},
		{Key: "c", Label: "check"},
		{Key: "R", Label: "refresh"},
		{Key: "q", Label: "quit"},
	}
}

func (v *SkillsView) GetSidebarSections() []components.SidebarSection {
	var installedSkills []string

	for _, skill := range AvailableSkills {
		for _, status := range v.agents {
			if status.SkillStatus[skill.Name] {
				installedSkills = append(installedSkills, skill.Name)
				break
			}
		}
	}

	return []components.SidebarSection{
		{Title: "Installed Skills", Items: installedSkills},
	}
}

func (v *SkillsView) Message() string {
	return v.message
}

func (v *SkillsView) refreshStatus() {
	for i := range v.agents {
		v.agents[i].Exists = v.agents[i].Agent.Exists()
		v.agents[i].SupportsSkill = v.agents[i].Agent.SupportsSkills()

		for _, skill := range AvailableSkills {
			installed, err := v.agents[i].Agent.HasSkill(skill.Name)
			v.agents[i].SkillStatus[skill.Name] = installed
			if err != nil {
				v.agents[i].SkillError[skill.Name] = err
			} else {
				delete(v.agents[i].SkillError, skill.Name)
			}
		}
	}
}

// GetInstalledCount returns total skill installations across all agents
func (v *SkillsView) GetInstalledCount() int {
	count := 0
	for _, status := range v.agents {
		for _, installed := range status.SkillStatus {
			if installed {
				count++
			}
		}
	}
	return count
}

// GetTotalCount returns total possible installations (agents that support skills × skills)
func (v *SkillsView) GetTotalCount() int {
	supportingAgents := 0
	for _, status := range v.agents {
		if status.SupportsSkill {
			supportingAgents++
		}
	}
	return supportingAgents * len(AvailableSkills)
}

// GetSkillsCount returns the total number of available skills
func (v *SkillsView) GetSkillsCount() int {
	return len(AvailableSkills)
}
