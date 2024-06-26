package main

import (
	"go.uber.org/fx"

	config "github.com/alsey89/gogetter/config/viper"
	logger "github.com/alsey89/gogetter/logging/zap"
	server "github.com/alsey89/gogetter/server/echo"
	{{ if .IncludeJWTMiddleware }}jwt "github.com/alsey89/gogetter/jwt/echo" {{ end }}
	{{ if .IncludeMailer }}mailer "github.com/alsey89/gogetter/mail/gomail" {{ end }}
	{{ if .IncludeDBConnector }}database "github.com/alsey89/gogetter/database/postgres" {{ end }}
	// Internal domains can be imported below
)

var configuration *config.Module

func init() {
	//---------config---------
	//!PRECEDENCE: ENV > CONFIG FILE > FALLBACK > MODULE DEFAULTS
	config.SetSystemLogLevel("debug")
	configuration = config.SetUpConfig("SERVER", "yaml")
	configuration.SetFallbackConfigs(map[string]interface{}{
		//-----server-----
		"server.host":      "0.0.0.0",
		"server.port":      3001,
		"server.log_level": "DEV",
		//"server.allow_headers":   "*",
		//"server.allow_methods":   "*",
		//"server.allow_origins":   "http://localhost:3000, http://localhost:3001",
		//"server.csrf_protection": true,
		//"server.csrf_secure":     false,
		//"server.csrf_domain":     "localhost",

		{{- if .IncludeDBConnector }}

		//-----database-----
		//"database.host":         "postgres",
		//"database.port":         5432,
		//"database.dbname":       "postgres",
		//"database.user":         "postgres",
		//"database.password":     "password",
		//"database.sslmode":      "prefer",
		//"database.loglevel":     "error",
		//"database.auto_migrate": false,
		{{- end }}

		{{- if .IncludeMailer }}

		//-----mailer-----
		//"mailer.host":         "smtp.gmail.com",
		//"mailer.port":         587,
		//"mailer.username":     "example@example-gmail.com",
		//"mailer.app_password": "foo bar baz qux",
		//"mailer.tls":          true,
		{{- end }}

		{{- if .IncludeJWTMiddleware }}

		//-----echo_jwt-----
		//"echo_jwt.signing_key":    "authsecret",
		//"echo_jwt.token_lookup":   "cookie:jwt",
		//"echo_jwt.signing_method": "HS256",
		//"echo_jwt.exp_in_hours":   72,
		{{- end }}

		//-----internal domains-----
		// example
		//"auth.confirmation_signing_key":    "confirmationsecret",
		//"auth.confirmation_token_lookup":   "query:token",
		//"auth.confirmation_signing_method": "HS256",
		//"auth.confirmation_exp_in_hours":   1,
	})
}

func main() {
	app := fx.New(
		// ----- modules -----
		fx.Supply(configuration),
		logger.InitiateModule(),
		server.InitiateModule("server"),
		{{- if .IncludeDBConnector }}
		database.InitiateModuleAndSchema("database"),
		{{- end }}
		{{- if .IncludeJWTMiddleware }}
		jwt.InitiateModule("echo_jwt"),
		{{- end }}
		{{- if .IncludeMailer }}
		mailer.InitiateModule("mailer"),
		{{- end }}
		
		// ----- internal domains -----
		// auth.InitiateDomain("auth"),
		// company.InitiateDomain("company"),

		fx.NopLogger, //disable fx logs
	)
	app.Run()
}
