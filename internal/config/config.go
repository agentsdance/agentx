package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// PlaywrightMCPConfig is the configuration for Playwright MCP server
var PlaywrightMCPConfig = map[string]interface{}{
	"command": "npx",
	"args":    []interface{}{"@playwright/mcp@latest"},
}

// Context7MCPConfig is the configuration for Context7 MCP server
var Context7MCPConfig = map[string]interface{}{
	"command": "npx",
	"args":    []interface{}{"-y", "@upstash/context7-mcp"},
}

// Context7MCPConfigRemote is the remote configuration for Context7 MCP server
var Context7MCPConfigRemote = map[string]interface{}{
	"url": "https://mcp.context7.com/mcp",
}

// ReadConfig reads a JSON config file
func ReadConfig(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg map[string]interface{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// WriteConfig writes a JSON config file with pretty formatting
func WriteConfig(path string, cfg map[string]interface{}) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// HasMCP checks if a specific MCP server is configured
func HasMCP(cfg map[string]interface{}, name string) bool {
	mcpServers, ok := cfg["mcpServers"].(map[string]interface{})
	if !ok {
		return false
	}
	_, exists := mcpServers[name]
	return exists
}

// AddMCP adds an MCP server to the config
func AddMCP(cfg map[string]interface{}, name string, mcpConfig map[string]interface{}) {
	mcpServers, ok := cfg["mcpServers"].(map[string]interface{})
	if !ok {
		mcpServers = make(map[string]interface{})
		cfg["mcpServers"] = mcpServers
	}
	mcpServers[name] = mcpConfig
}

// RemoveMCP removes an MCP server from the config
func RemoveMCP(cfg map[string]interface{}, name string) {
	mcpServers, ok := cfg["mcpServers"].(map[string]interface{})
	if !ok {
		return
	}
	delete(mcpServers, name)
}

// HasPlaywrightMCP checks if Playwright MCP is configured
func HasPlaywrightMCP(cfg map[string]interface{}) bool {
	return HasMCP(cfg, "playwright")
}

// AddPlaywrightMCP adds Playwright MCP to the config
func AddPlaywrightMCP(cfg map[string]interface{}) {
	AddMCP(cfg, "playwright", PlaywrightMCPConfig)
}

// RemovePlaywrightMCP removes Playwright MCP from the config
func RemovePlaywrightMCP(cfg map[string]interface{}) {
	RemoveMCP(cfg, "playwright")
}

// HasContext7MCP checks if Context7 MCP is configured
func HasContext7MCP(cfg map[string]interface{}) bool {
	return HasMCP(cfg, "context7")
}

// AddContext7MCP adds Context7 MCP to the config
func AddContext7MCP(cfg map[string]interface{}) {
	AddMCP(cfg, "context7", Context7MCPConfig)
}

// RemoveContext7MCP removes Context7 MCP from the config
func RemoveContext7MCP(cfg map[string]interface{}) {
	RemoveMCP(cfg, "context7")
}
