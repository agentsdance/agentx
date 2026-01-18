package cmd

import (
	"fmt"
	"os"

	"github.com/agentsdance/agentx/internal/version"
	"github.com/agentsdance/agentx/ui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "agentx",
	Aliases: []string{"agents", "ax"},
	Short:   "Unified MCP Servers & Agent Skills Manager for AI coding agents",
	Long: `agentx is a CLI tool for managing MCP servers and skills across AI coding agents
(Claude Code, Codex, Cursor, Gemini cli, opencode).

Run without arguments to launch the TUI interface.

Aliases: agents, ax`,
	Version: version.GetFullVersion(),
	Run: func(cmd *cobra.Command, args []string) {
		// Launch TUI when no subcommand is provided
		if err := ui.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate("{{.Version}}\n")
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(skillsCmd)
	rootCmd.AddCommand(pluginsCmd)
}
