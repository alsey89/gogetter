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

### Usage

To use the modules in your project, simply import them from main, set up the configurations, and initiate them in your Fx App.

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

Refer to the template/example file for a working example.

### Per Module Documentation

COMING SOON!

## CLI

The CLI tool, built with [Cobra](https://github.com/spf13/cobra), is a convenient way to spin up an entire service in one go.

## Contribution

Contributions are welcome! Please fork the repository and submit pull requests with your proposed changes. For major changes, please open an issue first to discuss what you would like to change.

Ensure to update tests as appropriate.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
