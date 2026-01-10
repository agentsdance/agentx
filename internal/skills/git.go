package skills

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// GitClone clones a repository to a temporary directory and returns the path
func GitClone(repoURL string) (string, error) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "agentx-skill-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Run git clone with --quiet to suppress progress output
	// This prevents messing up the TUI when running in interactive mode
	cmd := exec.Command("git", "clone", "--depth", "1", "--quiet", repoURL, tmpDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", fmt.Errorf("git clone failed: %w\n%s", err, string(output))
	}

	return tmpDir, nil
}

// FindSkillInRepo finds a skill directory in a cloned repository
// If skillName is empty, it looks for skills in standard locations
func FindSkillInRepo(repoPath, skillName string) (string, error) {
	// If a specific skill name is provided, look for it
	if skillName != "" {
		// Check in root
		skillPath := filepath.Join(repoPath, skillName)
		if isSkillDir(skillPath) {
			return skillPath, nil
		}

		// Check in skills/ subdirectory
		skillPath = filepath.Join(repoPath, "skills", skillName)
		if isSkillDir(skillPath) {
			return skillPath, nil
		}

		// Check in .claude/skills/ subdirectory
		skillPath = filepath.Join(repoPath, ".claude", "skills", skillName)
		if isSkillDir(skillPath) {
			return skillPath, nil
		}

		return "", fmt.Errorf("skill '%s' not found in repository", skillName)
	}

	// If no skill name, check if the repo root is a skill
	if isSkillDir(repoPath) {
		return repoPath, nil
	}

	// Check if there's a single skill in skills/
	skillsDir := filepath.Join(repoPath, "skills")
	if entries, err := os.ReadDir(skillsDir); err == nil && len(entries) == 1 {
		skillPath := filepath.Join(skillsDir, entries[0].Name())
		if isSkillDir(skillPath) {
			return skillPath, nil
		}
	}

	return "", fmt.Errorf("no skill found in repository root; use URL#skill-name to specify")
}

// FindCommandInRepo finds a command file in a cloned repository
func FindCommandInRepo(repoPath, commandName string) (string, error) {
	// If a specific command name is provided
	if commandName != "" {
		// Check in root
		cmdPath := filepath.Join(repoPath, commandName+".md")
		if isCommandFile(cmdPath) {
			return cmdPath, nil
		}

		// Check in commands/ subdirectory
		cmdPath = filepath.Join(repoPath, "commands", commandName+".md")
		if isCommandFile(cmdPath) {
			return cmdPath, nil
		}

		// Check in .claude/commands/ subdirectory
		cmdPath = filepath.Join(repoPath, ".claude", "commands", commandName+".md")
		if isCommandFile(cmdPath) {
			return cmdPath, nil
		}

		return "", fmt.Errorf("command '%s' not found in repository", commandName)
	}

	return "", fmt.Errorf("command name required for repository source")
}

// isSkillDir checks if a directory contains a SKILL.md file
func isSkillDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return false
	}

	skillMD := filepath.Join(path, "SKILL.md")
	_, err = os.Stat(skillMD)
	return err == nil
}

// isCommandFile checks if a file exists and has .md extension
func isCommandFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return false
	}
	return filepath.Ext(path) == ".md"
}

// CleanupTempDir removes a temporary directory
func CleanupTempDir(path string) {
	os.RemoveAll(path)
}
