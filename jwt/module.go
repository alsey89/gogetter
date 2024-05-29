package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/alsey89/gogetter/common"
)

const (
	defaultSigningKey    = "secret"
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

type Module struct {
	scope   string
	logger  *zap.Logger
	configs map[string]*Config
}

type Params struct {
	fx.In

	Logger    *zap.Logger
	Lifecycle fx.Lifecycle
}

/*
Initializes the JWT module with the provided jwtScopes.
Takes multiple jwtScopes and loads configuration for each.
Tokens can be generated for each scope.

For example:
- jwt.GenerateToken("auth", jwt.MapClaims{"user_id": 123})
- jwt.GenerateToken("email", jwt.MapClaims{"email": "john@doe.com"})
*/
func InitiateModule(moduleScope string, jwtScopes ...string) fx.Option {
	return fx.Module(
		"jwt",
		fx.Provide(func(p Params) (*Module, error) {
			logger := p.Logger.Named("[" + moduleScope + "]")
			configs := make(map[string]*Config)

			for _, scope := range jwtScopes {
				config := loadConfig(scope)
				configs[scope] = config
			}

			m := &Module{
				scope:   moduleScope,
				logger:  logger,
				configs: configs,
			}

			return m, nil
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

	// Set defaults
	viper.SetDefault(common.GetConfigPath(scope, "token_lookup"), defaultTokenLookup)
	viper.SetDefault(common.GetConfigPath(scope, "signing_key"), defaultSigningKey)
	viper.SetDefault(common.GetConfigPath(scope, "signing_method"), defaultSigningMethod)
	viper.SetDefault(common.GetConfigPath(scope, "exp_in_hours"), defaultExpInHours)

	return &Config{
		TokenLookup:   viper.GetString(common.GetConfigPath(scope, "token_lookup")),
		SigningKey:    viper.GetString(common.GetConfigPath(scope, "signing_key")),
		SigningMethod: viper.GetString(common.GetConfigPath(scope, "signing_method")),
		ExpInHours:    viper.GetInt(common.GetConfigPath(scope, "exp_in_hours")),
	}
}

func (m *Module) onStart(ctx context.Context) error {
	m.logger.Info("JWT Module initiated")

	m.PrintDebugLogs()

	return nil
}

func (m *Module) onStop(ctx context.Context) error {
	m.logger.Info("JWT Module stopped")
	return nil
}

/*
Prints debug logs for all configurations.
*/
func (m *Module) PrintDebugLogs() {
	for scope, config := range m.configs {
		m.logger.Debug("----- JWT Module Configuration -----", zap.String("Scope", scope))
		m.logger.Debug("TokenLookup", zap.String("TokenLookup", config.TokenLookup))
		m.logger.Debug("SigningKey", zap.String("SigningKey", config.SigningKey))
		m.logger.Debug("SigningMethod", zap.String("SigningMethod", config.SigningMethod))
		m.logger.Debug("ExpInHours", zap.Int("ExpInHours", config.ExpInHours))
	}
}

/*
Generates a JWT token with the provided additional claims for a specific scope.
Use jwt.MapClaims from "github.com/golang-jwt/jwt/v5"
*/
func (m *Module) GenerateToken(scope string, additionalClaims jwt.MapClaims) (*string, error) {
	config, exists := m.configs[scope]
	if !exists {
		m.logger.Error("Config not found", zap.String("Scope", scope))
		return nil, fmt.Errorf("config for scope %s not found", scope)
	}

	claims := jwt.MapClaims{
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(config.ExpInHours))),
	}

	for key, value := range additionalClaims {
		claims[key] = value
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(config.SigningMethod), claims)
	t, err := token.SignedString([]byte(config.SigningKey))
	if err != nil {
		m.logger.Error("Failed to generate token", zap.Error(err))
		return nil, err
	}
	return &t, nil
}

/*
parse JWT token for a specific scope
*/
func (m *Module) ParseToken(scope string, tokenString string) (jwt.MapClaims, error) {
	config, exists := m.configs[scope]
	if !exists {
		m.logger.Error("Config not found", zap.String("Scope", scope))
		return nil, fmt.Errorf("config for scope %s not found", scope)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(config.SigningMethod) != token.Method {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.SigningKey), nil
	}, jwt.WithValidMethods([]string{config.SigningMethod}))

	if err != nil {
		m.logger.Error("Failed to parse token", zap.Error(err))
		return nil, err
	}

	if !token.Valid {
		m.logger.Error("Invalid token")
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

/*
Retrieves the configuration for a specific scope.
*/
func (m *Module) GetConfig(scope string) (*Config, error) {
	config, exists := m.configs[scope]
	if !exists {
		return nil, fmt.Errorf("config for scope %s not found", scope)
	}
	return config, nil
}
