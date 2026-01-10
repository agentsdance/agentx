package plugins

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// SourceType represents the type of plugin source
type SourceType int

const (
	// SourceTypeLocal is a local filesystem path
	SourceTypeLocal SourceType = iota
	// SourceTypeGitRepo is a git repository URL
	SourceTypeGitRepo
	// SourceTypeGitRepoWithFragment is a git URL with #plugin-name fragment
	SourceTypeGitRepoWithFragment
)

// SourceInfo contains parsed source information
type SourceInfo struct {
	Type       SourceType
	Path       string // Local path (absolute)
	RepoURL    string // Git repository URL
	Fragment   string // Plugin name for repo#plugin-name format
	PluginPath string // Path within repo (for tree URLs)
}

// parseGitHubTreeURL parses a GitHub tree URL like:
// https://github.com/org/repo/tree/branch/path/to/plugin
// Returns repoURL, pluginPath, and whether it was a tree URL
func parseGitHubTreeURL(source string) (repoURL string, pluginPath string, isTreeURL bool) {
	if !strings.Contains(source, "github.com") {
		return "", "", false
	}

	u, err := url.Parse(source)
	if err != nil {
		return "", "", false
	}

	// Path should be like: /org/repo/tree/branch/path/to/plugin
	parts := strings.Split(strings.TrimPrefix(u.Path, "/"), "/")
	if len(parts) < 4 || parts[2] != "tree" {
		return "", "", false
	}

	// Extract org, repo
	org := parts[0]
	repo := parts[1]

	// Plugin path is everything after the branch
	if len(parts) > 4 {
		pluginPath = strings.Join(parts[4:], "/")
	}

	repoURL = fmt.Sprintf("https://github.com/%s/%s", org, repo)
	return repoURL, pluginPath, true
}

// ParseSource parses a plugin source string into SourceInfo
func ParseSource(source string) (*SourceInfo, error) {
	// Check if it's a local path first
	absPath, err := filepath.Abs(source)
	if err == nil {
		if info, err := os.Stat(absPath); err == nil {
			if info.IsDir() && IsPluginDir(absPath) {
				return &SourceInfo{
					Type: SourceTypeLocal,
					Path: absPath,
				}, nil
			}
			// Even if it doesn't look like a plugin yet, allow local dirs
			if info.IsDir() {
				return &SourceInfo{
					Type: SourceTypeLocal,
					Path: absPath,
				}, nil
			}
		}
	}

	// Check if it's a GitHub tree URL first
	if repoURL, pluginPath, isTreeURL := parseGitHubTreeURL(source); isTreeURL {
		if pluginPath != "" {
			return &SourceInfo{
				Type:       SourceTypeGitRepoWithFragment,
				RepoURL:    repoURL,
				PluginPath: pluginPath,
			}, nil
		}
		return &SourceInfo{
			Type:    SourceTypeGitRepo,
			RepoURL: repoURL,
		}, nil
	}

	// Check for fragment (repo#plugin-name)
	fragment := ""
	sourceWithoutFragment := source
	if idx := strings.LastIndex(source, "#"); idx > 0 {
		potentialFragment := source[idx+1:]
		if !strings.Contains(potentialFragment, "/") {
			fragment = potentialFragment
			sourceWithoutFragment = source[:idx]
		}
	}

	// Parse as URL
	u, err := url.Parse(sourceWithoutFragment)
	if err != nil {
		return nil, fmt.Errorf("invalid source: %s", source)
	}

	// Check if it's a valid HTTP(S) URL
	if u.Scheme == "https" || u.Scheme == "http" {
		if fragment != "" {
			return &SourceInfo{
				Type:     SourceTypeGitRepoWithFragment,
				RepoURL:  sourceWithoutFragment,
				Fragment: fragment,
			}, nil
		}
		return &SourceInfo{
			Type:    SourceTypeGitRepo,
			RepoURL: sourceWithoutFragment,
		}, nil
	}

	// Check for git@ style URLs (SSH)
	if strings.HasPrefix(source, "git@") {
		if fragment != "" {
			return &SourceInfo{
				Type:     SourceTypeGitRepoWithFragment,
				RepoURL:  sourceWithoutFragment,
				Fragment: fragment,
			}, nil
		}
		return &SourceInfo{
			Type:    SourceTypeGitRepo,
			RepoURL: sourceWithoutFragment,
		}, nil
	}

	return nil, fmt.Errorf("cannot determine source type for: %s", source)
}

// IsGitSource returns true if the source is a git repository
func (s *SourceInfo) IsGitSource() bool {
	return s.Type == SourceTypeGitRepo || s.Type == SourceTypeGitRepoWithFragment
}

// IsLocalSource returns true if the source is a local path
func (s *SourceInfo) IsLocalSource() bool {
	return s.Type == SourceTypeLocal
}
