package cmd

import (
	"fmt"

	"github.com/agentsdance/agentx/internal/agent"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [mcp-server]",
	Short: "Remove an MCP server from agents",
	Long:  `Remove an MCP server from all agents.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverName := args[0]
		if serverName != "playwright" && serverName != "context7" && serverName != "remix-icon" {
			fmt.Printf("Unknown MCP server: %s (supported: playwright, context7, remix-icon)\n", serverName)
			return
		}

		agents := agent.GetAllAgents()

		for _, a := range agents {
			var has bool
			var err error
			switch serverName {
			case "playwright":
				has, err = a.HasPlaywright()
			case "context7":
				has, err = a.HasContext7()
			case "remix-icon":
				has, err = a.HasRemixIcon()
			}
			if err != nil {
				fmt.Printf("%-12s error: %v\n", a.Name(), err)
				continue
			}
			if !has {
				fmt.Printf("%-12s not installed\n", a.Name())
				continue
			}

			switch serverName {
			case "playwright":
				err = a.RemovePlaywright()
			case "context7":
				err = a.RemoveContext7()
			case "remix-icon":
				err = a.RemoveRemixIcon()
			}

			if err != nil {
				fmt.Printf("%-12s failed: %v\n", a.Name(), err)
			} else {
				fmt.Printf("%-12s removed\n", a.Name())
			}
		}
	},
}
