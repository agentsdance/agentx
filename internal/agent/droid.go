package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/agentsdance/agentx/internal/plugins"
	"github.com/agentsdance/agentx/internal/skills"
)

// DroidAgent represents Factory Droid agent
type DroidAgent struct {
	configPath string
}

// NewDroidAgent creates a new Factory Droid agent
func NewDroidAgent() *DroidAgent {
	home, _ := os.UserHomeDir()
	return &DroidAgent{
		configPath: filepath.Join(home, ".factory", "mcp.json"),
	}
}

func (a *DroidAgent) Name() string {
	return "Droid"
}

func (a *DroidAgent) ConfigPath() string {
	return a.configPath
}

func (a *DroidAgent) Exists() bool {
	// Check if .factory directory exists
	home, _ := os.UserHomeDir()
	factoryDir := filepath.Join(home, ".factory")
	_, err := os.Stat(factoryDir)
	return err == nil
}

func (a *DroidAgent) readConfig() (map[string]interface{}, error) {
	data, err := os.ReadFile(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]interface{}{}, nil
		}
		return nil, err
	}
	var cfg map[string]interface{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (a *DroidAgent) writeConfig(cfg map[string]interface{}) error {
	// Ensure directory exists
	dir := filepath.Dir(a.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(a.configPath, data, 0644)
}

func (a *DroidAgent) getMCPServers(cfg map[string]interface{}) map[string]interface{} {
	if servers, ok := cfg["mcpServers"].(map[string]interface{}); ok {
		return servers
	}
	return nil
}

func (a *DroidAgent) ensureMCPServers(cfg map[string]interface{}) map[string]interface{} {
	if cfg["mcpServers"] == nil {
		cfg["mcpServers"] = make(map[string]interface{})
	}
	return cfg["mcpServers"].(map[string]interface{})
}

func (a *DroidAgent) HasPlaywright() (bool, error) {
	cfg, err := a.readConfig()
	if err != nil {
		return false, err
	}
	servers := a.getMCPServers(cfg)
	if servers == nil {
		return false, nil
	}
	_, ok := servers["playwright"]
	return ok, nil
}

func (a *DroidAgent) InstallPlaywright() error {
	cfg, err := a.readConfig()
	if err != nil {
		return err
	}
	servers := a.ensureMCPServers(cfg)
	servers["playwright"] = map[string]interface{}{
		"type":    "stdio",
		"command": "npx",
		"args":    []interface{}{"-y", "@playwright/mcp@latest"},
	}
	return a.writeConfig(cfg)
}

func (a *DroidAgent) RemovePlaywright() error {
	cfg, err := a.readConfig()
	if err != nil {
		return err
	}
	servers := a.getMCPServers(cfg)
	if servers != nil {
		delete(servers, "playwright")
	}
	return a.writeConfig(cfg)
}

func (a *DroidAgent) HasContext7() (bool, error) {
	cfg, err := a.readConfig()
	if err != nil {
		return false, err
	}
	servers := a.getMCPServers(cfg)
	if servers == nil {
		return false, nil
	}
	_, ok := servers["context7"]
	return ok, nil
}

func (a *DroidAgent) InstallContext7() error {
	cfg, err := a.readConfig()
	if err != nil {
		return err
	}
	servers := a.ensureMCPServers(cfg)
	servers["context7"] = map[string]interface{}{
		"type":    "stdio",
		"command": "npx",
		"args":    []interface{}{"-y", "@context7/mcp@latest"},
	}
	return a.writeConfig(cfg)
}

func (a *DroidAgent) RemoveContext7() error {
	cfg, err := a.readConfig()
	if err != nil {
		return err
	}
	servers := a.getMCPServers(cfg)
	if servers != nil {
		delete(servers, "context7")
	}
	return a.writeConfig(cfg)
}

func (a *DroidAgent) HasRemixIcon() (bool, error) {
	cfg, err := a.readConfig()
	if err != nil {
		return false, err
	}
	servers := a.getMCPServers(cfg)
	if servers == nil {
		return false, nil
	}
	_, ok := servers["remix-icon"]
	return ok, nil
}

func (a *DroidAgent) InstallRemixIcon() error {
	cfg, err := a.readConfig()
	if err != nil {
		return err
	}
	servers := a.ensureMCPServers(cfg)
	servers["remix-icon"] = map[string]interface{}{
		"type":    "stdio",
		"command": "npx",
		"args":    []interface{}{"-y", "@nicepkg/gpt-runner", "mcp", "--remix-icon"},
	}
	return a.writeConfig(cfg)
}

func (a *DroidAgent) RemoveRemixIcon() error {
	cfg, err := a.readConfig()
	if err != nil {
		return err
	}
	servers := a.getMCPServers(cfg)
	if servers != nil {
		delete(servers, "remix-icon")
	}
	return a.writeConfig(cfg)
}

func (a *DroidAgent) HasMCP(name string) (bool, error) {
	cfg, err := a.readConfig()
	if err != nil {
		return false, err
	}
	servers := a.getMCPServers(cfg)
	if servers == nil {
		return false, nil
	}
	_, ok := servers[name]
	return ok, nil
}

func (a *DroidAgent) InstallMCP(name string, mcpConfig map[string]interface{}) error {
	if mcpConfig == nil {
		return fmt.Errorf("missing mcp config")
	}
	cfg, err := a.readConfig()
	if err != nil {
		return err
	}
	servers := a.ensureMCPServers(cfg)
	servers[name] = cloneMCPConfig(mcpConfig)
	return a.writeConfig(cfg)
}

func (a *DroidAgent) RemoveMCP(name string) error {
	cfg, err := a.readConfig()
	if err != nil {
		return err
	}
	servers := a.getMCPServers(cfg)
	if servers != nil {
		delete(servers, name)
	}
	return a.writeConfig(cfg)
}

func (a *DroidAgent) ListMCPs() (map[string]map[string]interface{}, error) {
	cfg, err := a.readConfig()
	if err != nil {
		return nil, err
	}
	servers := a.getMCPServers(cfg)
	if servers == nil {
		return map[string]map[string]interface{}{}, nil
	}
	result := make(map[string]map[string]interface{})
	for name, server := range servers {
		if serverMap, ok := server.(map[string]interface{}); ok {
			result[name] = serverMap
		}
	}
	return result, nil
}

func (a *DroidAgent) SupportsSkills() bool {
	return true
}

func (a *DroidAgent) HasSkill(skillName string) (bool, error) {
	mgr := skills.NewDroidSkillManager()
	skill, err := mgr.Get(skillName)
	if err != nil {
		return false, nil
	}
	return skill != nil, nil
}

func (a *DroidAgent) InstallSkill(skillName, source string) error {
	mgr := skills.NewDroidSkillManager()
	_, err := mgr.Install(source, skills.ScopePersonal)
	return err
}

func (a *DroidAgent) RemoveSkill(skillName string) error {
	mgr := skills.NewDroidSkillManager()
	return mgr.Remove(skillName, skills.ScopePersonal)
}

func (a *DroidAgent) SupportsPlugins() bool {
	return true
}

func (a *DroidAgent) HasPlugin(pluginName string) (bool, error) {
	mgr := plugins.NewPluginManager()
	plugin, err := mgr.Get(pluginName)
	if err != nil {
		return false, nil
	}
	return plugin != nil, nil
}

func (a *DroidAgent) InstallPlugin(pluginName, source string) error {
	mgr := plugins.NewPluginManager()
	_, err := mgr.Install(source)
	return err
}

func (a *DroidAgent) RemovePlugin(pluginName string) error {
	mgr := plugins.NewPluginManager()
	return mgr.Remove(pluginName)
}
