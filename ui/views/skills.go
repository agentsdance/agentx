package views

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/agentsdance/agentx/internal/skills"
	"github.com/agentsdance/agentx/ui/components"
)

// SkillsView displays installed skills and commands
type SkillsView struct {
	skills  []skills.Skill
	cursor  int
	width   int
	height  int
	message string
	manager *skills.DefaultSkillManager
}

// NewSkillsView creates a new skills view
func NewSkillsView() *SkillsView {
	mgr := skills.NewSkillManager()
	skillList, _ := mgr.List()

	return &SkillsView{
		skills:  skillList,
		cursor:  0,
		manager: mgr,
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
			if v.cursor > 0 {
				v.cursor--
			}
		case "down", "j":
			if len(v.skills) > 0 && v.cursor < len(v.skills)-1 {
				v.cursor++
			}
		case "r":
			// Remove selected skill
			if len(v.skills) > 0 {
				skill := v.skills[v.cursor]
				if err := v.manager.Remove(skill.Name, skill.Scope); err != nil {
					v.message = fmt.Sprintf("Failed to remove: %v", err)
				} else {
					v.message = fmt.Sprintf("Removed %s", skill.Name)
					v.refreshSkills()
					if v.cursor >= len(v.skills) && v.cursor > 0 {
						v.cursor--
					}
				}
			}
		case "c":
			v.refreshSkills()
			v.message = "Skills refreshed"
		}
	}
	return v, nil
}

func (v *SkillsView) View() string {
	var b strings.Builder

	// Styles
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF"))

	borderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#374151"))

	skillStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7C3AED"))

	commandStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#3B82F6"))

	personalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#10B981"))

	projectStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F59E0B"))

	mutedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280"))

	selectedStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#1F2937"))

	// Header
	b.WriteString(headerStyle.Render("  Installed Skills & Commands"))
	b.WriteString("\n")
	b.WriteString(borderStyle.Render("  " + strings.Repeat("─", 60)))
	b.WriteString("\n")

	if len(v.skills) == 0 {
		b.WriteString("\n")
		b.WriteString(mutedStyle.Render("  No skills installed"))
		b.WriteString("\n\n")
		b.WriteString(mutedStyle.Render("  Install with: agentx skills install <source>"))
		b.WriteString("\n")
		return b.String()
	}

	// Column headers
	b.WriteString(mutedStyle.Render("  NAME              TYPE       SCOPE      DESCRIPTION"))
	b.WriteString("\n")

	// Skill rows
	for i, skill := range v.skills {
		var row strings.Builder

		// Cursor
		if i == v.cursor {
			row.WriteString("▸ ")
		} else {
			row.WriteString("  ")
		}

		// Name
		name := fmt.Sprintf("%-16s", truncate(skill.Name, 16))
		row.WriteString(name)
		row.WriteString("  ")

		// Type with color
		var typeStr string
		if skill.Type == skills.SkillTypeSkill {
			typeStr = skillStyle.Render(fmt.Sprintf("%-9s", "skill"))
		} else {
			typeStr = commandStyle.Render(fmt.Sprintf("%-9s", "command"))
		}
		row.WriteString(typeStr)
		row.WriteString("  ")

		// Scope with color
		var scopeStr string
		if skill.Scope == skills.ScopePersonal {
			scopeStr = personalStyle.Render(fmt.Sprintf("%-9s", "personal"))
		} else {
			scopeStr = projectStyle.Render(fmt.Sprintf("%-9s", "project"))
		}
		row.WriteString(scopeStr)
		row.WriteString("  ")

		// Description
		desc := truncate(skill.Description, 30)
		row.WriteString(mutedStyle.Render(desc))

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

func (v *SkillsView) SetDimensions(width, height int) {
	v.width = width
	v.height = height
}

func (v *SkillsView) Title() string {
	return "Skills"
}

func (v *SkillsView) ShortHelp() []components.FooterAction {
	return []components.FooterAction{
		{Key: "r", Label: "remove"},
		{Key: "c", Label: "check"},
		{Key: "q", Label: "quit"},
	}
}

func (v *SkillsView) GetSidebarSections() []components.SidebarSection {
	var personalSkills, projectSkills, commands []string

	for _, s := range v.skills {
		if s.Type == skills.SkillTypeCommand {
			commands = append(commands, s.Name)
		} else if s.Scope == skills.ScopePersonal {
			personalSkills = append(personalSkills, s.Name)
		} else {
			projectSkills = append(projectSkills, s.Name)
		}
	}

	sections := []components.SidebarSection{}
	if len(personalSkills) > 0 {
		sections = append(sections, components.SidebarSection{Title: "Personal Skills", Items: personalSkills})
	}
	if len(projectSkills) > 0 {
		sections = append(sections, components.SidebarSection{Title: "Project Skills", Items: projectSkills})
	}
	if len(commands) > 0 {
		sections = append(sections, components.SidebarSection{Title: "Commands", Items: commands})
	}

	return sections
}

func (v *SkillsView) Message() string {
	return v.message
}

func (v *SkillsView) refreshSkills() {
	v.skills, _ = v.manager.List()
}

// GetSkillsCount returns the total number of skills
func (v *SkillsView) GetSkillsCount() int {
	return len(v.skills)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
