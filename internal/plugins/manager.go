package plugins

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/agentsdance/agentx/internal/skills"
)

// DefaultPluginManager implements PluginManager
type DefaultPluginManager struct{}

// NewPluginManager creates a new plugin manager
func NewPluginManager() *DefaultPluginManager {
	return &DefaultPluginManager{}
}

// List returns all installed plugins
func (m *DefaultPluginManager) List() ([]Plugin, error) {
	pluginsDir, err := GetPluginsDir()
	if err != nil {
		return nil, err
	}

	var plugins []Plugin

	entries, err := os.ReadDir(pluginsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return plugins, nil
		}
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			pluginPath := filepath.Join(pluginsDir, entry.Name())
			if IsPluginDir(pluginPath) {
				plugin, err := m.parsePluginDir(pluginPath)
				if err == nil {
					plugins = append(plugins, *plugin)
				}
			}
		}
	}

	return plugins, nil
}

// Install installs a plugin from a source
func (m *DefaultPluginManager) Install(source string) (*Plugin, error) {
	info, err := ParseSource(source)
	if err != nil {
		return nil, err
	}

	switch info.Type {
	case SourceTypeLocal:
		return m.installFromLocal(info.Path)
	case SourceTypeGitRepo:
		return m.installFromGit(info.RepoURL, "", "")
	case SourceTypeGitRepoWithFragment:
		pluginPath := info.PluginPath
		fragment := info.Fragment
		return m.installFromGit(info.RepoURL, fragment, pluginPath)
	default:
		return nil, fmt.Errorf("unsupported source type")
	}
}

// Remove removes a plugin by name
func (m *DefaultPluginManager) Remove(name string) error {
	pluginsDir, err := GetPluginsDir()
	if err != nil {
		return err
	}

	pluginPath := filepath.Join(pluginsDir, name)
	if _, err := os.Stat(pluginPath); err != nil {
		return fmt.Errorf("plugin not found: %s", name)
	}

	return os.RemoveAll(pluginPath)
}

// Check verifies plugins installation status
func (m *DefaultPluginManager) Check() ([]PluginStatus, error) {
	plugins, err := m.List()
	if err != nil {
		return nil, err
	}

	var statuses []PluginStatus
	for _, plugin := range plugins {
		status := PluginStatus{Plugin: plugin, Valid: true}

		// Validate plugin manifest exists
		manifestPath := GetPluginManifestPath(plugin.Path)
		if _, err := os.Stat(manifestPath); err != nil {
			status.Valid = false
			status.Error = err
			status.Issues = append(status.Issues, "Manifest not found")
		}

		// Validate required fields
		if plugin.Name == "" {
			status.Issues = append(status.Issues, "Missing name in manifest")
		}

		if plugin.Version == "" {
			status.Issues = append(status.Issues, "Missing version in manifest")
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}

// Get retrieves a plugin by name
func (m *DefaultPluginManager) Get(name string) (*Plugin, error) {
	plugins, err := m.List()
	if err != nil {
		return nil, err
	}

	for _, plugin := range plugins {
		if plugin.Name == name {
			return &plugin, nil
		}
	}

	return nil, fmt.Errorf("plugin not found: %s", name)
}

// parsePluginDir parses a plugin directory and returns a Plugin
func (m *DefaultPluginManager) parsePluginDir(path string) (*Plugin, error) {
	manifestPath := GetPluginManifestPath(path)
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}

	var manifest PluginManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("invalid plugin.json: %w", err)
	}

	// Use directory name if name not in manifest
	name := manifest.Name
	if name == "" {
		name = filepath.Base(path)
	}

	components := m.scanComponents(path)

	return &Plugin{
		Name:        name,
		Version:     manifest.Version,
		Description: manifest.Description,
		Author:      manifest.Author,
		Path:        path,
		Components:  components,
	}, nil
}

