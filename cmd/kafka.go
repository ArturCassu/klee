/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ArturCassu/klee/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const defaultKafkaComposePath = "docker/kafka.yaml"

func kafkaComposeFilePath() string {
	if p := viper.GetString("kafka.composeFile"); p != "" {
		return p
	}
	return defaultKafkaComposePath
}

// kafkaCmd represents the kafka command
var kafkaCmd = &cobra.Command{
	Use:   "kafka",
	Short: "Start Kafka and Kafka-UI containers using Docker",
	Long: `Starts Kafka and Kafka-UI containers using Docker.
Ensures Rancher Desktop and Docker are running before starting containers.`,
}

var kafkaStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Kafka and Kafka-UI containers",
	Long:  `Start Kafka and Kafka-UI containers using Docker.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Start Docker
		if !utils.IsDockerRunning() {
			err := utils.CheckAndStartDocker()
			if err != nil {
				panic(err)
			}
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Starting kafka container...")

		composeFile := kafkaComposeFilePath()

		composeArgs := []string{"-p", "kafka", "-f", composeFile, "up", "-d"}
		fmt.Fprintf(cmd.OutOrStdout(), "Running: docker-compose %s\n", utils.ShellQuoteArgs(composeArgs...))

		compose := exec.Command("docker-compose", composeArgs...)

		// Stream output directly to the user so they see progress.
		compose.Stdout = cmd.OutOrStdout()
		compose.Stderr = cmd.ErrOrStderr()

		// Since Stdout/Stderr are set, don't use CombinedOutput below.
		// Run will return an error if docker-compose fails.
		err := compose.Run()
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Failed to start kafka container: %v\n", err)
			return
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Kafka container started successfully.")
	},
}

var kafkaStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Kafka and Kafka-UI containers",
	Long:  `Stop Kafka and Kafka-UI containers using Docker.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If docker isn't running, there's nothing meaningful we can do.
		// (Match start behavior by ensuring Docker is up so we can cleanly stop.)
		if !utils.IsDockerRunning() {
			err := utils.CheckAndStartDocker()
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Failed to start Docker: %v\n", err)
				return
			}
		}

		removeVolumes, _ := cmd.Flags().GetBool("volumes")

		fmt.Fprintln(cmd.OutOrStdout(), "Stopping kafka containers...")

		composeFile := kafkaComposeFilePath()

		composeArgs := []string{"-p", "kafka", "-f", composeFile, "down"}
		if removeVolumes {
			composeArgs = append(composeArgs, "-v")
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Running: docker-compose %s\n", utils.ShellQuoteArgs(composeArgs...))

		compose := exec.Command("docker-compose", composeArgs...)
		compose.Stdout = cmd.OutOrStdout()
		compose.Stderr = cmd.ErrOrStderr()

		err := compose.Run()
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Failed to stop kafka containers: %v\n", err)
			return
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Kafka containers stopped successfully.")
	},
}

var kafkaConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure Kafka settings",
	Long:  `Configure Kafka settings for the current user.`,
	Run: func(cmd *cobra.Command, args []string) {
		current := kafkaComposeFilePath()
		input := utils.PromptInput("Path to Kafka docker-compose file:", current)
		input = filepath.Clean(input)

		// If user cleared the input, keep the previously configured/default value.
		if input == "." || input == "" {
			fmt.Fprintf(cmd.OutOrStdout(), "Keeping kafka compose file: %s\n", current)
			return
		}

		if _, err := os.Stat(input); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Invalid path: %s (%v)\n", input, err)
			return
		}

		if err := utils.PersistSetting("kafka.composeFile", input); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Failed to save config: %v\n", err)
			return
		}

		fmt.Fprintf(cmd.OutOrStdout(), "✓ Saved kafka compose file: %s\n", input)
	},
}

func init() {
	rootCmd.AddCommand(kafkaCmd)

	kafkaCmd.AddCommand(kafkaStartCmd)
	kafkaCmd.AddCommand(kafkaStopCmd)
	kafkaCmd.AddCommand(kafkaConfigCmd)

	kafkaStopCmd.Flags().BoolP("volumes", "v", false, "Remove named volumes declared in the compose file")
}
