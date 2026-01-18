package skills

// SkillType distinguishes between slash commands and full skills
type SkillType string

const (
	// SkillTypeCommand represents a single .md file slash command
	SkillTypeCommand SkillType = "command"
	// SkillTypeSkill represents a directory with SKILL.md
	SkillTypeSkill SkillType = "skill"
)

// SkillScope indicates where the skill is installed
type SkillScope string

const (
	// ScopePersonal is for skills in ~/.claude/
	ScopePersonal SkillScope = "personal"
	// ScopeProject is for skills in .claude/
	ScopeProject SkillScope = "project"
)

// Skill represents a Claude Code or Codex skill or command
type Skill struct {
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	AllowedTools []string   `json:"allowed_tools,omitempty"`
	Type         SkillType  `json:"type"`
	Scope        SkillScope `json:"scope"`
	Path         string     `json:"path"`
	Source       string     `json:"source,omitempty"`
}

// SkillStatus represents the health status of a skill
type SkillStatus struct {
	Skill  Skill
	Valid  bool
	Error  error
	Issues []string
}

// SkillManager handles skill operations
type SkillManager interface {
	// List returns all installed skills
	List() ([]Skill, error)

	// ListByScope returns skills filtered by scope
	ListByScope(scope SkillScope) ([]Skill, error)

	// Install installs a skill from a source
	Install(source string, scope SkillScope) (*Skill, error)

	// Remove removes a skill by name
	Remove(name string, scope SkillScope) error

	// Check verifies skills installation status
	Check() ([]SkillStatus, error)

	// Get retrieves a skill by name
	Get(name string) (*Skill, error)
}
