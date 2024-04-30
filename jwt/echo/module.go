package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	defaultSecret        = "secret"
	defaultTokenLookup   = "cookie:jwt"
	defaultSigningMethod = "HS256"
	defaultExpInHours    = 72
)

type Config struct {
	TokenLookup   string
	SigningKey    string
	SigningMethod string
	ExpInHours    int
}

type JWT struct {
	logger *zap.Logger
	config *Config

	scope      string
	middleware echo.MiddlewareFunc
}

type Params struct {
	fx.In

	Logger    *zap.Logger
	Lifecycle fx.Lifecycle
}

// InitiateModule initializes the JWT middleware with the provided scope.
// CONFIG --> scope.token_lookup, scope.signing_key, scope.signing_method, scope.exp_in_hours
func InitiateModule(scope string) fx.Option {
	return fx.Module(
		scope,
		fx.Provide(func(p Params) (*JWT, error) {
			logger := p.Logger.Named("[" + scope + "]")
			config := loadConfig(scope)

			middleware := echojwt.WithConfig(echojwt.Config{
				SigningKey:    []byte(config.SigningKey),
				SigningMethod: config.SigningMethod,
				TokenLookup:   config.TokenLookup,
				ErrorHandler: func(c echo.Context, err error) error {
					logger.Error("JWT Middleware Error", zap.Error(err))
					return err
				},
			})

			j := &JWT{
				logger:     logger,
				config:     config,
				scope:      scope,
				middleware: middleware,
			}

			return j, nil
		}),
		fx.Invoke(func(j *JWT, p Params) {
			p.Lifecycle.Append(
				fx.Hook{
					OnStart: j.onStart,
					OnStop:  j.onStop,
				},
			)
		}),
	)
}

func loadConfig(scope string) *Config {
	getConfigPath := func(key string) string {
		return fmt.Sprintf("%s.%s", scope, key)
	}

	//set defaults
	viper.SetDefault(getConfigPath("token_lookup"), defaultTokenLookup)
	viper.SetDefault(getConfigPath("signing_key"), defaultSecret)
	viper.SetDefault(getConfigPath("signing_method"), defaultSigningMethod)
	viper.SetDefault(getConfigPath("exp_in_hours"), defaultExpInHours)

	return &Config{
		TokenLookup:   viper.GetString(getConfigPath("token_lookup")),
		SigningKey:    viper.GetString(getConfigPath("signing_key")),
		SigningMethod: viper.GetString(getConfigPath("signing_method")),
		ExpInHours:    viper.GetInt(getConfigPath("exp_in_hours")),
	}
}

func (j *JWT) onStart(ctx context.Context) error {
	j.logger.Info("JWT Middleware initiated")

	j.PrintDebugLogs()

	return nil
}

func (j *JWT) onStop(ctx context.Context) error {
	j.logger.Info("Stopping JWT")
	return nil
}

//! ------------------------------------------------------------

func (j *JWT) PrintDebugLogs() {
	j.logger.Debug("----- JWT Middleware Configuration -----")
	j.logger.Debug("TokenLookup", zap.String("TokenLookup", j.config.TokenLookup))
	j.logger.Debug("SigningKey", zap.Any("SigningKey", j.config.SigningKey))
	j.logger.Debug("SigningMethod", zap.String("SigningMethod", j.config.SigningMethod))
}

// Middleware returns the echo.MiddlewareFunc for JWT authentication.
func (j *JWT) Middleware() echo.MiddlewareFunc {
	return j.middleware
}

// GenerateToken generates a JWT token with the provided additional claims.
// It takes additionalClaims as input, which is a map of custom claims to be added to the token.
// SigningKey, SigningMethod, and ExpInHours are taken from the config keys scope.signing_key, scope.signing_method, and scope.exp_in_hours respectively.
// The function returns a pointer to the generated token string and an error, if any.
func (j *JWT) GenerateToken(additionalClaims jwt.MapClaims) (*string, error) {

	claims := jwt.MapClaims{
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(j.config.ExpInHours))),
	}

	for key, value := range additionalClaims {
		claims[key] = value
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(j.config.SigningMethod), claims)
	t, err := token.SignedString([]byte(j.config.SigningKey))
	if err != nil {
		j.logger.Error("Failed to generate token", zap.Error(err))
		return nil, err
	}
	return &t, nil
}
