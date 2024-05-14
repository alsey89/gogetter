package cmd

import (
	"embed"
	"log"
	"os"
	"os/exec"
	"text/template"

	"github.com/alsey89/gogetter/cmd/survey"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init", //"init <project-name> [--path <path>]"
	Short: "Init command is used to initialize a new Go project.",
	Long:  `Init command is used to initialize a new Go project. It initializes go mod, creates a main.go file, and optionally sets up optional modules, git, docker and docker-compose.`,
	Run: func(cmd *cobra.Command, args []string) {
		// check if go mod already exists in the directory
		_, err := os.Stat("go.mod")
		if err == nil {
			log.Fatalf("go.mod already exists in the directory")
		}

		setUpProject()
	},
}

func setUpProject() {
	config := gatherRequirements()
	executeSetup(config)
	saveConfig(config)
}

func gatherRequirements() *ProjectConfig {
	var err error
	var boolResult bool
	var stringResult *string

	c := &ProjectConfig{
		IncludeLogger: true,
		Logger:        "zap",

		IncludeConfig: true,
		ConfigManager: "viper",

		IncludeHTTPServer: true,
		Framework:         "echo",
	}

	//----- greeting message -----
	boolResult, err = survey.SelectYesNo("Welcome to the GoGetter CLI. This will begin the setup process for your new Go service. Continue?", true)
	if !boolResult {
		log.Fatal("Exiting setup process.")
	}
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	//----- project name -----
	stringResult, err = survey.InputText("Enter the go module name for your project. [Example: github.com/alsey89/gogetter]", true)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	if stringResult == nil {
		log.Fatalf("Project name is required.")
	}
	c.Module = *stringResult

	//----- project directory -----
	stringResult, err = survey.InputText("Enter the directory for your project. Service will be initiated at the current directory if left empty.", false)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	if stringResult == nil {
		c.Dir = ""
	} else {
		c.Dir = *stringResult
	}

	//----- project modules -----

	// JWT middleware
	boolResult, err = survey.SelectYesNo("Do you want to include a JWT middleware module?", true)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	c.IncludeJWTMiddleware = boolResult
	//todo: offer different JWT middleware options
	c.JWTMiddleware = "echo"

	// DB connector
	boolResult, err = survey.SelectYesNo("Do you want to include a GORM Postgres database connector module?", true)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	c.IncludeDBConnector = boolResult
	//todo: offer different DB connector options
	c.DBConnector = "postgres"

	// Mailer
	boolResult, err = survey.SelectYesNo("Do you want to include a GoMail mailer module?", true)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	c.IncludeMailer = boolResult
	//todo: offer different mailer options
	c.Mailer = "gomail"

	//----- git setup -----
	boolResult, err = survey.SelectYesNo("Do you want to set up git for the project?", true)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	c.SetUpGit = boolResult

	//----- docker setup -----
	boolResult, err = survey.SelectYesNo("Do you want to set up Dockerfile for the project? Note: if no is selected, docker-compose setup will be skipped", true)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	c.SetUpDockerFile = boolResult

	//----- docker-compose setup -----
	if !c.SetUpDockerFile {
		_, err = survey.SelectYesNo("Notice: docker-compose setup will be skipped because it's dependent on Dockerfile.", true)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		c.SetUpDockerCompose = false
		return c
	} else {
		boolResult, err = survey.SelectYesNo("Do you want a docker-compose setup for local development? This will set up a docker-compose file for a local postgres and server with volume mapping. You can add the frontend yourself if you want.", true)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		c.SetUpDockerCompose = boolResult
		return c
	}
}

func executeSetup(config *ProjectConfig) {
	var err error

	if config.Dir != "" {
		if _, err := os.Stat(config.Dir); os.IsNotExist(err) {
			err = os.Mkdir(config.Dir, 0755)
			if err != nil {
				log.Fatalf("Error creating directory: %v", err)
			}
		} else {
			log.Fatalf("Directory already exists: %s", config.Dir)
		}

		err = os.Chdir(config.Dir)
		if err != nil {
			log.Fatalf("Error changing directory: %v", err)
		}
	}

	err = createGoModule(config.Module)
	if err != nil {
		log.Fatalf("Error creating Go module: %v", err)
	}

	err = createMainFile(config)
	if err != nil {
		log.Fatalf("Error creating main file: %v", err)
	}

	if config.SetUpDockerFile {
		err = createDockerfile()
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

	if config.SetUpDockerCompose {
		err = createDockerCompose()
		if err != nil {
			log.Fatalf("Error creating Docker Compose file: %v", err)
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

func createMainFile(c *ProjectConfig) error {
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

	if err := tmpl.Execute(file, c); err != nil {
		log.Printf("Failed to execute template: %v", err)
		return err
	}

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func saveConfig(config *ProjectConfig) {
	tpl, err := template.ParseFS(templateFS, "templates/gogetter.yaml.tpl")
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	file, err := os.Create("gogetter.yaml")
	if err != nil {
		log.Fatalf("Failed to create config file: %v", err)
	}
	defer file.Close()

	err = tpl.Execute(file, config)
	if err != nil {
		log.Fatalf("Failed to execute template: %v", err)
	}

	log.Println("Configuration saved to gogetter.yaml")
}

func createDockerfile() error {
	templatePath := "templates/Dockerfile.tpl"

	tmpl, err := template.ParseFS(templateFS, templatePath)
	if err != nil {
		log.Printf("Failed to parse template: %v", err)
		return err
	}

	file, err := os.Create("Dockerfile")
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		return err
	}
	defer file.Close()

	if err := tmpl.Execute(file, nil); err != nil {
		log.Printf("Failed to execute template: %v", err)
		return err
	}

	return nil
}

func createGitRepository() error {
	log.Printf("Creating a new Git repository\n")

	cmd := exec.Command("git", "init")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

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

	err = tmpl.Execute(file, nil)
	if err != nil {
		log.Printf("Failed to execute template: %v", err)
		return err
	}

	return nil
}

func createDockerCompose() error {
	templatePath := "templates/docker-compose.yaml.tpl"

	tmpl, err := template.ParseFS(templateFS, templatePath)
	if err != nil {
		log.Printf("Failed to parse template: %v", err)
		return err
	}

	file, err := os.Create("docker-compose.yaml")
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		return err
	}
	defer file.Close()

	if err := tmpl.Execute(file, nil); err != nil {
		log.Printf("Failed to execute template: %v", err)
		return err
	}

	return nil
}
