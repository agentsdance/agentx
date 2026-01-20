package agent

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/agentsdance/agentx/internal/config"
	"github.com/agentsdance/agentx/internal/skills"
)

const codexMCPKey = "mcp_servers"

var codexPlaywrightMCPConfig = map[string]interface{}{
	"command": "npx",
	"args":    []string{"@playwright/mcp@latest"},
}

var codexContext7MCPConfig = map[string]interface{}{
	"command": "npx",
	"args":    []string{"-y", "@upstash/context7-mcp"},
}

var codexRemixIconMCPConfig = map[string]interface{}{
	"command": "npx",
	"args":    []string{"-y", "remixicon-mcp"},
}

// CodexAgent represents Codex CLI agent
type CodexAgent struct {
	configPath string
}

// NewCodexAgent creates a new Codex agent
func NewCodexAgent() *CodexAgent {
	home, _ := os.UserHomeDir()
	codexHome := os.Getenv("CODEX_HOME")
	if codexHome == "" {
		codexHome = filepath.Join(home, ".codex")
	}

	return &CodexAgent{
		configPath: filepath.Join(codexHome, "config.toml"),
	}
}

func (a *CodexAgent) Name() string {
	return "Codex"
}

func (a *CodexAgent) ConfigPath() string {
	return a.configPath
}

func (a *CodexAgent) Exists() bool {
	_, err := os.Stat(a.configPath)
	return err == nil
}

func (a *CodexAgent) HasPlaywright() (bool, error) {
	cfg, err := config.ReadTOMLConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return hasCodexMCP(cfg, "playwright"), nil
}

func (a *CodexAgent) InstallPlaywright() error {
	cfg, err := config.ReadTOMLConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			cfg = make(map[string]interface{})
		} else {
			return err
		}
	}
	addCodexMCP(cfg, "playwright", codexPlaywrightMCPConfig)
	return config.WriteTOMLConfig(a.configPath, cfg)
}

func (a *CodexAgent) RemovePlaywright() error {
	cfg, err := config.ReadTOMLConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	removeCodexMCP(cfg, "playwright")
	return config.WriteTOMLConfig(a.configPath, cfg)
}

func (a *CodexAgent) HasContext7() (bool, error) {
	cfg, err := config.ReadTOMLConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return hasCodexMCP(cfg, "context7"), nil
}

func (a *CodexAgent) InstallContext7() error {
	cfg, err := config.ReadTOMLConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			cfg = make(map[string]interface{})
		} else {
			return err
		}
	}
	addCodexMCP(cfg, "context7", codexContext7MCPConfig)
	return config.WriteTOMLConfig(a.configPath, cfg)
}

func (a *CodexAgent) RemoveContext7() error {
	cfg, err := config.ReadTOMLConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	removeCodexMCP(cfg, "context7")
	return config.WriteTOMLConfig(a.configPath, cfg)
}

func (a *CodexAgent) HasRemixIcon() (bool, error) {
	cfg, err := config.ReadTOMLConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return hasCodexMCP(cfg, "remix-icon"), nil
}

func (a *CodexAgent) InstallRemixIcon() error {
	cfg, err := config.ReadTOMLConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			cfg = make(map[string]interface{})
		} else {
			return err
		}
	}
	addCodexMCP(cfg, "remix-icon", codexRemixIconMCPConfig)
	return config.WriteTOMLConfig(a.configPath, cfg)
}

func (a *CodexAgent) RemoveRemixIcon() error {
	cfg, err := config.ReadTOMLConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	removeCodexMCP(cfg, "remix-icon")
	return config.WriteTOMLConfig(a.configPath, cfg)
}

func (a *CodexAgent) HasMCP(name string) (bool, error) {
	cfg, err := config.ReadTOMLConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return hasCodexMCP(cfg, name), nil
}

func (a *CodexAgent) InstallMCP(name string, mcpConfig map[string]interface{}) error {
	if mcpConfig == nil {
		return fmt.Errorf("missing mcp config")
	}
	cfg, err := config.ReadTOMLConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			cfg = make(map[string]interface{})
		} else {
			return err
		}
	}
	addCodexMCP(cfg, name, normalizeArgsToStrings(mcpConfig))
	return config.WriteTOMLConfig(a.configPath, cfg)
}

func (a *CodexAgent) RemoveMCP(name string) error {
	cfg, err := config.ReadTOMLConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	removeCodexMCP(cfg, name)
	return config.WriteTOMLConfig(a.configPath, cfg)
}

func (a *CodexAgent) ListMCPs() (map[string]map[string]interface{}, error) {
	cfg, err := config.ReadTOMLConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]map[string]interface{}{}, nil
		}
		return nil, err
	}
	result := map[string]map[string]interface{}{}
	mcpServers, ok := cfg[codexMCPKey].(map[string]interface{})
	if !ok {
		return result, nil
	}
	for name, raw := range mcpServers {
		if raw == nil {
			continue
		}
		if serverCfg, ok := raw.(map[string]interface{}); ok {
			result[name] = serverCfg
		}
	}
	return result, nil
}

func (a *CodexAgent) SupportsSkills() bool {
	return true
}

func (a *CodexAgent) HasSkill(skillName string) (bool, error) {
	mgr := skills.NewCodexSkillManager()
	skill, err := mgr.Get(skillName)
	if err != nil {
		return false, nil // Not found is not an error
	}
	return skill != nil, nil
}

func (a *CodexAgent) InstallSkill(skillName, source string) error {
	mgr := skills.NewCodexSkillManager()
	_, err := mgr.Install(source, skills.ScopePersonal)
	return err
}

func (a *CodexAgent) RemoveSkill(skillName string) error {
	mgr := skills.NewCodexSkillManager()
	return mgr.Remove(skillName, skills.ScopePersonal)
}

func (a *CodexAgent) SupportsPlugins() bool {
	return false
}

func (a *CodexAgent) HasPlugin(pluginName string) (bool, error) {
	return false, nil
}

func (a *CodexAgent) InstallPlugin(pluginName, source string) error {
	return nil
}

func (a *CodexAgent) RemovePlugin(pluginName string) error {
	return nil
}

func hasCodexMCP(cfg map[string]interface{}, name string) bool {
	mcpServers, ok := cfg[codexMCPKey].(map[string]interface{})
	if !ok {
		return false
	}
	_, exists := mcpServers[name]
	return exists
}

func addCodexMCP(cfg map[string]interface{}, name string, mcpConfig map[string]interface{}) {
	mcpServers, ok := cfg[codexMCPKey].(map[string]interface{})
	if !ok {
		mcpServers = make(map[string]interface{})
		cfg[codexMCPKey] = mcpServers
	}
	mcpServers[name] = mcpConfig
}

func removeCodexMCP(cfg map[string]interface{}, name string) {
	mcpServers, ok := cfg[codexMCPKey].(map[string]interface{})
	if !ok {
		return
	}
	delete(mcpServers, name)
}
