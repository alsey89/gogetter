# GoGetter - Sensible Go Modules and CLI Tool

GoGetter is a collection of reusable Go modules, pre-configured to work with [Uber's Fx Framework](https://github.com/uber-go/fx). The aim is to provide a way to rapidly spin up service(s) using pre-built modules. This repository also includes a CLI tool to streamline common development tasks.

## All modules are dependent on:

- [Uber Fx](https://github.com/uber-go/fx)

## Modules

- **Configuration Management**
  - [Viper](https://github.com/spf13/viper)
- **Logging**
  - [Zap](https://github.com/uber-go/zap)
- **HTTPServer**
  - [Echo](https://echo.labstack.com/)
- **JWT**
  - WIP: General
  - [Echo JWT](https://echo.labstack.com/)
- **Database Integration** using [GORM](https://gorm.io/index.html)
  - [MySQL](https://www.mysql.com/)
  - [PostgreSQL](https://www.postgresql.org/)
- **Mailer**
  - [GoMail](https://github.com/go-gomail/gomail)

### Per Module Documentation

COMING SOON!

### Usage

To use the modules in your project, simply import them from main, set up the configurations, and initiate them in your Fx App.

Refer to [example.go](./example.go) for a working example.

### Configuration

There are 3 levels of configuration, listed in order of precedence:

1. **Environmental variables**:
   - Format: `prefix_scope_key`
   - Separator: `_` (underscore)
2. **Config Files**
   - Format: `scope.key`
   - Separator: `.`
3. **Fallback Config**
   - Format: `scope.key`
   - Separator: `.`

### Injection

Refer to [example.go](./example.go) for a working example.

## CLI

The CLI tool, built with [Cobra](https://github.com/spf13/cobra), is a convenient way to spin up an entire service in one go.

### Usage

To use the CLI tool, install it first:

```
go install github.com/alsey89/gogetter/cmd/gogetter@latest
```

#### Init

Init initializes the project. It sets up go module, creates a main.go file and installs the relevant dependencies. Optionally, it can set up a Dockerfile, a docker-compose.yaml, and git.

```
gogetter init
```

Here's an example of the process:

```
? Welcome to the GoGetter CLI. This will begin the setup process for your new Go service. Continue? Yes
? Enter the go module name for your project. [Example: github.com/alsey89/gogetter] test
? Enter the directory for your project. Service will be initiated at the current directory if left empty.
? Do you want to include a Echo-JWT middleware module? Yes
? Do you want to include a GORM Postgres database connector module? Yes
? Do you want to include a GoMail mailer module? Yes
? Do you want to set up git for the project? Yes
? Do you want to set up Dockerfile for the project? Note: if no is selected, docker-compose setup will be skipped Yes
? Do you want a docker-compose setup for local development? This will set up a docker-compose file for a local postgres and server with volume mapping. You can add the frontend yourself if you want. Yes
```

#### Run

Run spins up the docker-compose service, defaulting to a dev setup with automatic rebuild and reload.

```
gogetter run dev
```

Arguments:

- dev: sets BUILD_ENV=development
- development: sets BUILD_ENV=development
- prod: sets BUILD_ENV=production
- production: sets BUILD_ENV=production

Effects:
Check the [Dockerfile template](./cmd/templates/Dockerfile.tpl) to see how the BUILD_ENV affects the container setup.

### Troubleshooting

If the command is not found after installation, check Go Environmental variables and system $PATH.

## Contribution

Contributions are welcome! Please fork the repository and submit pull requests with your proposed changes. For major changes, please open an issue first to discuss what you would like to change.

Ensure to update tests as appropriate.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
