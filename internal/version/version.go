package version

import "fmt"

// These variables are set at build time using ldflags
var (
	// Version is the current version of AgentX (set via -ldflags or defaults to "dev")
	Version = "dev"
	// GitCommit is the git commit hash
	GitCommit = "unknown"
	// BuildDate is the build timestamp
	BuildDate = "unknown"
)

// GetFullVersion returns formatted version info with commit and build date
func GetFullVersion() string {
	base := fmt.Sprintf("agentx version %s", Version)
	if GitCommit == "unknown" && BuildDate == "unknown" {
		return base
	}
	if GitCommit != "unknown" && BuildDate != "unknown" {
		return fmt.Sprintf("%s\nCommit: %s\nBuilt: %s", base, GitCommit, BuildDate)
	}
	if GitCommit != "unknown" {
		return fmt.Sprintf("%s\nCommit: %s", base, GitCommit)
	}
	if BuildDate != "unknown" {
		return fmt.Sprintf("%s\nBuilt: %s", base, BuildDate)
	}
	return base
}
