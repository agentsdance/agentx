package version

// These variables are set at build time using ldflags
var (
	// Version is the current version of AgentX (set via -ldflags or defaults to "dev")
	Version = "dev"
	// GitCommit is the git commit hash
	GitCommit = "unknown"
	// BuildDate is the build timestamp
	BuildDate = "unknown"
)
