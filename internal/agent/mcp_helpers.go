package agent

import "fmt"

// MCPConfigEntry represents a discovered MCP configuration.
type MCPConfigEntry struct {
	Config map[string]interface{}
	Source string
}

// CollectMCPConfigs collects MCP configs from all agents, keyed by server name.
func CollectMCPConfigs(agents []Agent) map[string]MCPConfigEntry {
	configs := make(map[string]MCPConfigEntry)
	for _, a := range agents {
		entries, err := a.ListMCPs()
		if err != nil {
			continue
		}
		for name, cfg := range entries {
			if _, exists := configs[name]; exists {
				continue
			}
			if cfg == nil {
				continue
			}
			configs[name] = MCPConfigEntry{
				Config: cloneMCPConfig(cfg),
				Source: a.Name(),
			}
		}
	}
	return configs
}

func cloneMCPConfig(cfg map[string]interface{}) map[string]interface{} {
	if cfg == nil {
		return nil
	}
	cloned := make(map[string]interface{}, len(cfg))
	for key, value := range cfg {
		switch v := value.(type) {
		case map[string]interface{}:
			cloned[key] = cloneMCPConfig(v)
		case []interface{}:
			items := make([]interface{}, len(v))
			copy(items, v)
			cloned[key] = items
		case []string:
			items := make([]string, len(v))
			copy(items, v)
			cloned[key] = items
		default:
			cloned[key] = value
		}
	}
	return cloned
}

func normalizeArgsToStrings(cfg map[string]interface{}) map[string]interface{} {
	cloned := cloneMCPConfig(cfg)
	if cloned == nil {
		return nil
	}
	args, ok := cloned["args"]
	if !ok {
		return cloned
	}
	switch v := args.(type) {
	case []string:
		return cloned
	case []interface{}:
		values := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				values = append(values, s)
			} else {
				values = append(values, fmt.Sprint(item))
			}
		}
		cloned["args"] = values
	}
	return cloned
}
