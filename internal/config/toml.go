package config

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// ReadTOMLConfig reads a TOML config file.
func ReadTOMLConfig(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if len(bytes.TrimSpace(data)) == 0 {
		return map[string]interface{}{}, nil
	}

	var cfg map[string]interface{}
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg == nil {
		cfg = map[string]interface{}{}
	}
	return cfg, nil
}

// WriteTOMLConfig writes a TOML config file with pretty formatting.
func WriteTOMLConfig(path string, cfg map[string]interface{}) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	if cfg == nil {
		cfg = map[string]interface{}{}
	}

	data, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
