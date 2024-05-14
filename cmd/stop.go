package cmd

import (
	"log"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(downCmd)
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the Docker Compose services.",
	Long:  `Stop the Docker Compose services and remove orphans.`,
	Run: func(cmd *cobra.Command, args []string) {
		executeDockerComposeDown(cmd)
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop the Docker Compose services.",
	Long:  `Stop the Docker Compose services and remove orphans.`,
	Run: func(cmd *cobra.Command, args []string) {
		executeDockerComposeDown(cmd)
	},
}

func executeDockerComposeDown(cmd *cobra.Command) {
	// Construct the Docker Compose command
	command := exec.Command("sh", "-c", "docker-compose down --remove-orphans")

	// Redirect output and errors to Cobra command's stdout and stderr
	command.Stdout = cmd.OutOrStdout()
	command.Stderr = cmd.OutOrStderr()

	// Execute the command
	if err := command.Run(); err != nil {
		log.Fatalf("Failed to run docker-compose: %v", err)
	}
}