// scanComponents scans a plugin directory for its components
func (m *DefaultPluginManager) scanComponents(pluginPath string) PluginComponents {
	var components PluginComponents

	// Scan commands
	commandsDir := GetPluginCommandsDir(pluginPath)
	if entries, err := os.ReadDir(commandsDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
				name := strings.TrimSuffix(entry.Name(), ".md")
				components.Commands = append(components.Commands, name)
			}
		}
	}

	// Scan agents
	agentsDir := GetPluginAgentsDir(pluginPath)
	if entries, err := os.ReadDir(agentsDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
				name := strings.TrimSuffix(entry.Name(), ".md")
				components.Agents = append(components.Agents, name)
			}
		}
	}

	// Scan skills
	skillsDir := GetPluginSkillsDir(pluginPath)
	if entries, err := os.ReadDir(skillsDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				skillPath := filepath.Join(skillsDir, entry.Name())
				skillMD := filepath.Join(skillPath, "SKILL.md")
				if _, err := os.Stat(skillMD); err == nil {
					components.Skills = append(components.Skills, entry.Name())
				}
			}
		}
	}

	// Scan hooks
	hooksDir := GetPluginHooksDir(pluginPath)
	if entries, err := os.ReadDir(hooksDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
				name := strings.TrimSuffix(entry.Name(), ".json")
				components.Hooks = append(components.Hooks, name)
			}
		}
	}

	// Scan MCP servers
	mcpPath := GetPluginMCPPath(pluginPath)
	if data, err := os.ReadFile(mcpPath); err == nil {
		var mcpConfig map[string]interface{}
		if err := json.Unmarshal(data, &mcpConfig); err == nil {
			if servers, ok := mcpConfig["mcpServers"].(map[string]interface{}); ok {
				for name := range servers {
					components.MCPServers = append(components.MCPServers, name)
				}
			}
		}
	}

	return components
}

// installFromLocal installs a plugin from a local path
func (m *DefaultPluginManager) installFromLocal(sourcePath string) (*Plugin, error) {
	info, err := os.Stat(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("source not found: %s", sourcePath)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("plugin source must be a directory: %s", sourcePath)
	}

	if !IsPluginDir(sourcePath) {
		return nil, fmt.Errorf("not a valid plugin directory (missing .claude-plugin/plugin.json): %s", sourcePath)
	}

	return m.installPluginDir(sourcePath)
}

// installFromGit installs a plugin from a git repository
func (m *DefaultPluginManager) installFromGit(repoURL, fragment, pluginPath string) (*Plugin, error) {
	// Clone the repository using skills package
	tmpDir, err := skills.GitClone(repoURL)
	if err != nil {
		return nil, err
	}
	defer skills.CleanupTempDir(tmpDir)

	// Determine the plugin location
	var targetPath string

	if pluginPath != "" {
		// For tree URLs, use the full path within the repo
		targetPath = filepath.Join(tmpDir, pluginPath)
	} else if fragment != "" {
		// Try to find plugin by fragment name
		targetPath = filepath.Join(tmpDir, fragment)
		if !IsPluginDir(targetPath) {
			// Try in plugins/ subdirectory
			targetPath = filepath.Join(tmpDir, "plugins", fragment)
		}
	} else {
		// Check if repo root is a plugin
		targetPath = tmpDir
	}

	if !IsPluginDir(targetPath) {
		return nil, fmt.Errorf("plugin not found at path: %s", targetPath)
	}

	plugin, err := m.installPluginDir(targetPath)
	if err != nil {
		return nil, err
	}

	// Set source URL
	plugin.Source = repoURL
	if fragment != "" {
		plugin.Source = repoURL + "#" + fragment
	} else if pluginPath != "" {
		plugin.Source = repoURL + "/tree/main/" + pluginPath
	}

	return plugin, nil
}

// installPluginDir installs a plugin from a directory
func (m *DefaultPluginManager) installPluginDir(sourcePath string) (*Plugin, error) {
	// Parse the plugin first to get its name
	plugin, err := m.parsePluginDir(sourcePath)
	if err != nil {
		return nil, err
	}

	// Get target directory
	pluginsDir, err := GetPluginsDir()
	if err != nil {
		return nil, err
	}

	// Ensure plugins directory exists
	if err := os.MkdirAll(pluginsDir, 0755); err != nil {
		return nil, err
	}

	// Copy the plugin directory
	targetPath := filepath.Join(pluginsDir, plugin.Name)

	// Check if already exists
	if _, err := os.Stat(targetPath); err == nil {
		return nil, fmt.Errorf("plugin already exists: %s", plugin.Name)
	}

	if err := copyDir(sourcePath, targetPath); err != nil {
		return nil, fmt.Errorf("failed to copy plugin: %w", err)
	}

	plugin.Path = targetPath
	return plugin, nil
}

// copyDir copies a directory recursively
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip .git directory
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		// Calculate relative path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		return copyFile(path, targetPath)
	})
}

// copyFile copies a single file
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Get source file info for permissions
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
