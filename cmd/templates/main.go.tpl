package main

import (
	"go.uber.org/fx"
	config "github.com/alsey89/gogetter/config/viper"
	logger "github.com/alsey89/gogetter/logging/zap"
	server "github.com/alsey89/gogetter/server/echo"
	{{ if .IncludeDBConnector }}postgres "github.com/alsey89/gogetter/database/postgres" {{ end }}
	{{ if .IncludeJWTMiddleware }}jwt "github.com/alsey89/gogetter/jwt/echo" {{ end }}
	{{ if .IncludeMailer }}mailer "github.com/alsey89/gogetter/mail/gomail" {{ end }}
)

var configuration *viper.Config

func init() {
	config.SetSystemLogLevel("debug")
	configuration = config.SetUpConfig("SERVICE", "yaml")
	configuration.SetFallbackConfigs(map[string]interface{}{
		// Server configuration
		"server.host": "0.0.0.0",
		"server.port": 3001,
		{{- if .IncludeDBConnector }}

		// Database configuration
		"database.host":     "0.0.0.0",
		"database.port":     5432,
		"database.dbname":   "postgres",
		"database.user":     "postgres",
		"database.password": "password",
		"database.sslmode":  "prefer",
		"databse.loglevel":  "error",
		{{- end }}
		{{- if .IncludeMailer }}

		// Mailer configuration
		"mailer.host":         "smtp.gmail.com",
		"mailer.port":         587,
		"mailer.username":     "your.email@gmail.com",
		"mailer.app_password": "your_app_password",
		"mailer.tls":          true,
		{{- end }}
		{{- if .IncludeJWTMiddleware }}

		// JWT configuration
		"auth_jwt.signing_key":  "thisisasecret",
		"auth_jwt.token_lookup": "cookie:jwt",
		{{- end }}
	})
}

func main() {
	app := fx.New(
		fx.Supply(configuration),
		logger.InitiateModule(),
		{{- if .IncludeDBConnector }}

		postgres.InitiateModuleAndSchema(
			"database",
			// ...schema,
			// example: &User{},
			// example: &Post{},
			// example: &Comment{},
		),
		{{- end }}
		{{- if .IncludeJWTMiddleware }}

		jwt.InitiateModule("auth_jwt"),
		{{- end }}
		{{- if .IncludeMailer }}
		
		mailer.InitiateModule("mailer"),
		{{- end }}

		server.InitiateModule("server"),
		fx.NopLogger,
	)
	app.Run()
}
