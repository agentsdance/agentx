package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/agentsdance/agentx/internal/agent"
)

var checkCmd = &cobra.Command{
	Use:   "check [agent]",
	Short: "Check Playwright MCP installation status",
	Long:  `Check if Playwright MCP is configured for all agents or a specific agent.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var agents []agent.Agent
		if len(args) > 0 {
			a := agent.GetAgentByName(args[0])
			if a == nil {
				fmt.Printf("Unknown agent: %s\n", args[0])
				return
			}
			agents = []agent.Agent{a}
		} else {
			agents = agent.GetAllAgents()
		}

		fmt.Println("Playwright MCP Status")
		fmt.Println("---------------------")
		for _, a := range agents {
			status := "not configured"
			has, err := a.HasPlaywright()
			if err != nil {
				status = fmt.Sprintf("error: %v", err)
			} else if has {
				status = "installed"
			} else if !a.Exists() {
				status = "config not found"
			}
			fmt.Printf("%-12s %s\n", a.Name(), status)
		}
	},
}
