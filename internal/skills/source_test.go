package skills

import (
	"testing"
)

func TestParseGitHubTreeURL(t *testing.T) {
	tests := []struct {
		name          string
		source        string
		wantRepoURL   string
		wantSkillPath string
		wantIsTreeURL bool
	}{
		{
			name:          "GitHub tree URL with skills path",
			source:        "https://github.com/anthropics/skills/tree/main/skills/frontend-design",
			wantRepoURL:   "https://github.com/anthropics/skills",
			wantSkillPath: "skills/frontend-design",
			wantIsTreeURL: true,
		},
		{
			name:          "GitHub tree URL with single path",
			source:        "https://github.com/user/repo/tree/master/my-skill",
			wantRepoURL:   "https://github.com/user/repo",
			wantSkillPath: "my-skill",
			wantIsTreeURL: true,
		},
		{
			name:          "GitHub tree URL without path",
			source:        "https://github.com/user/repo/tree/main",
			wantRepoURL:   "https://github.com/user/repo",
			wantSkillPath: "",
			wantIsTreeURL: true,
		},
		{
			name:          "Regular GitHub repo URL",
			source:        "https://github.com/user/repo",
			wantRepoURL:   "",
			wantSkillPath: "",
			wantIsTreeURL: false,
		},
		{
			name:          "Non-GitHub URL",
			source:        "https://gitlab.com/user/repo/tree/main/skills",
			wantRepoURL:   "",
			wantSkillPath: "",
			wantIsTreeURL: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRepoURL, gotSkillPath, gotIsTreeURL := parseGitHubTreeURL(tt.source)
			if gotRepoURL != tt.wantRepoURL {
				t.Errorf("parseGitHubTreeURL() repoURL = %v, want %v", gotRepoURL, tt.wantRepoURL)
			}
			if gotSkillPath != tt.wantSkillPath {
				t.Errorf("parseGitHubTreeURL() skillPath = %v, want %v", gotSkillPath, tt.wantSkillPath)
			}
			if gotIsTreeURL != tt.wantIsTreeURL {
				t.Errorf("parseGitHubTreeURL() isTreeURL = %v, want %v", gotIsTreeURL, tt.wantIsTreeURL)
			}
		})
	}
}

func TestParseSourceGitHubTreeURL(t *testing.T) {
	tests := []struct {
		name          string
		source        string
		wantType      SourceType
		wantRepoURL   string
		wantSkillPath string
		wantErr       bool
	}{
		{
			name:          "GitHub tree URL with skills path",
			source:        "https://github.com/anthropics/skills/tree/main/skills/frontend-design",
			wantType:      SourceTypeGitRepoWithFragment,
			wantRepoURL:   "https://github.com/anthropics/skills",
			wantSkillPath: "skills/frontend-design",
			wantErr:       false,
		},
		{
			name:        "Regular GitHub repo URL",
			source:      "https://github.com/user/repo",
			wantType:    SourceTypeGitRepo,
			wantRepoURL: "https://github.com/user/repo",
			wantErr:     false,
		},
		{
			name:        "GitHub repo URL with fragment",
			source:      "https://github.com/anthropics/skills#frontend-design",
			wantType:    SourceTypeGitRepoWithFragment,
			wantRepoURL: "https://github.com/anthropics/skills",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSource(tt.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if got.Type != tt.wantType {
				t.Errorf("ParseSource() Type = %v, want %v", got.Type, tt.wantType)
			}
			if got.RepoURL != tt.wantRepoURL {
				t.Errorf("ParseSource() RepoURL = %v, want %v", got.RepoURL, tt.wantRepoURL)
			}
			if tt.wantSkillPath != "" && got.SkillPath != tt.wantSkillPath {
				t.Errorf("ParseSource() SkillPath = %v, want %v", got.SkillPath, tt.wantSkillPath)
			}
		})
	}
}
