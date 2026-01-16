package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// EnsureConfigReady ensures Viper has a config file set and that the directory/file exist.
// It is safe to call multiple times.
func EnsureConfigReady() error {
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("error getting home directory: %w", err)
		}
		configFile = filepath.Join(home, ".config", "klee", "config.yaml")
		viper.SetConfigFile(configFile)
		viper.SetConfigType("yaml")
	}

	configDir := filepath.Dir(configFile)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := os.WriteFile(configFile, []byte("# Klee Configuration\n"), 0644); err != nil {
			return fmt.Errorf("error creating config file: %w", err)
		}
	}

	return nil
}

// PersistSetting updates a viper key and writes it to the configured config file.
func PersistSetting(key string, value any) error {
	if err := EnsureConfigReady(); err != nil {
		return err
	}

	viper.Set(key, value)

	// If the file doesn't exist or wasn't read successfully yet, WriteConfig can fail.
	// WriteConfigAs is safer because it will create the file.
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		return fmt.Errorf("no config file configured")
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return viper.WriteConfigAs(configFile)
	}

	if err := viper.WriteConfig(); err != nil {
		// Fallback in case the config file wasn't read yet.
		return viper.WriteConfigAs(configFile)
	}

	return nil
}
