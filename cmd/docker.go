/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// dockerCmd represents the docker command
var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "A brief description",
	Long:  `A longer description`,
}

var dockerStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a Docker container",
	Long:  `Start a Docker container with the specified configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("docker start called")
	},
}

var dockerStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a Docker container",
	Long:  `Stop a running Docker container.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("docker stop called")
	},
}

var dockerConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure Docker settings",
	Long:  `Configure Docker settings for the current user.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("docker config called")
	},
}

func init() {
	rootCmd.AddCommand(dockerCmd)
	dockerCmd.AddCommand(dockerStartCmd)
	dockerCmd.AddCommand(dockerStopCmd)
	dockerCmd.AddCommand(dockerConfigCmd)
}
