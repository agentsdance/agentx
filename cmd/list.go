package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/agentsdance/agentx/internal/agent"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Playwright MCP status across all agents",
	Long:  `Display Playwright MCP installation status for all supported agents.`,
	Run: func(cmd *cobra.Command, args []string) {
		agents := agent.GetAllAgents()

		fmt.Println()
		fmt.Println("  Agent           Status         Config Path")
		fmt.Println("  ─────────────────────────────────────────────────────────────")

		for _, a := range agents {
			status := "○ not configured"
			has, err := a.HasPlaywright()
			if err != nil {
				status = fmt.Sprintf("✗ error: %v", err)
			} else if has {
				status = "✓ installed"
			} else if !a.Exists() {
				status = "○ config not found"
			}
			fmt.Printf("  %-14s  %-17s  %s\n", a.Name(), status, a.ConfigPath())
		}
		fmt.Println()
	},
}
