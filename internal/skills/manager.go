package skills

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// DefaultSkillManager implements SkillManager
type DefaultSkillManager struct{}

// NewSkillManager creates a new skill manager
func NewSkillManager() *DefaultSkillManager {
	return &DefaultSkillManager{}
}

// List returns all installed skills from both personal and project scopes
func (m *DefaultSkillManager) List() ([]Skill, error) {
	var skills []Skill

	for _, scope := range []SkillScope{ScopePersonal, ScopeProject} {
		scopeSkills, err := m.ListByScope(scope)
		if err != nil {
			// Skip errors for missing directories
			continue
		}
		skills = append(skills, scopeSkills...)
	}

	return skills, nil
}

// ListByScope returns skills filtered by scope
func (m *DefaultSkillManager) ListByScope(scope SkillScope) ([]Skill, error) {
	var skills []Skill

	// List commands (single .md files)
	commandsDir, err := GetCommandsDir(scope)
	if err != nil {
		return nil, err
	}

	if entries, err := os.ReadDir(commandsDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
				skill, err := m.parseCommandFile(filepath.Join(commandsDir, entry.Name()), scope)
				if err == nil {
					skills = append(skills, *skill)
				}
			}
		}
	}

	// List skills (directories with SKILL.md)
	skillsDir, err := GetSkillsDir(scope)
	if err != nil {
		return nil, err
	}

	if entries, err := os.ReadDir(skillsDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				skillPath := filepath.Join(skillsDir, entry.Name())
				skillMD := filepath.Join(skillPath, "SKILL.md")
				if _, err := os.Stat(skillMD); err == nil {
					skill, err := m.parseSkillDir(skillPath, scope)
					if err == nil {
						skills = append(skills, *skill)
					}
				}
			}
		}
	}

	return skills, nil
}

// Install installs a skill from a source
func (m *DefaultSkillManager) Install(source string, scope SkillScope) (*Skill, error) {
	info, err := ParseSource(source)
	if err != nil {
		return nil, err
	}

	switch info.Type {
	case SourceTypeLocal:
		return m.installFromLocal(info.Path, scope)
	case SourceTypeGitRepo:
		return m.installFromGit(info.RepoURL, "", "", scope)
	case SourceTypeGitRepoWithFragment:
		// Use SkillPath if available (from tree URLs), otherwise use Fragment
		skillPath := info.SkillPath
		fragment := info.Fragment
		return m.installFromGit(info.RepoURL, fragment, skillPath, scope)
	default:
		return nil, fmt.Errorf("unsupported source type")
	}
}

// Remove removes a skill by name
func (m *DefaultSkillManager) Remove(name string, scope SkillScope) error {
	// Try removing from commands
	commandsDir, _ := GetCommandsDir(scope)
	commandPath := filepath.Join(commandsDir, name+".md")
	if _, err := os.Stat(commandPath); err == nil {
		return os.Remove(commandPath)
	}

	// Try removing from skills
	skillsDir, _ := GetSkillsDir(scope)
	skillPath := filepath.Join(skillsDir, name)
	if _, err := os.Stat(skillPath); err == nil {
		return os.RemoveAll(skillPath)
	}

	return fmt.Errorf("skill not found: %s", name)
}

