package main

import (
	"go.uber.org/fx"

	"github.com/alsey89/gogetter/pkg/config"
	"github.com/alsey89/gogetter/pkg/logger"
	"github.com/alsey89/gogetter/pkg/mailer"
	"github.com/alsey89/gogetter/pkg/server"
	"github.com/alsey89/gogetter/pkg/token"
)

func init() {
	//! CONFIG PRECEDENCE: ENV > CONFIG FILE > FALLBACK
	config.SetSystemLogLevel("DEBUG")
	config.SetUpConfig("SERVER", "yaml", "./")
	config.SetFallbackConfigs(map[string]interface{}{
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

		"mailer.host":         "smtp.gmail.com",
		"mailer.port":         587,
		"mailer.username":     "example@example-gmail.com",
		"mailer.app_password": "foo bar baz qux",
		"mailer.tls":          true,

		"jwt_auth.signing_key":    "authsecret",
		"jwt_auth.token_lookup":   "cookie:jwt",
		"jwt_auth.signing_method": "HS256",
		"jwt_auth.exp_in_hours":   72,

		"jwt_email.signing_key":    "authsecret",
		"jwt_email.token_lookup":   "query:jwt",
		"jwt_email.signing_method": "HS256",
		"jwt_email.exp_in_hours":   1,

		"jwt_reset.signing_key":    "authsecret",
		"jwt_reset.token_lookup":   "query:jwt",
		"jwt_reset.signing_method": "HS256",
		"jwt_reset.exp_in_hours":   1,
	})
}
func main() {
	app := fx.New(
		logger.InjectModule("logger"),
		server.InjectModule("server"),
		// pgconn.InitiateModuleAndSchema(
		// 	"database",
		// // schema.User{},
		// // schema.ContactInfo{},
		// // schema.EmergencyContact{},
		// ),
		token.InjectModule("jwt", "jwt_auth", "jwt_email", "jwt_reset"),
		mailer.InjectModule("mailer", false),

		//-- Internal Domains Start --

		// auth.InitiateDomain("auth"),
		// company.InitiateDomain("company"),

		//-- Internal Domains End --

		fx.NopLogger,
	)
	app.Run()
}
