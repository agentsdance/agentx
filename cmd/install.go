package cmd

import (
	"fmt"

	"github.com/agentsdance/agentx/internal/agent"
	"github.com/spf13/cobra"
)

var agentFlag string

var installCmd = &cobra.Command{
	Use:   "install [mcp-server]",
	Short: "Install an MCP server to agents",
	Long:  `Install an MCP server to all agents or a specific agent.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverName := args[0]
		if serverName != "playwright" && serverName != "context7" && serverName != "remix-icon" {
			fmt.Printf("Unknown MCP server: %s (supported: playwright, context7, remix-icon)\n", serverName)
			return
		}

		var agents []agent.Agent
		if agentFlag != "" {
			a := agent.GetAgentByName(agentFlag)
			if a == nil {
				fmt.Printf("Unknown agent: %s\n", agentFlag)
				return
			}
			agents = []agent.Agent{a}
		} else {
			agents = agent.GetAllAgents()
		}

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
			if has {
				fmt.Printf("%-12s already installed\n", a.Name())
				continue
			}

			switch serverName {
			case "playwright":
				err = a.InstallPlaywright()
			case "context7":
				err = a.InstallContext7()
			case "remix-icon":
				err = a.InstallRemixIcon()
			}

			if err != nil {
				fmt.Printf("%-12s failed: %v\n", a.Name(), err)
			} else {
				fmt.Printf("%-12s installed\n", a.Name())
			}
		}
	},
}

func init() {
	installCmd.Flags().StringVarP(&agentFlag, "agent", "a", "", "Target agent (claude, codex, cursor, gemini, opencode)")
}
