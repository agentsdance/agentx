package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/agentsdance/agentx/internal/skills"
	"github.com/spf13/cobra"
)

var skillsScope string
var skillsAgent string

var skillsCmd = &cobra.Command{
	Use:   "skills",
	Short: "Manage Claude Code and Codex skills",
	Long: `Manage Claude Code and Codex skills and slash commands.

Use --agent to switch between agents (default: Claude Code). Codex does not
support command files.

Skills are stored in:
  Personal: ~/.claude/skills/ and ~/.claude/commands/
  Project:  .claude/skills/ and .claude/commands/

Codex skills are stored in:
  Personal: $CODEX_HOME/skills/ (default ~/.codex/skills/)
  Project:  .codex/skills/`,
}

var skillsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed skills",
	Long:  `List all installed skills and commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		mgr, err := resolveSkillsManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		var skillList []skills.Skill

		if skillsScope != "" {
			scope := skills.SkillScope(skillsScope)
			if scope != skills.ScopePersonal && scope != skills.ScopeProject {
				fmt.Fprintf(os.Stderr, "Invalid scope: %s (use 'personal' or 'project')\n", skillsScope)
				os.Exit(1)
			}
			skillList, err = mgr.ListByScope(scope)
		} else {
			skillList, err = mgr.List()
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if len(skillList) == 0 {
			fmt.Println("No skills installed")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tTYPE\tSCOPE\tDESCRIPTION")
		fmt.Fprintln(w, "----\t----\t-----\t-----------")
		for _, s := range skillList {
			desc := s.Description
			if len(desc) > 50 {
				desc = desc[:47] + "..."
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", s.Name, s.Type, s.Scope, desc)
		}
		w.Flush()
	},
}

var skillsInstallCmd = &cobra.Command{
	Use:   "install <source>",
	Short: "Install a skill from URL or path",
	Long: `Install a skill from a local path or Git repository.

Examples:
  agentx skills install ./my-skill/
  agentx skills install ./my-command.md
  agentx skills install https://github.com/user/skills-repo
  agentx skills install https://github.com/user/repo#skill-name

The source can be:
  - A local directory containing SKILL.md (skill)
  - A local .md file (command)
  - A Git repository URL
  - A Git repository URL with #skill-name fragment`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		source := args[0]
		scope := skills.ScopePersonal
		if skillsScope == "project" {
			scope = skills.ScopeProject
		}

		mgr, err := resolveSkillsManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		skill, err := mgr.Install(source, scope)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Installed %s '%s' to %s scope\n", skill.Type, skill.Name, skill.Scope)
		fmt.Printf("  Path: %s\n", skill.Path)
		if skill.Description != "" {
			fmt.Printf("  Description: %s\n", skill.Description)
		}
	},
}

var skillsRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a skill",
	Long:  `Remove a skill or command by name.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		scope := skills.ScopePersonal
		if skillsScope == "project" {
			scope = skills.ScopeProject
		}

		mgr, err := resolveSkillsManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if err := mgr.Remove(name, scope); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Removed skill: %s\n", name)
	},
}

var skillsCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check skills installation status",
	Long:  `Verify that all installed skills are valid and properly configured.`,
	Run: func(cmd *cobra.Command, args []string) {
		mgr, err := resolveSkillsManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		statuses, err := mgr.Check()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if len(statuses) == 0 {
			fmt.Println("No skills installed")
			return
		}

		fmt.Println("Skills Health Check")
		fmt.Println("===================")
		fmt.Println()

		hasIssues := false
		for _, s := range statuses {
			status := "OK"
			statusColor := "\033[32m" // Green
			if !s.Valid {
				status = "ERROR"
				statusColor = "\033[31m" // Red
				hasIssues = true
			} else if len(s.Issues) > 0 {
				status = "WARNING"
				statusColor = "\033[33m" // Yellow
				hasIssues = true
			}
			resetColor := "\033[0m"

			fmt.Printf("%-20s %s%s%s (%s, %s)\n",
				s.Skill.Name,
				statusColor, status, resetColor,
				s.Skill.Type, s.Skill.Scope)

			for _, issue := range s.Issues {
				fmt.Printf("  - %s\n", issue)
			}
			if s.Error != nil {
				fmt.Printf("  - Error: %v\n", s.Error)
			}
		}

		if !hasIssues {
			fmt.Println()
			fmt.Println("All skills are healthy!")
		}
	},
}

func init() {
	skillsCmd.PersistentFlags().StringVarP(&skillsScope, "scope", "s", "",
		"Scope for the operation (personal, project)")
	skillsCmd.PersistentFlags().StringVarP(&skillsAgent, "agent", "a", "claude",
		"Target agent for skills (claude, codex)")

	skillsCmd.AddCommand(skillsListCmd)
	skillsCmd.AddCommand(skillsInstallCmd)
	skillsCmd.AddCommand(skillsRemoveCmd)
	skillsCmd.AddCommand(skillsCheckCmd)
}

func resolveSkillsManager() (*skills.DefaultSkillManager, error) {
	switch normalizeSkillsAgent(skillsAgent) {
	case "", "claude", "claudecode":
		return skills.NewSkillManager(), nil
	case "codex":
		return skills.NewCodexSkillManager(), nil
	default:
		return nil, fmt.Errorf("unknown agent: %s (use 'claude' or 'codex')", skillsAgent)
	}
}

func normalizeSkillsAgent(name string) string {
	normalized := strings.ToLower(strings.TrimSpace(name))
	normalized = strings.ReplaceAll(normalized, "-", "")
	normalized = strings.ReplaceAll(normalized, "_", "")
	normalized = strings.ReplaceAll(normalized, " ", "")
	return normalized
}
