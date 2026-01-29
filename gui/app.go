package main

import (
	"context"

	"github.com/agentsdance/agentx/internal/agent"
	"github.com/agentsdance/agentx/internal/plugins"
	"github.com/agentsdance/agentx/internal/skills"
	"github.com/agentsdance/agentx/internal/version"
)

// App struct holds the application state
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// AgentInfo represents agent information for the frontend
type AgentInfo struct {
	Name       string `json:"name"`
	ConfigPath string `json:"configPath"`
	Exists     bool   `json:"exists"`
}

// MCPInfo represents MCP server information
type MCPInfo struct {
	Name   string                 `json:"name"`
	Config map[string]interface{} `json:"config"`
}

// SkillInfo represents skill information for the frontend
type SkillInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Source      string `json:"source"`
}

// PluginInfo represents plugin information for the frontend
type PluginInfo struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Source      string   `json:"source"`
	Version     string   `json:"version"`
	Author      string   `json:"author"`
	Components  []string `json:"components"`
}

// GetAgents returns all available agents with their status
func (a *App) GetAgents() []AgentInfo {
	agents := agent.GetAllAgents()
	result := make([]AgentInfo, len(agents))
	for i, ag := range agents {
		result[i] = AgentInfo{
			Name:       ag.Name(),
			ConfigPath: ag.ConfigPath(),
			Exists:     ag.Exists(),
		}
	}
	return result
}

// GetMCPs returns all MCP servers for an agent
func (a *App) GetMCPs(agentName string) ([]MCPInfo, error) {
	ag := agent.GetAgentByName(agentName)
	if ag == nil {
		return nil, nil
	}

	mcps, err := ag.ListMCPs()
	if err != nil {
		return nil, err
	}

	result := make([]MCPInfo, 0, len(mcps))
	for name, config := range mcps {
		result = append(result, MCPInfo{
			Name:   name,
			Config: config,
		})
	}
	return result, nil
}

// InstallMCP installs an MCP server for an agent
func (a *App) InstallMCP(agentName, mcpName string) error {
	ag := agent.GetAgentByName(agentName)
	if ag == nil {
		return nil
	}

	switch mcpName {
	case "playwright":
		return ag.InstallPlaywright()
	case "context7":
		return ag.InstallContext7()
	case "remix-icon":
		return ag.InstallRemixIcon()
	}
	return nil
}

// RemoveMCP removes an MCP server from an agent
func (a *App) RemoveMCP(agentName, mcpName string) error {
	ag := agent.GetAgentByName(agentName)
	if ag == nil {
		return nil
	}

	switch mcpName {
	case "playwright":
		return ag.RemovePlaywright()
	case "context7":
		return ag.RemoveContext7()
	case "remix-icon":
		return ag.RemoveRemixIcon()
	default:
		return ag.RemoveMCP(mcpName)
	}
}

// HasMCP checks if an MCP server is installed for an agent
func (a *App) HasMCP(agentName, mcpName string) (bool, error) {
	ag := agent.GetAgentByName(agentName)
	if ag == nil {
		return false, nil
	}

	switch mcpName {
	case "playwright":
		return ag.HasPlaywright()
	case "context7":
		return ag.HasContext7()
	case "remix-icon":
		return ag.HasRemixIcon()
	default:
		return ag.HasMCP(mcpName)
	}
}

// GetSkillsRegistry returns available skills from the registry
func (a *App) GetSkillsRegistry() ([]SkillInfo, error) {
	registrySkills, err := skills.FetchSkillsRegistryWithFallback()
	if err != nil {
		return nil, err
	}

	result := make([]SkillInfo, len(registrySkills))
	for i, s := range registrySkills {
		result[i] = SkillInfo{
			Name:        s.Name,
			Description: s.Description,
			Author:      s.Author,
			Source:      s.Source,
		}
	}
	return result, nil
}

// InstallSkill installs a skill for an agent
func (a *App) InstallSkill(agentName, skillSource string) error {
	ag := agent.GetAgentByName(agentName)
	if ag == nil {
		return nil
	}

	if !ag.SupportsSkills() {
		return nil
	}

	return ag.InstallSkill("", skillSource)
}

// RemoveSkill removes a skill from an agent
func (a *App) RemoveSkill(agentName, skillName string) error {
	ag := agent.GetAgentByName(agentName)
	if ag == nil {
		return nil
	}

	return ag.RemoveSkill(skillName)
}

