package skills

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// SourceType represents the type of skill source
type SourceType int

const (
	// SourceTypeLocal is a local filesystem path
	SourceTypeLocal SourceType = iota
	// SourceTypeGitRepo is a git repository URL
	SourceTypeGitRepo
	// SourceTypeGitRepoWithFragment is a git URL with #skill-name fragment
	SourceTypeGitRepoWithFragment
)

// SourceInfo contains parsed source information
type SourceInfo struct {
	Type     SourceType
	Path     string // Local path (absolute)
	RepoURL  string // Git repository URL
	Fragment string // Skill name for repo#skill-name format
}

// ParseSource parses a skill source string into SourceInfo
func ParseSource(source string) (*SourceInfo, error) {
	// Check if it's a local path first
	// Try to resolve it as an absolute path
	absPath, err := filepath.Abs(source)
	if err == nil {
		if info, err := os.Stat(absPath); err == nil {
			if info.IsDir() {
				return &SourceInfo{
					Type: SourceTypeLocal,
					Path: absPath,
				}, nil
			}
			// It's a file, could be a single command .md file
			return &SourceInfo{
				Type: SourceTypeLocal,
				Path: absPath,
			}, nil
		}
	}

	// Check for fragment (repo#skill-name)
	fragment := ""
	sourceWithoutFragment := source
	if idx := strings.LastIndex(source, "#"); idx > 0 {
		// Make sure it's not part of the URL path (like github.com/user#section)
		// by checking if there's a / after the #
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
