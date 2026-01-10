package plugins

import (
	"os"
	"path/filepath"
)

// GetPluginsDir returns the plugins directory path (~/.agentx/plugins)
func GetPluginsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".agentx", "plugins"), nil
}

// GetPluginManifestPath returns the path to plugin.json for a plugin
func GetPluginManifestPath(pluginPath string) string {
	return filepath.Join(pluginPath, ".claude-plugin", "plugin.json")
}

// IsPluginDir checks if a directory is a valid plugin directory
func IsPluginDir(path string) bool {
	manifestPath := GetPluginManifestPath(path)
	_, err := os.Stat(manifestPath)
	return err == nil
}

// GetPluginCommandsDir returns the path to the commands directory
func GetPluginCommandsDir(pluginPath string) string {
	return filepath.Join(pluginPath, "commands")
}

// GetPluginAgentsDir returns the path to the agents directory
func GetPluginAgentsDir(pluginPath string) string {
	return filepath.Join(pluginPath, "agents")
}

// GetPluginSkillsDir returns the path to the skills directory
func GetPluginSkillsDir(pluginPath string) string {
	return filepath.Join(pluginPath, "skills")
}

// GetPluginHooksDir returns the path to the hooks directory
func GetPluginHooksDir(pluginPath string) string {
	return filepath.Join(pluginPath, "hooks")
}

// GetPluginMCPPath returns the path to .mcp.json
func GetPluginMCPPath(pluginPath string) string {
	return filepath.Join(pluginPath, ".mcp.json")
}
