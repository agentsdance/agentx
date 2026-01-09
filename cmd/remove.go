package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/agentsdance/agentx/internal/agent"
)

var removeCmd = &cobra.Command{
	Use:   "remove [mcp-server]",
	Short: "Remove an MCP server from agents",
	Long:  `Remove an MCP server from all agents.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverName := args[0]
		if serverName != "playwright" {
			fmt.Printf("Unknown MCP server: %s (only 'playwright' is supported in v1)\n", serverName)
			return
		}

		agents := agent.GetAllAgents()

		for _, a := range agents {
			has, err := a.HasPlaywright()
			if err != nil {
				fmt.Printf("%-12s error: %v\n", a.Name(), err)
				continue
			}
			if !has {
				fmt.Printf("%-12s not installed\n", a.Name())
				continue
			}

			if err := a.RemovePlaywright(); err != nil {
				fmt.Printf("%-12s failed: %v\n", a.Name(), err)
			} else {
				fmt.Printf("%-12s removed\n", a.Name())
			}
		}
	},
}
