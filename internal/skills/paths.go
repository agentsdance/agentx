package skills

import (
	"os"
	"path/filepath"
)

// GetClaudeBasePaths returns the base paths for Claude Code configuration
func GetClaudeBasePaths() (personal, project string, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", "", err
	}

	personal = filepath.Join(home, ".claude")

	// Project path is relative to current directory
	cwd, err := os.Getwd()
	if err != nil {
		return personal, "", err
	}
	project = filepath.Join(cwd, ".claude")

	return personal, project, nil
}

// GetCommandsDir returns the commands directory for a scope
func GetCommandsDir(scope SkillScope) (string, error) {
	personal, project, err := GetClaudeBasePaths()
	if err != nil {
		return "", err
	}

	base := personal
	if scope == ScopeProject {
		base = project
	}
	return filepath.Join(base, "commands"), nil
}

// GetSkillsDir returns the skills directory for a scope
func GetSkillsDir(scope SkillScope) (string, error) {
	personal, project, err := GetClaudeBasePaths()
	if err != nil {
		return "", err
	}

	base := personal
	if scope == ScopeProject {
		base = project
	}
	return filepath.Join(base, "skills"), nil
}

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}
