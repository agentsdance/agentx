package agent

import (
	"encoding/json"
	"fmt"
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
	return a.HasMCP("playwright")
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
	return a.HasMCP("context7")
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
	return a.HasMCP("remix-icon")
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

func (a *GeminiAgent) HasMCP(name string) (bool, error) {
	cfg, err := config.ReadConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			cfg = map[string]interface{}{}
		} else {
			return false, err
		}
	}
	if config.HasMCP(cfg, name) {
		return true, nil
	}
	extensions, err := a.listExtensionMCPs()
	if err != nil {
		return false, err
	}
	_, ok := extensions[name]
	return ok, nil
}

func (a *GeminiAgent) InstallMCP(name string, mcpConfig map[string]interface{}) error {
	if mcpConfig == nil {
		return fmt.Errorf("missing mcp config")
	}
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
	config.AddMCP(cfg, name, cloneMCPConfig(mcpConfig))
	return config.WriteConfig(a.configPath, cfg)
}

func (a *GeminiAgent) RemoveMCP(name string) error {
	cfg, err := config.ReadConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			cfg = map[string]interface{}{}
		} else {
			return err
		}
	}
	if !config.HasMCP(cfg, name) {
		extensions, extErr := a.listExtensionMCPs()
		if extErr == nil {
			if _, ok := extensions[name]; ok {
				return fmt.Errorf("mcp %s is managed by a Gemini extension", name)
			}
		}
		return nil
	}
	config.RemoveMCP(cfg, name)
	return config.WriteConfig(a.configPath, cfg)
}

func (a *GeminiAgent) ListMCPs() (map[string]map[string]interface{}, error) {
	cfg, err := config.ReadConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			cfg = map[string]interface{}{}
		} else {
			return nil, err
		}
	}
	result := config.GetMCPServers(cfg)
	extensions, err := a.listExtensionMCPs()
	if err != nil {
		return result, nil
	}
	for name, cfg := range extensions {
		if _, exists := result[name]; exists {
			continue
		}
		result[name] = cfg
	}
	return result, nil
}

func (a *GeminiAgent) listExtensionMCPs() (map[string]map[string]interface{}, error) {
	extensionsDir := filepath.Join(filepath.Dir(a.configPath), "extensions")
	entries, err := os.ReadDir(extensionsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]map[string]interface{}{}, nil
		}
		return nil, err
	}

	enabled, _ := readGeminiExtensionEnablement(filepath.Join(extensionsDir, "extension-enablement.json"))
	result := make(map[string]map[string]interface{})
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if enabled != nil {
			if _, ok := enabled[name]; !ok {
				continue
			}
		}
		configPath := filepath.Join(extensionsDir, name, "gemini-extension.json")
		extCfg, err := readGeminiExtensionConfig(configPath)
		if err != nil {
			continue
		}
		for serverName, serverCfg := range extCfg {
			if serverCfg == nil {
				continue
			}
			result[serverName] = serverCfg
		}
	}
	return result, nil
}

func readGeminiExtensionEnablement(path string) (map[string]struct{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	result := make(map[string]struct{}, len(raw))
	for key := range raw {
		result[key] = struct{}{}
	}
	return result, nil
}

func readGeminiExtensionConfig(path string) (map[string]map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var payload struct {
		MCPServers map[string]map[string]interface{} `json:"mcpServers"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	if payload.MCPServers == nil {
		return map[string]map[string]interface{}{}, nil
	}
	return payload.MCPServers, nil
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
