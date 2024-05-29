package echo

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/alsey89/gogetter/common"
)

const (
	//server
	DefaultHost     = "0.0.0.0"
	DefaultLogLevel = "DEV"
	DefaultPort     = 3001
	//csrf middleware
	DefaultCSRFProtection = false
	DefaultCSRFSecure     = false
	DefaultCSRFDomain     = "localhost"
	//jwt middleware
	DefaultSigningKey    = "secret"
	DefaultTokenLookup   = "cookie:jwt"
	DefaultSigningMethod = "HS256"
)

type Config struct {
	AllowHeaders   string
	AllowMethods   string
	AllowOrigins   string
	CSRFProtection bool
	CSRFSecure     bool
	CSRFDomain     string
	Host           string
	LogLevel       string
	Port           int
}

type Module struct {
	config *Config
	logger *zap.Logger
	scope  string
	server *echo.Echo
}

type Params struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *zap.Logger
}

func InitiateModule(scope string) fx.Option {
	return fx.Module(
		scope,
		fx.Provide(func(p Params) *Module {
			logger := p.Logger.Named("[" + scope + "]")
			server := echo.New()
			config := loadConfig(scope)

			m := &Module{
				logger: logger,
				scope:  scope,

				config: config,
				server: server,
			}

			return m
		}),
		fx.Invoke(func(m *Module, p Params) {
			p.Lifecycle.Append(
				fx.Hook{
					OnStart: m.onStart,
					OnStop:  m.onStop,
				},
			)
		}),
	)
}

func loadConfig(scope string) *Config {
	//set defaults
	viper.SetDefault(common.GetConfigPath(scope, "allow_headers"), "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token, X-Requested-With, Origin, Cache-Control, Pragma, Expires, Set-Cookie, Cookie, jwt")
	viper.SetDefault(common.GetConfigPath(scope, "allow_methods"), "GET,PUT,POST,DELETE, OPTIONS, PATCH")
	viper.SetDefault(common.GetConfigPath(scope, "allow_origins"), "*")

	viper.SetDefault(common.GetConfigPath(scope, "csrf_protection"), DefaultCSRFProtection)
	viper.SetDefault(common.GetConfigPath(scope, "csrf_secure"), DefaultCSRFSecure)
	viper.SetDefault(common.GetConfigPath(scope, "csrf_domain"), DefaultCSRFDomain)

	viper.SetDefault(common.GetConfigPath(scope, "host"), DefaultHost)
	viper.SetDefault(common.GetConfigPath(scope, "log_level"), DefaultLogLevel)
	viper.SetDefault(common.GetConfigPath(scope, "port"), DefaultPort)

	return &Config{
		AllowHeaders: viper.GetString(common.GetConfigPath(scope, "allow_headers")),
		AllowMethods: viper.GetString(common.GetConfigPath(scope, "allow_methods")),
		AllowOrigins: viper.GetString(common.GetConfigPath(scope, "allow_origins")),

		CSRFProtection: viper.GetBool(common.GetConfigPath(scope, "csrf_protection")),
		CSRFSecure:     viper.GetBool(common.GetConfigPath(scope, "csrf_secure")),
		CSRFDomain:     viper.GetString(common.GetConfigPath(scope, "csrf_domain")),

		Host:     viper.GetString(common.GetConfigPath(scope, "host")),
		Port:     viper.GetInt(common.GetConfigPath(scope, "port")),
		LogLevel: viper.GetString(common.GetConfigPath(scope, "log_level")),
	}
}

func (m *Module) onStart(ctx context.Context) error {
	m.logger.Info("Server Module initiated")

	// middlewares
	m.setUpCorsMiddleware()
	m.setUpCSRFMiddleware()
	m.setUpRequestLoggerMiddleware()
	// server
	go m.startServer(true, true)

	m.PrintDebugLogs()

	return nil
}

func (m *Module) onStop(context.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.server.Shutdown(ctx)
	if err != nil {
		m.logger.Error("server shutdown error", zap.Error(err))
	}

	m.logger.Info("Server module stopped")
	return nil
}

func (m *Module) setUpCSRFMiddleware() {
	if !m.config.CSRFProtection {
		return
	}

	csrfConfig := middleware.CSRFConfig{
		TokenLookup:    "cookie:_csrf",
		CookiePath:     "/",
		CookieDomain:   m.config.CSRFDomain,
		CookieSecure:   m.config.CSRFSecure,
		CookieSameSite: http.SameSiteLaxMode,
		CookieHTTPOnly: true,
	}
	m.server.Use(middleware.CSRFWithConfig(csrfConfig))
}

func (m *Module) setUpCorsMiddleware() {
	// configure CORS middleware
	corsConfig := middleware.CORSConfig{
		AllowOrigins:     strings.Split(m.config.AllowOrigins, ","),
		AllowMethods:     strings.Split(m.config.AllowMethods, ","),
		AllowHeaders:     strings.Split(m.config.AllowHeaders, ","),
		AllowCredentials: true,
	}
	if m.config.AllowOrigins == "" || m.config.AllowOrigins == "*" {
		corsConfig.AllowOrigins = []string{"*"}
	}
	if m.config.AllowMethods == "" || m.config.AllowMethods == "*" {
		corsConfig.AllowMethods = []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete}
	}
	if m.config.AllowHeaders == "" || m.config.AllowHeaders == "*" {
		corsConfig.AllowHeaders = []string{"accept", "content-type", "authorization", "x-csrf-token", "x-requested-with", "origin", "cache-control", "pragma", "expires", "set-cookie", "cookie", "jwt"}
	}
	// add CORS middleware
	m.server.Use(middleware.CORSWithConfig(corsConfig))
}

