package skills

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// DefaultSkillsRegistryURL is the default URL for the skills registry
const DefaultSkillsRegistryURL = "https://raw.githubusercontent.com/agentsdance/agentskills/master/skills.json"

// RegistrySkill represents a skill entry in the registry
type RegistrySkill struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Author      string `json:"author,omitempty"`
	License     string `json:"license,omitempty"`
	Source      string `json:"source"`
}

// SkillsRegistry represents the skills registry
type SkillsRegistry struct {
	Version string          `json:"version"`
	Skills  []RegistrySkill `json:"skills"`
}

// FetchSkillsRegistry fetches the skills registry from the default URL
func FetchSkillsRegistry() ([]RegistrySkill, error) {
	return FetchSkillsRegistryFromURL(DefaultSkillsRegistryURL)
}

// FetchSkillsRegistryFromURL fetches the skills registry from a specific URL
func FetchSkillsRegistryFromURL(registryURL string) ([]RegistrySkill, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(registryURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch skills registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("skills registry returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read skills registry response: %w", err)
	}

	var registry SkillsRegistry
	if err := json.Unmarshal(body, &registry); err != nil {
		return nil, fmt.Errorf("failed to parse skills registry: %w", err)
	}

	// Cache the registry locally
	if err := cacheSkillsRegistry(body); err != nil {
		// Non-fatal, just log if needed
	}

	return registry.Skills, nil
}

// GetCachedSkillsRegistry returns the cached skills registry if available
func GetCachedSkillsRegistry() ([]RegistrySkill, error) {
	cachePath, err := getSkillsRegistryCachePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, err
	}

	var registry SkillsRegistry
	if err := json.Unmarshal(data, &registry); err != nil {
		return nil, err
	}

	return registry.Skills, nil
}

// FetchSkillsRegistryWithFallback tries to fetch from network, falls back to cache, then to local file
func FetchSkillsRegistryWithFallback() ([]RegistrySkill, error) {
	skills, err := FetchSkillsRegistry()
	if err == nil && len(skills) > 0 {
		return skills, nil
	}

	// Try cached version
	cached, cacheErr := GetCachedSkillsRegistry()
	if cacheErr == nil && len(cached) > 0 {
		return cached, nil
	}

	// Try local bundled registry file (for development)
	local, localErr := GetLocalSkillsRegistry()
	if localErr == nil && len(local) > 0 {
		return local, nil
	}

	if err != nil {
		return nil, err
	}
	return skills, nil
}

// GetLocalSkillsRegistry reads the skills registry from the local bundled file
func GetLocalSkillsRegistry() ([]RegistrySkill, error) {
	// Try common locations for the registry file
	paths := []string{
		"registry/skills.json",                       // Current working directory
		"./registry/skills.json",                     // Explicit current directory
		filepath.Join("..", "registry/skills.json"), // Parent directory
	}

	// Also try relative to executable
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		paths = append(paths, filepath.Join(exeDir, "registry", "skills.json"))
		paths = append(paths, filepath.Join(exeDir, "..", "registry", "skills.json"))
	}

	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var registry SkillsRegistry
		if err := json.Unmarshal(data, &registry); err != nil {
			continue
		}

		return registry.Skills, nil
	}

	return nil, fmt.Errorf("local skills registry not found")
}

// cacheSkillsRegistry saves the skills registry data to local cache
func cacheSkillsRegistry(data []byte) error {
	cachePath, err := getSkillsRegistryCachePath()
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

// getSkillsRegistryCachePath returns the path to the cached skills registry
func getSkillsRegistryCachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".agentx", "cache", "skills-registry.json"), nil
}
