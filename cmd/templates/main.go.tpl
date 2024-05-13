package main

import (
	"go.uber.org/fx"

	"github.com/alsey89/gogetter/config/viper"
	"github.com/alsey89/gogetter/logging/zap"
	"github.com/alsey89/gogetter/server/echo"
	{{ if .IncludeJWTMiddleware }} "github.com/alsey89/gogetter/jwt/echo" {{ end }}
	{{ if .IncludeMailer }} "github.com/alsey89/gogetter/mail/gomail" {{ end }}
	{{ if .IncludeDBConnector }} "github.com/alsey89/gogetter/database/postgres" {{ end }}
	// Internal domains can be imported below
)

var configuration *viper.Module

func init() {
	viper.SetSystemLogLevel("debug")
	configuration = viper.SetUpConfig("SERVER", "yaml")
	//! CONFIG PRECEDENCE: ENV > CONFIG FILE > FALLBACK
	configuration.SetFallbackConfigs(map[string]interface{}{
		"server.host":      "0.0.0.0",
		"server.port":      3001,
		"server.log_level": "DEV",

		"server.allow_headers":   "*",
		"server.allow_methods":   "*",
		"server.allow_origins":   "http://localhost:3000, http://localhost:3001",
		"server.csrf_protection": true,
		"server.csrf_secure":     false,
		"server.csrf_domain":     "localhost",

		{{- if .IncludeDBConnector }}
		"database.host":         "postgres",
		"database.port":         5432,
		"database.dbname":       "postgres",
		"database.user":         "postgres",
		"database.password":     "password",
		"database.sslmode":      "prefer",
		"database.loglevel":     "error",
		"database.auto_migrate": false,
		{{- end }}
		{{- if .IncludeMailer }}
		"mailer.host":         "smtp.gmail.com",
		"mailer.port":         587,
		"mailer.username":     "example@example-gmail.com",
		"mailer.app_password": "foo bar baz qux",
		"mailer.tls":          true,
		{{- end }}
		{{- if .IncludeJWTMiddleware }}
		"echo_jwt.signing_key":    "authsecret",
		"echo_jwt.token_lookup":   "cookie:jwt",
		"echo_jwt.signing_method": "HS256",
		"echo_jwt.exp_in_hours":   72,
		{{- end }}
	})
}

func main() {
	app := fx.New(
		fx.Supply(configuration),
		zap.InitiateModule(),
		echo.InitiateModule("server"),
		{{- if .IncludeDBConnector }}
		postgres.InitiateModuleAndSchema("database"),
		{{- end }}
		{{- if .IncludeJWTMiddleware }}
		echo_jwt.InitiateModule("echo_jwt"),
		{{- end }}
		{{- if .IncludeMailer }}
		gomail.InitiateModule("mailer"),
		{{- end }}
		
		// Internal domains can be included here
		// auth.InitiateDomain("auth"),
		// company.InitiateDomain("company"),
		// Internal domains end here

		fx.NopLogger,
	)
	app.Run()
}
