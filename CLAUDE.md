# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Run Commands

```bash
# Build the binary (injects version from git tags)
make build

# Run the TUI
./agentx

# Run CLI commands
./agentx install playwright
./agentx skills list
./agentx plugins list

# Run tests
go test ./...

# Run a single test
go test -run TestParseGitHubTreeURL ./internal/skills/

# Clean build artifacts
make clean
```

## Git Operations

**Never run `git commit` or `git push`** - the user handles all git operations manually.

## Architecture

AgentX is a CLI tool for managing MCP servers, skills, and plugins across multiple AI coding agents (Claude Code, Codex, Cursor, Gemini CLI, OpenCode).

### Core Layers

**CLI Layer** (`cmd/`): Cobra-based commands. `root.go` launches the TUI when run without arguments, or delegates to subcommands (`install`, `check`, `list`, `remove`, `skills`, `plugins`).

**Agent Abstraction** (`internal/agent/`): The `Agent` interface defines operations for all supported agents. Each agent (Claude, Codex, Cursor, Gemini, OpenCode) implements this interface with its own config file location and format:
- Claude Code: `~/.claude.json`
- Codex: `~/.codex/config.toml`
- Cursor: `~/.cursor/mcp.json`
- Gemini CLI: `~/.gemini/settings.json`
- OpenCode: `~/.opencode/config.json`

**Domain Logic** (`internal/`):
- `skills/`: Skill installation from local paths or Git repos, YAML frontmatter parsing
- `plugins/`: Plugin management with remote registry fetch and local caching
- `config/`: JSON config file reading/writing for agent configurations
- `mcp/`: MCP server definitions (Playwright, Context7)

**TUI Layer** (`ui/`): Built with Bubble Tea and Lip Gloss.
- `app.go`: Main model coordinating views, tabs, header, sidebar, and footer
- `views/`: Four tab views (MCP, Skills, Plugins, Agents) - each implements the View interface
- `components/`: Reusable UI pieces (TabBar, Header, Sidebar, Footer)
- `theme/`: Color scheme and styling

### Data Flow

1. TUI initializes all views, each view queries agents for current state
2. User actions in views call agent methods (e.g., `InstallPlaywright()`)
3. Agent methods modify JSON config files in user's home directory
4. Views refresh by re-querying agent state

### Registry

`registry/plugins.json` contains available plugins. The TUI fetches from the remote URL (`https://raw.githubusercontent.com/agentsdance/agentx/master/registry/plugins.json`) with local cache fallback.