func (m *Module) setUpRequestLoggerMiddleware() {

	// configure request logger according to log level
	requestLoggerConfig := middleware.RequestLoggerConfig{
		LogProtocol:  true,
		LogMethod:    true,
		LogURI:       true,
		LogStatus:    true,
		LogRequestID: true,
		LogRemoteIP:  true,
		LogLatency:   true,
		LogError:     true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			switch m.config.LogLevel {
			case "DEV":
				log.Printf("|--------------------------------------------\n")
				m.logger.Info("request",
					zap.String("URI", v.URI),
					zap.String("method", v.Method),
					zap.Int("status", v.Status),
					zap.Any("error", v.Error),
					zap.String("remote_ip", v.RemoteIP),
					zap.String("request_id", v.RequestID),
					zap.Duration("latency", v.Latency),
					zap.String("protocol", v.Protocol),
				)
				log.Printf("--------------------------------------------|\n")
			case "PROD":
				log.Printf("|--------------------------------------------\n")
				m.logger.Info("request",
					zap.String("URI", v.URI),
					zap.Int("status", v.Status),
					zap.Any("error", v.Error),
					zap.String("request_id", v.RequestID),
					zap.Duration("latency", v.Latency),
				)
				log.Printf("--------------------------------------------|\n")
			case "DEBUG":
				log.Printf("|--------------------------------------------\n")
				m.logger.Debug("request",
					zap.String("URI", v.URI),
					zap.String("method", v.Method),
					zap.Int("status", v.Status),
					zap.String("remote_ip", v.RemoteIP),
					zap.String("request_id", v.RequestID),
					zap.Duration("latency", v.Latency),
					zap.String("protocol", v.Protocol),
					zap.Any("error", v.Error),
					zap.Any("request_body", c.Request().Body),
					// todo: add more debug logs if needed
				)
				log.Printf("--------------------------------------------|\n")
			default:
				m.logger.Error("invalid log level", zap.String("log_level", m.config.LogLevel))
			}
			return nil
		},
	}
	// add request logger middleware
	m.server.Use(middleware.RequestLoggerWithConfig(requestLoggerConfig))
}

func (m *Module) startServer(HideBanner bool, HidePort bool) {
	m.server.HideBanner = HideBanner || false
	m.server.HidePort = HidePort || false

	addr := fmt.Sprintf("%s:%d", m.config.Host, m.config.Port)
	err := m.server.Start(addr)
	if err != nil && err != http.ErrServerClosed {
		m.logger.Fatal(err.Error())
	}
}

// ----------------------------------------------------------

func (m *Module) GetServer() *echo.Echo {
	return m.server
}

/*
Returns a pointer to an echo.MiddlewareFunc that provides JWT authentication middleware for Echo framework.
Middleware validates the JWT token, parses claims, and stores them in context under the key "user".
Takes in signing key, signing method, and token lookup.

	Example/Default values:
	- SigningKey = "secret"
	- SigningMethod = "HS256"
	- TokenLookup = "cookie:jwt"
*/
func (m *Module) GetEchoJWTMiddleware(signingKey string, signingMethod string, tokenLookup string) *echo.MiddlewareFunc {
	config := echojwt.Config{}

	if signingKey == "" {
		config.SigningKey = []byte(DefaultSigningKey)
	} else {
		config.SigningKey = []byte(signingKey)
	}

	if signingMethod == "" {
		config.SigningMethod = DefaultSigningMethod
	} else {
		config.SigningMethod = signingMethod
	}

	if tokenLookup == "" {
		config.TokenLookup = DefaultTokenLookup
	} else {
		config.TokenLookup = tokenLookup
	}

	middleware := echojwt.WithConfig(config)

	return &middleware
}

func (m *Module) PrintDebugLogs() {
	//* Debug Logs
	//server
	m.logger.Debug("----- Server Configuration -----")
	m.logger.Debug("Host", zap.String("Host", m.config.Host))
	m.logger.Debug("Port", zap.Int("Port", m.config.Port))
	//cors
	m.logger.Debug("----- Cors Configuration -----")
	m.logger.Debug("AllowOrigins", zap.String("AllowOrigins", m.config.AllowOrigins))
	m.logger.Debug("AllowMethods", zap.String("AllowMethods", m.config.AllowMethods))
	m.logger.Debug("AllowHeaders", zap.String("AllowHeaders", m.config.AllowHeaders))
	//csrf
	m.logger.Debug("----- CSRF Configuration -----")
	m.logger.Debug("CSRFProtection", zap.Bool("CSRFProtection", m.config.CSRFProtection))
	if m.config.CSRFProtection {
		m.logger.Debug("CSRFTokenLookup", zap.String("CSRFTokenLookup", "cookie:_csrf"))
		m.logger.Debug("CSRFCookiePath", zap.String("CSRFCookiePath", "/"))
		m.logger.Debug("CSRFCookieDomain", zap.String("CSRFCookieDomain", m.config.CSRFDomain))
		m.logger.Debug("CSRFCookieSecure", zap.Bool("CSRFCookieSecure", m.config.CSRFSecure))
		m.logger.Debug("CSRFCookieHTTPOnly", zap.Bool("CSRFCookieHTTPOnly", true))
	}
}
