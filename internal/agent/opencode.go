package agent

import (
	"os"
	"path/filepath"

	"github.com/agentsdance/agentx/internal/config"
)

// OpenCodeAgent represents OpenCode agent
type OpenCodeAgent struct {
	configPath string
}

// NewOpenCodeAgent creates a new OpenCode agent
func NewOpenCodeAgent() *OpenCodeAgent {
	home, _ := os.UserHomeDir()
	return &OpenCodeAgent{
		configPath: filepath.Join(home, ".opencode", "config.json"),
	}
}

func (a *OpenCodeAgent) Name() string {
	return "opencode"
}

func (a *OpenCodeAgent) ConfigPath() string {
	return a.configPath
}

func (a *OpenCodeAgent) Exists() bool {
	_, err := os.Stat(a.configPath)
	return err == nil
}

func (a *OpenCodeAgent) HasPlaywright() (bool, error) {
	cfg, err := config.ReadConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return config.HasPlaywrightMCP(cfg), nil
}

func (a *OpenCodeAgent) InstallPlaywright() error {
	cfg, err := config.ReadConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			cfg = make(map[string]interface{})
			// Ensure directory exists
			if err := os.MkdirAll(filepath.Dir(a.configPath), 0755); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	config.AddPlaywrightMCP(cfg)
	return config.WriteConfig(a.configPath, cfg)
}

func (a *OpenCodeAgent) RemovePlaywright() error {
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

func (a *OpenCodeAgent) HasContext7() (bool, error) {
	cfg, err := config.ReadConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return config.HasContext7MCP(cfg), nil
}

func (a *OpenCodeAgent) InstallContext7() error {
	cfg, err := config.ReadConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			cfg = make(map[string]interface{})
			if err := os.MkdirAll(filepath.Dir(a.configPath), 0755); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	config.AddContext7MCP(cfg)
	return config.WriteConfig(a.configPath, cfg)
}

func (a *OpenCodeAgent) RemoveContext7() error {
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

func (a *OpenCodeAgent) HasRemixIcon() (bool, error) {
	cfg, err := config.ReadConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return config.HasRemixIconMCP(cfg), nil
}

func (a *OpenCodeAgent) InstallRemixIcon() error {
	cfg, err := config.ReadConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			cfg = make(map[string]interface{})
			if err := os.MkdirAll(filepath.Dir(a.configPath), 0755); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	config.AddRemixIconMCP(cfg)
	return config.WriteConfig(a.configPath, cfg)
}

func (a *OpenCodeAgent) RemoveRemixIcon() error {
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

func (a *OpenCodeAgent) SupportsSkills() bool {
	return false
}

func (a *OpenCodeAgent) HasSkill(skillName string) (bool, error) {
	return false, nil
}

func (a *OpenCodeAgent) InstallSkill(skillName, source string) error {
	return nil
}

func (a *OpenCodeAgent) RemoveSkill(skillName string) error {
	return nil
}

func (a *OpenCodeAgent) SupportsPlugins() bool {
	return false
}

func (a *OpenCodeAgent) HasPlugin(pluginName string) (bool, error) {
	return false, nil
}

func (a *OpenCodeAgent) InstallPlugin(pluginName, source string) error {
	return nil
}

func (a *OpenCodeAgent) RemovePlugin(pluginName string) error {
	return nil
}
