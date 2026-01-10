package plugins

// PluginAuthor represents the author of a plugin
type PluginAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
	URL   string `json:"url,omitempty"`
}

// PluginManifest represents the plugin.json manifest file
type PluginManifest struct {
	Name        string       `json:"name"`
	Version     string       `json:"version"`
	Description string       `json:"description"`
	Author      PluginAuthor `json:"author"`
}

// PluginComponents represents what a plugin provides
type PluginComponents struct {
	Commands   []string `json:"commands,omitempty"`
	Agents     []string `json:"agents,omitempty"`
	Skills     []string `json:"skills,omitempty"`
	Hooks      []string `json:"hooks,omitempty"`
	MCPServers []string `json:"mcp_servers,omitempty"`
}

// Plugin represents an installed plugin
type Plugin struct {
	Name        string           `json:"name"`
	Version     string           `json:"version"`
	Description string           `json:"description"`
	Author      PluginAuthor     `json:"author"`
	Path        string           `json:"path"`
	Source      string           `json:"source,omitempty"`
	Components  PluginComponents `json:"components"`
}

// PluginStatus represents the health status of a plugin
type PluginStatus struct {
	Plugin Plugin
	Valid  bool
	Error  error
	Issues []string
}

// PluginManager handles plugin operations
type PluginManager interface {
	// List returns all installed plugins
	List() ([]Plugin, error)

	// Install installs a plugin from a source
	Install(source string) (*Plugin, error)

	// Remove removes a plugin by name
	Remove(name string) error

	// Check verifies plugins installation status
	Check() ([]PluginStatus, error)

	// Get retrieves a plugin by name
	Get(name string) (*Plugin, error)
}
