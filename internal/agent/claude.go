package agent

import (
	"os"
	"path/filepath"

	"github.com/agentsdance/agentx/internal/config"
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
