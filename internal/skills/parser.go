package skills

import (
	"bufio"
	"strings"

	"gopkg.in/yaml.v3"
)

// Frontmatter represents the YAML frontmatter in skill/command files
type Frontmatter struct {
	Name         string `yaml:"name"`
	Description  string `yaml:"description"`
	AllowedTools string `yaml:"allowed-tools"` // Comma-separated
}

// ParseResult contains the parsed frontmatter and body content
type ParseResult struct {
	Frontmatter *Frontmatter
	Body        string
}

// ParseSkillFile parses a skill/command markdown file
func ParseSkillFile(content string) (*ParseResult, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))

	// Check for frontmatter delimiter
	if !scanner.Scan() {
		return &ParseResult{Body: content}, nil
	}

	firstLine := strings.TrimSpace(scanner.Text())
	if firstLine != "---" {
		// No frontmatter, return content as-is
		return &ParseResult{Body: content}, nil
	}

	// Collect frontmatter lines
	var frontmatterLines []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "---" {
			break
		}
		frontmatterLines = append(frontmatterLines, line)
	}

	// Parse YAML frontmatter
	var fm Frontmatter
	frontmatterYAML := strings.Join(frontmatterLines, "\n")
	if err := yaml.Unmarshal([]byte(frontmatterYAML), &fm); err != nil {
		return nil, err
	}

	// Collect remaining content (body)
	var bodyLines []string
	for scanner.Scan() {
		bodyLines = append(bodyLines, scanner.Text())
	}

	return &ParseResult{
		Frontmatter: &fm,
		Body:        strings.Join(bodyLines, "\n"),
	}, nil
}

// ParseAllowedTools splits the comma-separated allowed-tools string
func ParseAllowedTools(tools string) []string {
	if tools == "" {
		return nil
	}

	parts := strings.Split(tools, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
