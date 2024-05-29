package main

import (
	"log"

	"go.uber.org/fx"

	config "github.com/alsey89/gogetter/config/viper"
	jwt "github.com/alsey89/gogetter/jwt"
	logger "github.com/alsey89/gogetter/logging/zap"
	mailer "github.com/alsey89/gogetter/mail/gomail"
	server "github.com/alsey89/gogetter/server/echo"

	jwtv5 "github.com/golang-jwt/jwt/v5"
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
		fx.Supply(configuration),
		logger.InitiateModule(),
		server.InitiateModule("server"),
		// postgres.InitiateModuleAndSchema(
		// 	"database",
		// // schema.User{},
		// // schema.ContactInfo{},
		// // schema.EmergencyContact{},
		// ),
		jwt.InitiateModule("jwt", "jwt_auth", "jwt_email", "jwt_reset"),
		mailer.InitiateModule("mailer"),

		//-- Internal Domains Start --

		// auth.InitiateDomain("auth"),
		// company.InitiateDomain("company"),

		//-- Internal Domains End --

		//manual testing of jwt module
		fx.Invoke(func(jwt *jwt.Module) {
			authToken, _ := jwt.GenerateToken("jwt_auth", jwtv5.MapClaims{"user_id": 11111})
			emailToken, _ := jwt.GenerateToken("jwt_email", jwtv5.MapClaims{"user_id": 22222})

			claims, err := jwt.ParseToken("jwt_auth", *authToken)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("claims: %v\n", claims)
			}

			claims, err = jwt.ParseToken("jwt_email", *emailToken)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("claims: %v\n", claims)
			}

		}),

		fx.NopLogger,
	)
	app.Run()
}