// Check verifies skills installation status
func (m *DefaultSkillManager) Check() ([]SkillStatus, error) {
	skills, err := m.List()
	if err != nil {
		return nil, err
	}

	var statuses []SkillStatus
	for _, skill := range skills {
		status := SkillStatus{Skill: skill, Valid: true}

		// Validate skill file exists
		if _, err := os.Stat(skill.Path); err != nil {
			status.Valid = false
			status.Error = err
			status.Issues = append(status.Issues, "File not found")
		}

		// Validate frontmatter
		if skill.Description == "" {
			status.Issues = append(status.Issues, "Missing description")
		}

		if skill.Name == "" {
			status.Issues = append(status.Issues, "Missing name")
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}

// Get retrieves a skill by name
func (m *DefaultSkillManager) Get(name string) (*Skill, error) {
	skills, err := m.List()
	if err != nil {
		return nil, err
	}

	for _, skill := range skills {
		if skill.Name == name {
			return &skill, nil
		}
	}

	return nil, fmt.Errorf("skill not found: %s", name)
}

// parseCommandFile parses a command .md file
func (m *DefaultSkillManager) parseCommandFile(path string, scope SkillScope) (*Skill, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	result, err := ParseSkillFile(string(content))
	if err != nil {
		return nil, err
	}

	// Get name from filename if not in frontmatter
	name := strings.TrimSuffix(filepath.Base(path), ".md")
	description := ""
	var allowedTools []string

	if result.Frontmatter != nil {
		if result.Frontmatter.Name != "" {
			name = result.Frontmatter.Name
		}
		description = result.Frontmatter.Description
		allowedTools = ParseAllowedTools(result.Frontmatter.AllowedTools)
	}

	return &Skill{
		Name:         name,
		Description:  description,
		AllowedTools: allowedTools,
		Type:         SkillTypeCommand,
		Scope:        scope,
		Path:         path,
	}, nil
}

// parseSkillDir parses a skill directory with SKILL.md
func (m *DefaultSkillManager) parseSkillDir(path string, scope SkillScope) (*Skill, error) {
	skillMD := filepath.Join(path, "SKILL.md")
	content, err := os.ReadFile(skillMD)
	if err != nil {
		return nil, err
	}

	result, err := ParseSkillFile(string(content))
	if err != nil {
		return nil, err
	}

	// Get name from directory if not in frontmatter
	name := filepath.Base(path)
	description := ""
	var allowedTools []string

	if result.Frontmatter != nil {
		if result.Frontmatter.Name != "" {
			name = result.Frontmatter.Name
		}
		description = result.Frontmatter.Description
		allowedTools = ParseAllowedTools(result.Frontmatter.AllowedTools)
	}

	return &Skill{
		Name:         name,
		Description:  description,
		AllowedTools: allowedTools,
		Type:         SkillTypeSkill,
		Scope:        scope,
		Path:         path,
	}, nil
}

// installFromLocal installs a skill from a local path
func (m *DefaultSkillManager) installFromLocal(sourcePath string, scope SkillScope) (*Skill, error) {
	info, err := os.Stat(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("source not found: %s", sourcePath)
	}

	if info.IsDir() {
		// It's a skill directory
		return m.installSkillDir(sourcePath, scope)
	}

	// It's a command file
	if strings.HasSuffix(sourcePath, ".md") {
		return m.installCommandFile(sourcePath, scope)
	}

	return nil, fmt.Errorf("unsupported file type: %s", sourcePath)
}

// installFromGit installs a skill from a git repository
// fragment is the skill name after # (e.g., repo#skill-name)
// skillPath is the path within repo (e.g., from tree URLs like /skills/frontend-design)
func (m *DefaultSkillManager) installFromGit(repoURL, fragment, skillPath string, scope SkillScope) (*Skill, error) {
	// Clone the repository
	tmpDir, err := GitClone(repoURL)
	if err != nil {
		return nil, err
	}
	defer CleanupTempDir(tmpDir)

	// Determine the skill location
	var targetPath string
	var lookupName string

	if skillPath != "" {
		// For tree URLs, use the full path within the repo
		targetPath = filepath.Join(tmpDir, skillPath)
		// Extract skill name from the path (last component)
		lookupName = filepath.Base(skillPath)
	} else {
		lookupName = fragment
	}

	// If we have a direct path from tree URL, check it first
	if targetPath != "" {
		if isSkillDir(targetPath) {
			skill, err := m.installSkillDir(targetPath, scope)
			if err != nil {
				return nil, err
			}
			skill.Source = repoURL
			if skillPath != "" {
				skill.Source = repoURL + "/tree/main/" + skillPath
			}
			return skill, nil
		}
		// Check if it's a command file
		if isCommandFile(targetPath + ".md") {
			skill, err := m.installCommandFile(targetPath+".md", scope)
			if err != nil {
				return nil, err
			}
			skill.Source = repoURL
			return skill, nil
		}
		return nil, fmt.Errorf("skill not found at path: %s", skillPath)
	}

	// Find the skill in the repository using fragment
	foundPath, err := FindSkillInRepo(tmpDir, lookupName)
	if err != nil {
		// Maybe it's a command file
		cmdPath, cmdErr := FindCommandInRepo(tmpDir, lookupName)
		if cmdErr != nil {
			return nil, err // Return original error
		}
		skill, err := m.installCommandFile(cmdPath, scope)
		if err != nil {
			return nil, err
		}
		skill.Source = repoURL
		if fragment != "" {
			skill.Source = repoURL + "#" + fragment
		}
		return skill, nil
	}

	skill, err := m.installSkillDir(foundPath, scope)
	if err != nil {
		return nil, err
	}
	skill.Source = repoURL
	if fragment != "" {
		skill.Source = repoURL + "#" + fragment
	}
	return skill, nil
}

// installSkillDir installs a skill directory
func (m *DefaultSkillManager) installSkillDir(sourcePath string, scope SkillScope) (*Skill, error) {
	// Parse the skill first to get its name
	skill, err := m.parseSkillDir(sourcePath, scope)
	if err != nil {
		return nil, err
	}

	// Get target directory
	skillsDir, err := GetSkillsDir(scope)
	if err != nil {
		return nil, err
	}

	// Ensure skills directory exists
	if err := EnsureDir(skillsDir); err != nil {
		return nil, err
	}

	// Copy the skill directory
	targetPath := filepath.Join(skillsDir, skill.Name)

	// Check if already exists
	if _, err := os.Stat(targetPath); err == nil {
		return nil, fmt.Errorf("skill already exists: %s", skill.Name)
	}

	if err := copyDir(sourcePath, targetPath); err != nil {
		return nil, fmt.Errorf("failed to copy skill: %w", err)
	}

	skill.Path = targetPath
	return skill, nil
}

// installCommandFile installs a command .md file
func (m *DefaultSkillManager) installCommandFile(sourcePath string, scope SkillScope) (*Skill, error) {
	// Parse the command first to get its name
	skill, err := m.parseCommandFile(sourcePath, scope)
	if err != nil {
		return nil, err
	}

	// Get target directory
	commandsDir, err := GetCommandsDir(scope)
	if err != nil {
		return nil, err
	}

	// Ensure commands directory exists
	if err := EnsureDir(commandsDir); err != nil {
		return nil, err
	}

	// Copy the command file
	targetPath := filepath.Join(commandsDir, skill.Name+".md")

	// Check if already exists
	if _, err := os.Stat(targetPath); err == nil {
		return nil, fmt.Errorf("command already exists: %s", skill.Name)
	}

	if err := copyFile(sourcePath, targetPath); err != nil {
		return nil, fmt.Errorf("failed to copy command: %w", err)
	}

	skill.Path = targetPath
	return skill, nil
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