// HasSkill checks if a skill is installed for an agent
func (a *App) HasSkill(agentName, skillName string) (bool, error) {
	ag := agent.GetAgentByName(agentName)
	if ag == nil {
		return false, nil
	}

	return ag.HasSkill(skillName)
}

// GetPluginsRegistry returns available plugins from the registry
func (a *App) GetPluginsRegistry() ([]PluginInfo, error) {
	registryPlugins, err := plugins.FetchRegistryWithFallback()
	if err != nil {
		return nil, err
	}

	result := make([]PluginInfo, len(registryPlugins))
	for i, p := range registryPlugins {
		result[i] = PluginInfo{
			Name:        p.Name,
			Description: p.Description,
			Source:      p.Source,
			Version:     p.Version,
			Author:      p.Author,
			Components:  p.Components,
		}
	}
	return result, nil
}

// InstallPlugin installs a plugin for an agent
func (a *App) InstallPlugin(agentName, pluginSource string) error {
	ag := agent.GetAgentByName(agentName)
	if ag == nil {
		return nil
	}

	if !ag.SupportsPlugins() {
		return nil
	}

	return ag.InstallPlugin("", pluginSource)
}

// RemovePlugin removes a plugin from an agent
func (a *App) RemovePlugin(agentName, pluginName string) error {
	ag := agent.GetAgentByName(agentName)
	if ag == nil {
		return nil
	}

	return ag.RemovePlugin(pluginName)
}

// HasPlugin checks if a plugin is installed for an agent
func (a *App) HasPlugin(agentName, pluginName string) (bool, error) {
	ag := agent.GetAgentByName(agentName)
	if ag == nil {
		return false, nil
	}

	return ag.HasPlugin(pluginName)
}

// SupportsSkills checks if an agent supports skills
func (a *App) SupportsSkills(agentName string) bool {
	ag := agent.GetAgentByName(agentName)
	if ag == nil {
		return false
	}
	return ag.SupportsSkills()
}

// SupportsPlugins checks if an agent supports plugins
func (a *App) SupportsPlugins(agentName string) bool {
	ag := agent.GetAgentByName(agentName)
	if ag == nil {
		return false
	}
	return ag.SupportsPlugins()
}

// MCPStatus represents MCP installation status across all agents
type MCPStatus struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Agents      map[string]string `json:"agents"` // agentName -> "installed" | "not_installed" | "n/a" | "error"
}

// SkillStatus represents skill installation status across all agents
type SkillStatus struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Source      string            `json:"source"`
	Agents      map[string]string `json:"agents"`
}

// PluginStatus represents plugin installation status across all agents
type PluginStatus struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Source      string            `json:"source"`
	Agents      map[string]string `json:"agents"`
}

// GetMCPMatrix returns MCP installation status matrix (like TUI)
func (a *App) GetMCPMatrix() []MCPStatus {
	agents := agent.GetAllAgents()
	mcpServers := []struct {
		name string
		desc string
	}{
		{"playwright", "Browser automation"},
		{"context7", "Library documentation"},
		{"remix-icon", "Icon library"},
	}

	result := make([]MCPStatus, len(mcpServers))
	for i, mcp := range mcpServers {
		status := MCPStatus{
			Name:        mcp.name,
			Description: mcp.desc,
			Agents:      make(map[string]string),
		}

		for _, ag := range agents {
			if !ag.Exists() {
				status.Agents[ag.Name()] = "n/a"
				continue
			}

			has, err := ag.HasMCP(mcp.name)
			if err != nil {
				status.Agents[ag.Name()] = "error"
			} else if has {
				status.Agents[ag.Name()] = "installed"
			} else {
				status.Agents[ag.Name()] = "not_installed"
			}
		}
		result[i] = status
	}

	return result
}

// GetSkillsMatrix returns skills installation status matrix
func (a *App) GetSkillsMatrix() []SkillStatus {
	agents := agent.GetAllAgents()
	registrySkills, err := skills.FetchSkillsRegistryWithFallback()
	if err != nil {
		return []SkillStatus{}
	}

	result := make([]SkillStatus, len(registrySkills))
	for i, skill := range registrySkills {
		source := skill.Source
		if source == "local" {
			source = "https://github.com/agentsdance/agentskills/tree/master/skills/" + skill.Name
		}

		status := SkillStatus{
			Name:        skill.Name,
			Description: skill.Description,
			Source:      source,
			Agents:      make(map[string]string),
		}

		for _, ag := range agents {
			if !ag.Exists() {
				status.Agents[ag.Name()] = "n/a"
				continue
			}
			if !ag.SupportsSkills() {
				status.Agents[ag.Name()] = "n/a"
				continue
			}

			has, err := ag.HasSkill(skill.Name)
			if err != nil {
				status.Agents[ag.Name()] = "error"
			} else if has {
				status.Agents[ag.Name()] = "installed"
			} else {
				status.Agents[ag.Name()] = "not_installed"
			}
		}
		result[i] = status
	}

	return result
}

