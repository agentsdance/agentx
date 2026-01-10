package plugins

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// DefaultRegistryURL is the default URL for the plugin registry
const DefaultRegistryURL = "https://raw.githubusercontent.com/agentsdance/agentx/master/registry/plugins.json"

// RegistryPlugin represents a plugin entry in the registry
type RegistryPlugin struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Source      string   `json:"source"`
	Version     string   `json:"version"`
	Author      string   `json:"author"`
	Components  []string `json:"components"` // Summary like ["commands", "skills", "mcp"]
}

// Registry represents the plugin registry
type Registry struct {
	Plugins []RegistryPlugin `json:"plugins"`
}

// FetchRegistry fetches the plugin registry from the default URL
func FetchRegistry() ([]RegistryPlugin, error) {
	return FetchRegistryFromURL(DefaultRegistryURL)
}

// FetchRegistryFromURL fetches the plugin registry from a specific URL
func FetchRegistryFromURL(registryURL string) ([]RegistryPlugin, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(registryURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registry returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read registry response: %w", err)
	}

	var registry Registry
	if err := json.Unmarshal(body, &registry); err != nil {
		return nil, fmt.Errorf("failed to parse registry: %w", err)
	}

	// Cache the registry locally
	if err := cacheRegistry(body); err != nil {
		// Non-fatal, just log if needed
	}

	return registry.Plugins, nil
}

// GetCachedRegistry returns the cached registry if available
func GetCachedRegistry() ([]RegistryPlugin, error) {
	cachePath, err := getRegistryCachePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, err
	}

	var registry Registry
	if err := json.Unmarshal(data, &registry); err != nil {
		return nil, err
	}

	return registry.Plugins, nil
}

// FetchRegistryWithFallback tries to fetch from network, falls back to cache, then to local file
func FetchRegistryWithFallback() ([]RegistryPlugin, error) {
	plugins, err := FetchRegistry()
	if err == nil && len(plugins) > 0 {
		return plugins, nil
	}

	// Try cached version
	cached, cacheErr := GetCachedRegistry()
	if cacheErr == nil && len(cached) > 0 {
		return cached, nil
	}

	// Try local bundled registry file (for development)
	local, localErr := GetLocalRegistry()
	if localErr == nil && len(local) > 0 {
		return local, nil
	}

	if err != nil {
		return nil, err
	}
	return plugins, nil
}

// GetLocalRegistry reads the registry from the local bundled file
func GetLocalRegistry() ([]RegistryPlugin, error) {
	// Try common locations for the registry file
	paths := []string{
		"registry/plugins.json",                    // Current working directory
		"./registry/plugins.json",                  // Explicit current directory
		filepath.Join("..", "registry/plugins.json"), // Parent directory
	}

	// Also try relative to executable
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		paths = append(paths, filepath.Join(exeDir, "registry", "plugins.json"))
		paths = append(paths, filepath.Join(exeDir, "..", "registry", "plugins.json"))
	}

	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var registry Registry
		if err := json.Unmarshal(data, &registry); err != nil {
			continue
		}

		return registry.Plugins, nil
	}

	return nil, fmt.Errorf("local registry not found")
}

// cacheRegistry saves the registry data to local cache
func cacheRegistry(data []byte) error {
	cachePath, err := getRegistryCachePath()
	if err != nil {
		return err
	}

	// Ensure cache directory exists
	cacheDir := filepath.Dir(cachePath)
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return err
	}

	return os.WriteFile(cachePath, data, 0644)
}

// getRegistryCachePath returns the path to the cached registry
func getRegistryCachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".agentx", "cache", "plugin-registry.json"), nil
}

// ComponentsSummary returns a human-readable summary of components
func ComponentsSummary(components PluginComponents) string {
	var parts []string

	if len(components.Commands) > 0 {
		parts = append(parts, fmt.Sprintf("%d cmd", len(components.Commands)))
	}
	if len(components.Agents) > 0 {
		parts = append(parts, fmt.Sprintf("%d agent", len(components.Agents)))
	}
	if len(components.Skills) > 0 {
		parts = append(parts, fmt.Sprintf("%d skill", len(components.Skills)))
	}
	if len(components.Hooks) > 0 {
		parts = append(parts, fmt.Sprintf("%d hook", len(components.Hooks)))
	}
	if len(components.MCPServers) > 0 {
		parts = append(parts, fmt.Sprintf("%d mcp", len(components.MCPServers)))
	}

	if len(parts) == 0 {
		return "empty"
	}

	return joinParts(parts)
}

func joinParts(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		return parts[0]
	}

	result := parts[0]
	for i := 1; i < len(parts); i++ {
		result += ", " + parts[i]
	}
	return result
}
