package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init <project-name> [--path <path>]",
	Short: "Init command is used to initialize a new Go project.",
	Long:  `Init command is used to initialize a new Go project. It creates a new directory with the project name and initializes a new Go module in it.`,
	Run: func(cmd *cobra.Command, args []string) {
		// check if go mod already exists in the directory
		_, err := os.Stat("go.mod")
		if err == nil {
			log.Fatalf("go.mod already exists in the directory")
		}

		setUpProject()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

// ----------------------------------------------------------------------------

type ProjectConfig struct {
	Module      string
	Path        string
	Modules     []string
	SetUpGit    bool
	SetUpDocker bool
}

func setUpProject() {
	var yesNoSelection bool

	// greeting message
	yesNoSelection = selectYesNo("Welcome to the GoGetter CLI. This will begin the setup process for your new Go service. Continue?", true)
	if !yesNoSelection {
		fmt.Println("Project setup cancelled.")
		return
	}

	//* Collect Project Configurations
	projectConfig := &ProjectConfig{}

	// project name
	moduleName := inputText("Enter the go module name for your project. [Example: github.com/alsey89/gogetter]", true)
	projectConfig.Module = moduleName

	// project directory
	path := inputText("Enter the path for your project. Service will be initiated at the current directory if left empty.", true)
	projectConfig.Path = path

	// project modules
	// todo - add more modules and different options for each module
	selectedModules := selectMultipleOptions("Select the modules you want to include in your project:", []string{
		"HTTPServer: echo",
		"Database:   postgres + gorm",
		"Logger:     zap",
		"Config:     viper",
		"Auth: 	 	 jwt",
		"Mailer:     gomail",
	})
	if len(selectedModules) == 0 {
		yesNoSelection = selectYesNo("No modules selected. Do you want to abort setup? An empty go project will be initiated if you select 'No'", true)
		if yesNoSelection {
			fmt.Println("Project setup cancelled.")
			return
		}
	}
	projectConfig.Modules = selectedModules

	// git setup
	yesNoSelection = selectYesNo("Do you want to set up git for the project?", true)
	projectConfig.SetUpGit = yesNoSelection

	// docker setup
	yesNoSelection = selectYesNo("Do you want to set up docker for the project?", true)
	projectConfig.SetUpDocker = yesNoSelection

	// Step 2: Execute setup based on collected configurations
	executeSetup(projectConfig)
}

func executeSetup(config *ProjectConfig) {
	// Step 1: Create a new directory with the project name if path is not empty or "."
	if config.Path != "" && config.Path != "." {
		err := os.MkdirAll(config.Path, os.ModePerm)
		if err != nil {
			log.Fatalf("Error creating project directory: %v", err)
		}
	}

	// Step 2: Change directory to the project directory
	os.Chdir(config.Path)

	// Step 3: Create a new Go module
	createGoModule(config.Module)

	// Step 4: Create a new main.go file
	// createMainFile(config.Modules)

	// Step 5: Initialize git
	// if config.SetUpGit {
	// 	initializeGit()
	// }

	// Step 6: Initialize docker
	// if config.SetUpDocker {
	// 	initializeDocker()
	// }
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

func createMainFile(modules []string) {
	// Create a new main.go file
	fmt.Println("Creating a new main.go file")

	// Create a new main.go file
	f, err := os.Create("main.go")
	if err != nil {
		log.Fatalf("Error creating main.go file: %v", err)
	}

	// Write the content to the main.go file
	f.Write([]byte(`package main`))

	// Add the selected modules to the main.go file
	for _, module := range modules {
		f.Write([]byte(fmt.Sprintf("\n\n// Module: %s", module)))
	}
}

// ----------------------------------------------------------------------------

func selectMultipleOptions(message string, options []string) []string {
	var selectedOptions []string

	// set up multi select prompt
	prompt := &survey.MultiSelect{
		Message: message,
		Options: options,
		Help:    "Use the arrow keys to navigate and space to select. Press enter when done.",
	}

	// run the prompt, save result to selectedOptions
	err := survey.AskOne(prompt, &selectedOptions)
	if err != nil {
		log.Fatalf("Failed to select options: %v", err)
	}

	return selectedOptions
}

func selectOneOption(message string, options []string) string {
	var selectedOption string

	// set up select prompt
	prompt := &survey.Select{
		Message: message,
		Options: options,
		Help:    "Use the arrow keys to navigate and space to select. Press enter when done.",
	}

	// run the prompt, save result to selectedOption
	err := survey.AskOne(prompt, &selectedOption)
	if err != nil {
		log.Fatalf("Failed to select option: %v", err)
	}

	return selectedOption
}

func selectYesNo(message string, defaultSelection bool) bool {
	var selectedOption bool

	// set up confirm prompt
	prompt := &survey.Confirm{
		Message: message,
		Default: defaultSelection,
	}

	// run the prompt, save result to selectedOption
	err := survey.AskOne(prompt, &selectedOption)
	if err != nil {
		log.Fatalf("Failed to select option: %v", err)
	}

	return selectedOption
}

func inputText(message string, isMandatory bool) string {
	var input string
	var opts []survey.AskOpt

	prompt := &survey.Input{
		Message: message,
	}

	if isMandatory {
		// this makes the input a required field
		opts = append(opts, survey.WithValidator(survey.Required))
	}

	// run the prompt, save result to input
	err := survey.AskOne(prompt, &input, opts...)
	if err != nil {
		log.Fatalf("Failed to input text: %v", err)
	}

	return input
}
