package config

import (
	"os"
	"path/filepath"
)

func ConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "payslip", "config.json"), nil

}
