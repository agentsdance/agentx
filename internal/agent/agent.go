package agent

// Agent represents an AI coding agent that supports MCP servers
type Agent interface {
	// Name returns the display name of the agent
	Name() string
	// ConfigPath returns the path to the agent's config file
	ConfigPath() string
	// Exists returns true if the agent's config file exists
	Exists() bool
	// HasPlaywright checks if Playwright MCP is configured
	HasPlaywright() (bool, error)
	// InstallPlaywright adds Playwright MCP to the config
	InstallPlaywright() error
	// RemovePlaywright removes Playwright MCP from the config
	RemovePlaywright() error
	// HasContext7 checks if Context7 MCP is configured
	HasContext7() (bool, error)
	// InstallContext7 adds Context7 MCP to the config
	InstallContext7() error
	// RemoveContext7 removes Context7 MCP from the config
	RemoveContext7() error
	// SupportsSkills returns true if the agent supports skills
	SupportsSkills() bool
	// HasSkill checks if a skill is installed
	HasSkill(skillName string) (bool, error)
	// InstallSkill installs a skill from a source URL
	InstallSkill(skillName, source string) error
	// RemoveSkill removes a skill by name
	RemoveSkill(skillName string) error
	// SupportsPlugins returns true if the agent supports plugins
	SupportsPlugins() bool
	// HasPlugin checks if a plugin is installed
	HasPlugin(pluginName string) (bool, error)
	// InstallPlugin installs a plugin from a source URL
	InstallPlugin(pluginName, source string) error
	// RemovePlugin removes a plugin by name
	RemovePlugin(pluginName string) error
}

// GetAllAgents returns all supported agents
func GetAllAgents() []Agent {
	return []Agent{
		NewClaudeAgent(),
		NewCodexAgent(),
		NewCursorAgent(),
		NewGeminiAgent(),
		NewOpenCodeAgent(),
	}
}

// GetAgentByName returns an agent by name (case-insensitive)
func GetAgentByName(name string) Agent {
	for _, a := range GetAllAgents() {
		if matchAgentName(a.Name(), name) {
			return a
		}
	}
	return nil
}

func matchAgentName(agentName, input string) bool {
	// Simple case-insensitive prefix matching
	input = toLower(input)
	switch input {
	case "claude", "claudecode", "claude-code", "claude_code":
		return agentName == "Claude Code"
	case "codex", "codexcli", "codex-cli", "codex_cli":
		return agentName == "Codex"
	case "cursor":
		return agentName == "Cursor"
	case "gemini", "geminicli", "gemini-cli", "gemini_cli":
		return agentName == "Gemini cli"
	case "opencode", "open-code", "open_code":
		return agentName == "opencode"
	}
	return false
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
}
