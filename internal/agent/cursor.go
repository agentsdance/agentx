package agent

import (
	"os"
	"path/filepath"

	"github.com/agentsdance/agentx/internal/config"
)

// CursorAgent represents Cursor editor agent
type CursorAgent struct {
	configPath string
}

// NewCursorAgent creates a new Cursor agent
func NewCursorAgent() *CursorAgent {
	home, _ := os.UserHomeDir()
	return &CursorAgent{
		configPath: filepath.Join(home, ".cursor", "mcp.json"),
	}
}

func (a *CursorAgent) Name() string {
	return "Cursor"
}

func (a *CursorAgent) ConfigPath() string {
	return a.configPath
}

func (a *CursorAgent) Exists() bool {
	// Check if .cursor directory exists
	dir := filepath.Dir(a.configPath)
	_, err := os.Stat(dir)
	return err == nil
}

func (a *CursorAgent) HasPlaywright() (bool, error) {
	return a.HasMCP("playwright")
}

func (a *CursorAgent) InstallPlaywright() error {
	return a.InstallMCP("playwright", config.PlaywrightMCPConfig)
}

func (a *CursorAgent) RemovePlaywright() error {
	return a.RemoveMCP("playwright")
}

func (a *CursorAgent) HasContext7() (bool, error) {
	return a.HasMCP("context7")
}

func (a *CursorAgent) InstallContext7() error {
	return a.InstallMCP("context7", config.Context7MCPConfig)
}

func (a *CursorAgent) RemoveContext7() error {
	return a.RemoveMCP("context7")
}

// HasMCP checks if a specific MCP server is configured
func (a *CursorAgent) HasMCP(name string) (bool, error) {
	cfg, err := config.ReadConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return config.HasMCP(cfg, name), nil
}

// InstallMCP adds an MCP server to the config
func (a *CursorAgent) InstallMCP(name string, mcpConfig map[string]interface{}) error {
	cfg, err := config.ReadConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			cfg = make(map[string]interface{})
		} else {
			return err
		}
	}
	config.AddMCP(cfg, name, mcpConfig)
	return config.WriteConfig(a.configPath, cfg)
}

// RemoveMCP removes an MCP server from the config
func (a *CursorAgent) RemoveMCP(name string) error {
	cfg, err := config.ReadConfig(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	config.RemoveMCP(cfg, name)
	return config.WriteConfig(a.configPath, cfg)
}

func (a *CursorAgent) SupportsSkills() bool {
	return false
}

func (a *CursorAgent) HasSkill(skillName string) (bool, error) {
	return false, nil
}

func (a *CursorAgent) InstallSkill(skillName, source string) error {
	return nil
}

func (a *CursorAgent) RemoveSkill(skillName string) error {
	return nil
}
