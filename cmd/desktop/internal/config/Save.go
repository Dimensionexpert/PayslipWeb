package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func Save(config Config) error {
	path, err := ConfigPath()
	if err != nil {
		return fmt.Errorf("getting config path: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return fmt.Errorf("encoding config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}
