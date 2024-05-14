package cmd

import (
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run <method>",
	Short: "run command is used to run the Go project.",
	Long:  `run command is used to run the local Go project using docker-compose.`,
	Args:  cobra.ExactArgs(1), // Ensures exactly one argument is passed
	Run: func(cmd *cobra.Command, args []string) {
		//! Check if either docker-compose.yaml or docker-compose.yml file exists
		if _, err := os.Stat("docker-compose.yaml"); os.IsNotExist(err) {
			if _, err := os.Stat("docker-compose.yml"); os.IsNotExist(err) {
				log.Fatalf("docker-compose file not found")
				return
			}
		}

		// Set the environment based on the method
		method := args[0]
		var env string
		switch method {
		case "":
			env = "development"
		case "dev":
			env = "development"
		case "development":
			env = "development"
		case "prod":
			env = "production"
		case "production":
			env = "production"
		default:
			log.Fatalf("Unsupported method: %s", method)
			return
		}

		// Construct the Docker Compose command
		command := exec.Command("sh", "-c", "docker-compose up --build")

		// Set the environment variables and PATH if necessary
		command.Env = append(os.Environ(), "BUILD_ENV="+env)

		// Run the command and capture output
		command.Stdout = cmd.OutOrStdout() // Redirect output to Cobra command's stdout
		command.Stderr = cmd.OutOrStderr() // Redirect errors to Cobra command's stderr

		// Execute the command
		if err := command.Run(); err != nil {
			log.Fatalf("Failed to run docker-compose: %v", err)
		}
	},
}
