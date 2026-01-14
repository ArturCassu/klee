/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "klee",
	Short: "",
	Long:  ``,
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/klee/config.yaml)")
	ensureConfigFile()
}

func initConfig() {
	// Config file reading logic
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Use default location: ~/.config/klee/config.yaml
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		configFile := home + "/.config/klee/config.yaml"
		viper.SetConfigFile(configFile)
	}

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	// Read config file (it's okay if it doesn't exist yet)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file was found but another error was produced
			fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
		}
	}
}

// ensureConfigFile creates the config directory and file if they don't exist
func ensureConfigFile() error {
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("error getting home directory: %w", err)
		}
		configFile = home + "/.config/klee/config.yaml"
	}

	// Create directory
	configDir := configFile[:len(configFile)-len("/config.yaml")]
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	// Create file if it doesn't exist
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := os.WriteFile(configFile, []byte("# Klee Configuration\n"), 0644); err != nil {
			return fmt.Errorf("error creating config file: %w", err)
		}
	}

	return nil
}
