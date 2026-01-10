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
	Type      SourceType
	Path      string // Local path (absolute)
	RepoURL   string // Git repository URL
	Fragment  string // Skill name for repo#skill-name format
	SkillPath string // Path within repo (for tree URLs like github.com/org/repo/tree/branch/path/to/skill)
}

// parseGitHubTreeURL parses a GitHub tree URL like:
// https://github.com/org/repo/tree/branch/path/to/skill
// Returns repoURL, skillPath, and whether it was a tree URL
func parseGitHubTreeURL(source string) (repoURL string, skillPath string, isTreeURL bool) {
	// Check if it's a GitHub tree URL
	if !strings.Contains(source, "github.com") {
		return "", "", false
	}

	u, err := url.Parse(source)
	if err != nil {
		return "", "", false
	}

	// Path should be like: /org/repo/tree/branch/path/to/skill
	parts := strings.Split(strings.TrimPrefix(u.Path, "/"), "/")
	if len(parts) < 4 || parts[2] != "tree" {
		return "", "", false
	}

	// Extract org, repo, branch, and skill path
	org := parts[0]
	repo := parts[1]
	// branch := parts[3] // We don't need branch for cloning (will use default)

	// Skill path is everything after the branch
	if len(parts) > 4 {
		skillPath = strings.Join(parts[4:], "/")
	}

	repoURL = fmt.Sprintf("https://github.com/%s/%s", org, repo)
	return repoURL, skillPath, true
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

	// Check if it's a GitHub tree URL first
	if repoURL, skillPath, isTreeURL := parseGitHubTreeURL(source); isTreeURL {
		if skillPath != "" {
			return &SourceInfo{
				Type:      SourceTypeGitRepoWithFragment,
				RepoURL:   repoURL,
				SkillPath: skillPath,
			}, nil
		}
		return &SourceInfo{
			Type:    SourceTypeGitRepo,
			RepoURL: repoURL,
		}, nil
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
