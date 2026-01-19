package agent

import (
	"os"
	"path/filepath"

	"github.com/agentsdance/agentx/internal/config"
	"github.com/agentsdance/agentx/internal/plugins"
	"github.com/agentsdance/agentx/internal/skills"
)

// ClaudeAgent represents Claude Code agent
type ClaudeAgent struct {
	configPath string
}

// NewClaudeAgent creates a new Claude Code agent
func NewClaudeAgent() *ClaudeAgent {
	home, _ := os.UserHomeDir()
	return &ClaudeAgent{
		configPath: filepath.Join(home, ".claude.json"),
	}
}

func (a *ClaudeAgent) Name() string {
	return "Claude Code"
}

func (a *ClaudeAgent) ConfigPath() string {
	return a.configPath
}

func (a *ClaudeAgent) Exists() bool {
	_, err := os.Stat(a.configPath)
	return err == nil
}

func (a *ClaudeAgent) HasPlaywright() (bool, error) {
	cfg, err := config.ReadConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return config.HasPlaywrightMCP(cfg), nil
}

func (a *ClaudeAgent) InstallPlaywright() error {
	cfg, err := config.ReadConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			cfg = make(map[string]interface{})
		} else {
			return err
		}
	}
	config.AddPlaywrightMCP(cfg)
	return config.WriteConfig(a.configPath, cfg)
}

func (a *ClaudeAgent) RemovePlaywright() error {
	cfg, err := config.ReadConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	config.RemovePlaywrightMCP(cfg)
	return config.WriteConfig(a.configPath, cfg)
}

func (a *ClaudeAgent) HasContext7() (bool, error) {
	cfg, err := config.ReadConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return config.HasContext7MCP(cfg), nil
}

func (a *ClaudeAgent) InstallContext7() error {
	cfg, err := config.ReadConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			cfg = make(map[string]interface{})
		} else {
			return err
		}
	}
	config.AddContext7MCP(cfg)
	return config.WriteConfig(a.configPath, cfg)
}

func (a *ClaudeAgent) RemoveContext7() error {
	cfg, err := config.ReadConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	config.RemoveContext7MCP(cfg)
	return config.WriteConfig(a.configPath, cfg)
}

func (a *ClaudeAgent) HasRemixIcon() (bool, error) {
	cfg, err := config.ReadConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return config.HasRemixIconMCP(cfg), nil
}

func (a *ClaudeAgent) InstallRemixIcon() error {
	cfg, err := config.ReadConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			cfg = make(map[string]interface{})
		} else {
			return err
		}
	}
	config.AddRemixIconMCP(cfg)
	return config.WriteConfig(a.configPath, cfg)
}

func (a *ClaudeAgent) RemoveRemixIcon() error {
	cfg, err := config.ReadConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	config.RemoveRemixIconMCP(cfg)
	return config.WriteConfig(a.configPath, cfg)
}

func (a *ClaudeAgent) SupportsSkills() bool {
	return true
}

func (a *ClaudeAgent) HasSkill(skillName string) (bool, error) {
	mgr := skills.NewSkillManager()
	skill, err := mgr.Get(skillName)
	if err != nil {
		return false, nil // Not found is not an error
	}
	return skill != nil, nil
}

func (a *ClaudeAgent) InstallSkill(skillName, source string) error {
	mgr := skills.NewSkillManager()
	_, err := mgr.Install(source, skills.ScopePersonal)
	return err
}

func (a *ClaudeAgent) RemoveSkill(skillName string) error {
	mgr := skills.NewSkillManager()
	return mgr.Remove(skillName, skills.ScopePersonal)
}

func (a *ClaudeAgent) SupportsPlugins() bool {
	return true
}

func (a *ClaudeAgent) HasPlugin(pluginName string) (bool, error) {
	mgr := plugins.NewPluginManager()
	plugin, err := mgr.Get(pluginName)
	if err != nil {
		return false, nil // Not found is not an error
	}
	return plugin != nil, nil
}

func (a *ClaudeAgent) InstallPlugin(pluginName, source string) error {
	mgr := plugins.NewPluginManager()
	_, err := mgr.Install(source)
	return err
}

func (a *ClaudeAgent) RemovePlugin(pluginName string) error {
	mgr := plugins.NewPluginManager()
	return mgr.Remove(pluginName)
}
