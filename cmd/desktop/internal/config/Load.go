package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func Load() (Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return Config{}, fmt.Errorf("getting config path: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, nil
		}

		return Config{}, fmt.Errorf("reading config file: %w", err)
	}

	var config Config

	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("decoding config file: %w", err)
	}

	return config, nil
}
