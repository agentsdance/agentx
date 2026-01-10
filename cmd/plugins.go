package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/agentsdance/agentx/internal/plugins"
)

var pluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "Manage Claude Code plugins",
	Long: `Manage Claude Code plugins.

Plugins are stored in: ~/.agentx/plugins/

Plugins are bundles that can contain multiple extension types:
  - Commands (slash commands)
  - Agents (specialized agents)
  - Skills (agent skills)
  - Hooks (event handlers)
  - MCP Servers (external tool integrations)`,
}

var pluginsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed plugins",
	Long:  `List all installed plugins.`,
	Run: func(cmd *cobra.Command, args []string) {
		mgr := plugins.NewPluginManager()

		pluginList, err := mgr.List()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if len(pluginList) == 0 {
			fmt.Println("No plugins installed")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tVERSION\tCOMPONENTS\tDESCRIPTION")
		fmt.Fprintln(w, "----\t-------\t----------\t-----------")
		for _, p := range pluginList {
			desc := p.Description
			if len(desc) > 40 {
				desc = desc[:37] + "..."
			}
			components := plugins.ComponentsSummary(p.Components)
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", p.Name, p.Version, components, desc)
		}
		w.Flush()
	},
}

var pluginsInstallCmd = &cobra.Command{
	Use:   "install <source>",
	Short: "Install a plugin from URL or path",
	Long: `Install a plugin from a local path or Git repository.

Examples:
  agentx plugins install ./my-plugin/
  agentx plugins install https://github.com/user/plugin-repo
  agentx plugins install https://github.com/user/repo#plugin-name
  agentx plugins install https://github.com/user/repo/tree/main/plugins/my-plugin

The source can be:
  - A local directory containing .claude-plugin/plugin.json
  - A Git repository URL
  - A Git repository URL with #plugin-name fragment
  - A GitHub tree URL pointing to a plugin directory`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		source := args[0]

		mgr := plugins.NewPluginManager()
		plugin, err := mgr.Install(source)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Installed plugin '%s' (v%s)\n", plugin.Name, plugin.Version)
		fmt.Printf("  Path: %s\n", plugin.Path)
		if plugin.Description != "" {
			fmt.Printf("  Description: %s\n", plugin.Description)
		}

		// Show components
		if len(plugin.Components.Commands) > 0 {
			fmt.Printf("  Commands: %s\n", strings.Join(plugin.Components.Commands, ", "))
		}
		if len(plugin.Components.Agents) > 0 {
			fmt.Printf("  Agents: %s\n", strings.Join(plugin.Components.Agents, ", "))
		}
		if len(plugin.Components.Skills) > 0 {
			fmt.Printf("  Skills: %s\n", strings.Join(plugin.Components.Skills, ", "))
		}
		if len(plugin.Components.Hooks) > 0 {
			fmt.Printf("  Hooks: %s\n", strings.Join(plugin.Components.Hooks, ", "))
		}
		if len(plugin.Components.MCPServers) > 0 {
			fmt.Printf("  MCP Servers: %s\n", strings.Join(plugin.Components.MCPServers, ", "))
		}
	},
}

var pluginsRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a plugin",
	Long:  `Remove a plugin by name.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		mgr := plugins.NewPluginManager()
		if err := mgr.Remove(name); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Removed plugin: %s\n", name)
	},
}

var pluginsCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check plugins installation status",
	Long:  `Verify that all installed plugins are valid and properly configured.`,
	Run: func(cmd *cobra.Command, args []string) {
		mgr := plugins.NewPluginManager()
		statuses, err := mgr.Check()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if len(statuses) == 0 {
			fmt.Println("No plugins installed")
			return
		}

		fmt.Println("Plugins Health Check")
		fmt.Println("====================")
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

			fmt.Printf("%-20s %s%s%s (v%s)\n",
				s.Plugin.Name,
				statusColor, status, resetColor,
				s.Plugin.Version)

			for _, issue := range s.Issues {
				fmt.Printf("  - %s\n", issue)
			}
			if s.Error != nil {
				fmt.Printf("  - Error: %v\n", s.Error)
			}
		}

		if !hasIssues {
			fmt.Println()
			fmt.Println("All plugins are healthy!")
		}
	},
}

var pluginsInfoCmd = &cobra.Command{
	Use:   "info <name>",
	Short: "Show detailed plugin information",
	Long:  `Show detailed information about an installed plugin.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		mgr := plugins.NewPluginManager()
		plugin, err := mgr.Get(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Plugin: %s\n", plugin.Name)
		fmt.Printf("Version: %s\n", plugin.Version)
		if plugin.Description != "" {
			fmt.Printf("Description: %s\n", plugin.Description)
		}
		if plugin.Author.Name != "" {
			author := plugin.Author.Name
			if plugin.Author.Email != "" {
				author += " <" + plugin.Author.Email + ">"
			}
			fmt.Printf("Author: %s\n", author)
		}
		fmt.Printf("Path: %s\n", plugin.Path)
		if plugin.Source != "" {
			fmt.Printf("Source: %s\n", plugin.Source)
		}

		fmt.Println()
		fmt.Println("Components:")
		fmt.Println("-----------")

		if len(plugin.Components.Commands) > 0 {
			fmt.Printf("Commands (%d):\n", len(plugin.Components.Commands))
			for _, c := range plugin.Components.Commands {
				fmt.Printf("  - %s\n", c)
			}
		}

		if len(plugin.Components.Agents) > 0 {
			fmt.Printf("Agents (%d):\n", len(plugin.Components.Agents))
			for _, a := range plugin.Components.Agents {
				fmt.Printf("  - %s\n", a)
			}
		}

		if len(plugin.Components.Skills) > 0 {
			fmt.Printf("Skills (%d):\n", len(plugin.Components.Skills))
			for _, s := range plugin.Components.Skills {
				fmt.Printf("  - %s\n", s)
			}
		}

		if len(plugin.Components.Hooks) > 0 {
			fmt.Printf("Hooks (%d):\n", len(plugin.Components.Hooks))
			for _, h := range plugin.Components.Hooks {
				fmt.Printf("  - %s\n", h)
			}
		}

		if len(plugin.Components.MCPServers) > 0 {
			fmt.Printf("MCP Servers (%d):\n", len(plugin.Components.MCPServers))
			for _, m := range plugin.Components.MCPServers {
				fmt.Printf("  - %s\n", m)
			}
		}

		if len(plugin.Components.Commands) == 0 &&
			len(plugin.Components.Agents) == 0 &&
			len(plugin.Components.Skills) == 0 &&
			len(plugin.Components.Hooks) == 0 &&
			len(plugin.Components.MCPServers) == 0 {
			fmt.Println("  (no components found)")
		}
	},
}

var pluginsAvailableCmd = &cobra.Command{
	Use:   "available",
	Short: "List available plugins from registry",
	Long:  `List all plugins available in the registry.`,
	Run: func(cmd *cobra.Command, args []string) {
		registryPlugins, err := plugins.FetchRegistryWithFallback()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching registry: %v\n", err)
			os.Exit(1)
		}

		if len(registryPlugins) == 0 {
			fmt.Println("No plugins available in registry")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tVERSION\tAUTHOR\tDESCRIPTION")
		fmt.Fprintln(w, "----\t-------\t------\t-----------")
		for _, p := range registryPlugins {
			desc := p.Description
			if len(desc) > 40 {
				desc = desc[:37] + "..."
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", p.Name, p.Version, p.Author, desc)
		}
		w.Flush()
	},
}

func init() {
	pluginsCmd.AddCommand(pluginsListCmd)
	pluginsCmd.AddCommand(pluginsAvailableCmd)
	pluginsCmd.AddCommand(pluginsInstallCmd)
	pluginsCmd.AddCommand(pluginsRemoveCmd)
	pluginsCmd.AddCommand(pluginsCheckCmd)
	pluginsCmd.AddCommand(pluginsInfoCmd)
}
