package jwt_manager

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestLoadConfig(t *testing.T) {
	scope := "test_scope"
	viper.Set("test_scope.token_lookup", "query:token")
	viper.Set("test_scope.signing_key", "testkey")
	viper.Set("test_scope.signing_method", "HS256")
	viper.Set("test_scope.exp_in_hours", 1)

	config := loadConfig(scope)

	assert.Equal(t, "query:token", config.TokenLookup)
	assert.Equal(t, "testkey", config.SigningKey)
	assert.Equal(t, "HS256", config.SigningMethod)
	assert.Equal(t, 1, config.ExpInHours)
}

func TestGenerateToken(t *testing.T) {
	config := &Config{
		TokenLookup:   "query:token",
		SigningKey:    "testkey",
		SigningMethod: "HS256",
		ExpInHours:    1,
	}

	logger, _ := zap.NewProduction()
	module := &Module{
		logger:  logger,
		configs: map[string]*Config{"test_scope": config},
	}

	claims := jwt.MapClaims{"user_id": 123}
	token, err := module.GenerateToken("test_scope", claims)
	assert.NoError(t, err)
	assert.NotNil(t, token)
}
