package cmd

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/template"

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
	Module string
	Path   string
	// mandatory modules
	IncludeLogger     bool
	IncludeConfig     bool
	IncludeHTTPServer bool
	// optional modules
	IncludeJWTMiddleware bool
	IncludeDBConnector   bool
	IncludeMailer        bool

	SetUpGit        bool
	SetUpDockerFile bool

	// todo add project specific modules
}

func setUpProject() {
	// greeting message
	if !selectYesNo("Welcome to the GoGetter CLI. This will begin the setup process for your new Go service. Continue?", true) {
		fmt.Println("Project setup cancelled.")
		return
	}

	//* Collect Project Configurations
	projectConfig := &ProjectConfig{
		IncludeLogger:     true,
		IncludeConfig:     true,
		IncludeHTTPServer: true,
	}

	// project name
	projectConfig.Module = inputText("Enter the go module name for your project. [Example: github.com/alsey89/gogetter]", true)

	// project directory
	projectConfig.Path = inputText("Enter the path for your project. Service will be initiated at the current directory if left empty.", false)

	// project modules
	projectConfig.IncludeJWTMiddleware = selectYesNo("Do you want to include a JWT middleware module?", true)
	projectConfig.IncludeDBConnector = selectYesNo("Do you want to include a database module with Postgres and GORM?", true)
	projectConfig.IncludeMailer = selectYesNo("Do you want to include a mailer module using Gomail?", true)

	// git setup
	projectConfig.SetUpGit = selectYesNo("Do you want to set up git for the project?", true)
	// docker setup
	projectConfig.SetUpDockerFile = selectYesNo("Do you want to set up docker for the project?", true)

	executeSetup(projectConfig)
}

func executeSetup(config *ProjectConfig) {
	var err error

	err = createGoModule(config.Module)
	if err != nil {
		log.Fatalf("Error creating Go module: %v", err)
	}

	err = createMainFile(config)
	if err != nil {
		log.Fatalf("Error creating main file: %v", err)
	}

	if config.SetUpDockerFile {
		err = createDockerfile(false)
		if err != nil {
			log.Fatalf("Error creating Dockerfile: %v", err)
		}
	}

	if config.SetUpGit {
		err = createGitRepository()
		if err != nil {
			log.Fatalf("Error creating Git repository: %v", err)
		}
	}
}

func createGoModule(name string) error {
	// Create a new Go module
	log.Printf("Creating a new Go module: %s\n", name)

	// run go mod init <name>
	cmd := exec.Command("go", "mod", "init", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

//go:embed templates/*
var templateFS embed.FS

func createMainFile(projectConfig *ProjectConfig) error {
	// Using the embedded file system to access the template
	tmpl, err := template.ParseFS(templateFS, "templates/main.go.tpl")
	if err != nil {
		log.Printf("Failed to parse template: %v", err)
		return err
	}

	file, err := os.Create("main.go")
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		return err
	}
	defer file.Close()

	// Execute the template with the project configuration
	if err := tmpl.Execute(file, projectConfig); err != nil {
		log.Printf("Failed to execute template: %v", err)
		return err
	}

	// run go mod init <name>
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func createDockerfile(forDevelopment bool) error {
	var templatePath string
	var dockerfileName string

	if forDevelopment {
		templatePath = "templates/Dockerfile.dev.tpl"
		dockerfileName = "Dockerfile.dev"
	} else {
		templatePath = "templates/Dockerfile.tpl"
		dockerfileName = "Dockerfile"
	}

	// Using the embedded file system to access the template
	tmpl, err := template.ParseFS(templateFS, templatePath)
	if err != nil {
		log.Printf("Failed to parse template: %v", err)
		return err
	}

	file, err := os.Create(dockerfileName)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		return err
	}
	defer file.Close()

	// Execute the template with the project configuration
	if err := tmpl.Execute(file, nil); err != nil {
		log.Printf("Failed to execute template: %v", err)
		return err
	}

	return nil
}

func createGitRepository() error {
	// Create a new Git repository
	log.Printf("Creating a new Git repository\n")

	// run git init
	cmd := exec.Command("git", "init")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	// create .gitignore file
	file, err := os.Create(".gitignore")
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		return err
	}
	defer file.Close()

	tmpl, err := template.ParseFS(templateFS, "templates/gitignore.tpl")
	if err != nil {
		log.Printf("Failed to parse template: %v", err)
		return err
	}

	if err := tmpl.Execute(file, nil); err != nil {
		log.Printf("Failed to execute template: %v", err)
		return err
	}

	return nil
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
