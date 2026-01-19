package agent

import (
	"os"
	"path/filepath"

	"github.com/agentsdance/agentx/internal/config"
)

// GeminiAgent represents Gemini CLI agent
type GeminiAgent struct {
	configPath string
}

// NewGeminiAgent creates a new Gemini CLI agent
func NewGeminiAgent() *GeminiAgent {
	home, _ := os.UserHomeDir()
	return &GeminiAgent{
		configPath: filepath.Join(home, ".gemini", "settings.json"),
	}
}

func (a *GeminiAgent) Name() string {
	return "Gemini cli"
}

func (a *GeminiAgent) ConfigPath() string {
	return a.configPath
}

func (a *GeminiAgent) Exists() bool {
	_, err := os.Stat(a.configPath)
	return err == nil
}

func (a *GeminiAgent) HasPlaywright() (bool, error) {
	cfg, err := config.ReadConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return config.HasPlaywrightMCP(cfg), nil
}

func (a *GeminiAgent) InstallPlaywright() error {
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

func (a *GeminiAgent) RemovePlaywright() error {
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

func (a *GeminiAgent) HasContext7() (bool, error) {
	cfg, err := config.ReadConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return config.HasContext7MCP(cfg), nil
}

func (a *GeminiAgent) InstallContext7() error {
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

func (a *GeminiAgent) RemoveContext7() error {
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

func (a *GeminiAgent) HasRemixIcon() (bool, error) {
	cfg, err := config.ReadConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return config.HasRemixIconMCP(cfg), nil
}

func (a *GeminiAgent) InstallRemixIcon() error {
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

func (a *GeminiAgent) RemoveRemixIcon() error {
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

func (a *GeminiAgent) SupportsSkills() bool {
	return false
}

func (a *GeminiAgent) HasSkill(skillName string) (bool, error) {
	return false, nil
}

func (a *GeminiAgent) InstallSkill(skillName, source string) error {
	return nil
}

func (a *GeminiAgent) RemoveSkill(skillName string) error {
	return nil
}

func (a *GeminiAgent) SupportsPlugins() bool {
	return false
}

func (a *GeminiAgent) HasPlugin(pluginName string) (bool, error) {
	return false, nil
}

func (a *GeminiAgent) InstallPlugin(pluginName, source string) error {
	return nil
}

func (a *GeminiAgent) RemovePlugin(pluginName string) error {
	return nil
}
