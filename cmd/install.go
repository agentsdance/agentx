package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/agentsdance/agentx/internal/agent"
)

var agentFlag string

var installCmd = &cobra.Command{
	Use:   "install [mcp-server]",
	Short: "Install an MCP server to agents",
	Long:  `Install an MCP server to all agents or a specific agent.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverName := args[0]
		if serverName != "playwright" {
			fmt.Printf("Unknown MCP server: %s (only 'playwright' is supported in v1)\n", serverName)
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
			has, err := a.HasPlaywright()
			if err != nil {
				fmt.Printf("%-12s error: %v\n", a.Name(), err)
				continue
			}
			if has {
				fmt.Printf("%-12s already installed\n", a.Name())
				continue
			}

			if err := a.InstallPlaywright(); err != nil {
				fmt.Printf("%-12s failed: %v\n", a.Name(), err)
			} else {
				fmt.Printf("%-12s installed\n", a.Name())
			}
		}
	},
}

func init() {
	installCmd.Flags().StringVarP(&agentFlag, "agent", "a", "", "Target agent (claude, gemini, opencode)")
}
