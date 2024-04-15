package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init <project-name> [--path <path>]",
	Short: "Init command is used to initialize a new Go project.",
	Long:  `Init command is used to initialize a new Go project. It creates a new directory with the project name and initializes a new Go module in it.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Project name is required")
		}
		name := args[0] // First positional argument is the project name
		path, _ := cmd.Flags().GetString("path")
		if path == "" {
			path = "." // Default path is the current directory
		}
		setUpProject(name, path)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Local flag for the path
	initCmd.Flags().StringP("path", "p", "", "Path to create the project")
}

func setUpProject(name, path string) {
	if path != "." {
		// Create the directory if it does not exist
		fullPath := filepath.Join(path, name)
		err := os.MkdirAll(fullPath, os.ModePerm)
		if err != nil {
			log.Fatalf("Error creating directory: %v", err)
		}
		os.Chdir(fullPath) // Change the working directory to the new project directory
	}

	// Create a new Go module
	createGoModule(name)

	// Initialize a new Go module
	cmd := fmt.Sprintf("go mod init %s", name)
	// Simulate command execution (you need to execute this in an actual shell or use exec.Command)
	fmt.Printf("Running command: %s\n", cmd)

	// Run the command

}

func createGoModule(name string) {
	// Create a new Go module
	fmt.Printf("Creating a new Go module: %s\n", name)

	// run go mod init <name>
	cmd := exec.Command("go", "mod", "init", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error creating Go module: %v", err)
	}
}
