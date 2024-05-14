package main

import (
	"go.uber.org/fx"

	config "github.com/alsey89/gogetter/config/viper"
	jwt "github.com/alsey89/gogetter/jwt/echo"
	logger "github.com/alsey89/gogetter/logging/zap"
	mailer "github.com/alsey89/gogetter/mail/gomail"
	server "github.com/alsey89/gogetter/server/echo"
)

var configuration *config.Module

func init() {
	//! CONFIG PRECEDENCE: ENV > CONFIG FILE > FALLBACK
	config.SetSystemLogLevel("debug")
	configuration = config.SetUpConfig("SERVER", "yaml")
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

		// Database
		"database.host":         "postgres",
		"database.port":         5432,
		"database.dbname":       "postgres",
		"database.user":         "postgres",
		"database.password":     "password",
		"database.sslmode":      "prefer",
		"databse.loglevel":      "error",
		"database.auto_migrate": false,

		// Mailer
		"mailer.host":         "smtp.gmail.com",
		"mailer.port":         587,
		"mailer.username":     "example@example-gmail.com",
		"mailer.app_password": "foo bar baz qux",
		"mailer.tls":          true,

		// Echo JWT
		"echo_jwt.signing_key":    "authsecret",
		"echo_jwt.token_lookup":   "cookie:jwt",
		"echo_jwt.signing_method": "HS256",
		"echo_jwt.exp_in_hours":   72,
	})
}
func main() {
	app := fx.New(
		fx.Supply(configuration),
		logger.InitiateModule(),
		server.InitiateModule("server"),
		// postgres.InitiateModuleAndSchema(
		// 	"database",
		// // schema.User{},
		// // schema.ContactInfo{},
		// // schema.EmergencyContact{},
		// ),
		jwt.InitiateModule("echo_jwt"),
		mailer.InitiateModule("mailer"),

		//-- Internal Domains Start --

		// auth.InitiateDomain("auth"),
		// company.InitiateDomain("company"),

		//-- Internal Domains End --

		fx.NopLogger,
	)
	app.Run()
}
