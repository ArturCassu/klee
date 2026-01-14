/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

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

		fmt.Println("Starting kafka container...")

		compose := exec.Command("docker-compose", "-f", "docker/kafka.yaml", "up", "-d")
		output, err := compose.CombinedOutput()
		if err != nil {
			fmt.Printf("Failed to start kafka container: %v\nOutput: %s\n", err, string(output))
			return
		}

		fmt.Println("Kafka container started successfully.")
	},
}

var kafkaStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Kafka and Kafka-UI containers",
	Long:  `Stop Kafka and Kafka-UI containers using Docker.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("kafka stop called")
	},
}

var kafkaConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure Kafka settings",
	Long:  `Configure Kafka settings for the current user.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("kafka config called")
	},
}

func init() {
	rootCmd.AddCommand(kafkaCmd)

	kafkaCmd.AddCommand(kafkaStartCmd)
	kafkaCmd.AddCommand(kafkaStopCmd)
	kafkaCmd.AddCommand(kafkaConfigCmd)
}