// GetPluginsMatrix returns plugins installation status matrix
func (a *App) GetPluginsMatrix() []PluginStatus {
	agents := agent.GetAllAgents()
	registryPlugins, err := plugins.FetchRegistryWithFallback()
	if err != nil {
		return []PluginStatus{}
	}

	result := make([]PluginStatus, len(registryPlugins))
	for i, plugin := range registryPlugins {
		status := PluginStatus{
			Name:        plugin.Name,
			Description: plugin.Description,
			Source:      plugin.Source,
			Agents:      make(map[string]string),
		}

		for _, ag := range agents {
			if !ag.Exists() {
				status.Agents[ag.Name()] = "n/a"
				continue
			}
			if !ag.SupportsPlugins() {
				status.Agents[ag.Name()] = "n/a"
				continue
			}

			has, err := ag.HasPlugin(plugin.Name)
			if err != nil {
				status.Agents[ag.Name()] = "error"
			} else if has {
				status.Agents[ag.Name()] = "installed"
			} else {
				status.Agents[ag.Name()] = "not_installed"
			}
		}
		result[i] = status
	}

	return result
}

// InstallMCPForAgent installs an MCP for an agent (for matrix view)
func (a *App) InstallMCPForAgent(agentName, mcpName string) error {
	return a.InstallMCP(agentName, mcpName)
}

// RemoveMCPFromAgent removes an MCP from an agent (for matrix view)
func (a *App) RemoveMCPFromAgent(agentName, mcpName string) error {
	return a.RemoveMCP(agentName, mcpName)
}

// InstallSkillForAgent installs a skill for an agent
func (a *App) InstallSkillForAgent(agentName, skillName, source string) error {
	ag := agent.GetAgentByName(agentName)
	if ag == nil || !ag.SupportsSkills() {
		return nil
	}
	return ag.InstallSkill(skillName, source)
}

// RemoveSkillFromAgent removes a skill from an agent
func (a *App) RemoveSkillFromAgent(agentName, skillName string) error {
	ag := agent.GetAgentByName(agentName)
	if ag == nil {
		return nil
	}
	return ag.RemoveSkill(skillName)
}

// InstallPluginForAgent installs a plugin for an agent
func (a *App) InstallPluginForAgent(agentName, pluginName, source string) error {
	ag := agent.GetAgentByName(agentName)
	if ag == nil || !ag.SupportsPlugins() {
		return nil
	}
	return ag.InstallPlugin(pluginName, source)
}

// RemovePluginFromAgent removes a plugin from an agent
func (a *App) RemovePluginFromAgent(agentName, pluginName string) error {
	ag := agent.GetAgentByName(agentName)
	if ag == nil {
		return nil
	}
	return ag.RemovePlugin(pluginName)
}

// GetVersion returns the current version of AgentX
func (a *App) GetVersion() string {
	return version.Version
}

// InstallMCPForAll installs an MCP server to all available agents
func (a *App) InstallMCPForAll(mcpName string) error {
	agents := agent.GetAllAgents()
	var lastErr error
	for _, ag := range agents {
		if !ag.Exists() {
			continue
		}
		var err error
		switch mcpName {
		case "playwright":
			err = ag.InstallPlaywright()
		case "context7":
			err = ag.InstallContext7()
		case "remix-icon":
			err = ag.InstallRemixIcon()
		}
		if err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// InstallSkillForAll installs a skill to all available agents that support skills
func (a *App) InstallSkillForAll(skillName, source string) error {
	agents := agent.GetAllAgents()
	var lastErr error
	for _, ag := range agents {
		if !ag.Exists() || !ag.SupportsSkills() {
			continue
		}
		if err := ag.InstallSkill(skillName, source); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// InstallPluginForAll installs a plugin to all available agents that support plugins
func (a *App) InstallPluginForAll(pluginName, source string) error {
	agents := agent.GetAllAgents()
	var lastErr error
	for _, ag := range agents {
		if !ag.Exists() || !ag.SupportsPlugins() {
			continue
		}
		if err := ag.InstallPlugin(pluginName, source); err != nil {
			lastErr = err
		}
	}
	return lastErr
}
