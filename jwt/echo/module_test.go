// file: jwt/echo/module_test.go
package jwt

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestLoadConfig(t *testing.T) {
	viper.Set("test_echo_jwt.token_lookup", "header:auth")
	viper.Set("test_echo_jwt.signing_key", "secret")
	viper.Set("test_echo_jwt.signing_method", "HS256")
	viper.Set("test_echo_jwt.exp_in_hours", 48)

	config := loadConfig("test_echo_jwt")

	// Asserts
	assert.Equal(t, "header:auth", config.TokenLookup)
	assert.Equal(t, "secret", config.SigningKey)
	assert.Equal(t, "HS256", config.SigningMethod)
	assert.Equal(t, 48, config.ExpInHours)

	viper.Reset()
}

func TestLoadConfigDefaults(t *testing.T) {
	config := loadConfig("test_echo_jwt")

	// Asserts
	assert.Equal(t, defaultTokenLookup, config.TokenLookup)
	assert.Equal(t, defaultSecret, config.SigningKey)
	assert.Equal(t, defaultSigningMethod, config.SigningMethod)
	assert.Equal(t, defaultExpInHours, config.ExpInHours)

	viper.Reset()
}

func TestInitiateModule(t *testing.T) {
	logger := zaptest.NewLogger(t)
	// Setup
	app := fxtest.New(t,
		fx.Provide(func() *zap.Logger {
			return logger
		}),

		InitiateModule("test_echo_jwt"),
		fx.Invoke(func(JWT *JWT) {
			assert.NotNil(t, JWT, "JWT should be initialized")
		}),
	)

	// Start the application
	startCtx, cancelStart := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancelStart()
	require.NoError(t, app.Start(startCtx), "failed to start the app")

	// Stop the application
	stopCtx, cancelStop := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancelStop()
	require.NoError(t, app.Stop(stopCtx), "failed to stop the app")
}

func TestGenerateToken(t *testing.T) {
	j := &JWT{
		config: &Config{
			SigningKey:    "secret",
			SigningMethod: "HS256",
			ExpInHours:    1,
		},
	}

	claims := jwt.MapClaims{
		"username": "testuser",
	}
	token, err := j.GenerateToken(claims)
	assert.Nil(t, err)
	assert.NotNil(t, token)

	parsedToken, err := jwt.ParseWithClaims(*token, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.config.SigningKey), nil
	})
	assert.Nil(t, err)
	assert.True(t, parsedToken.Valid)

	claimsInToken := parsedToken.Claims.(*jwt.MapClaims)
	assert.Equal(t, "testuser", (*claimsInToken)["username"].(string))
	assert.WithinDuration(t, time.Now().Add(time.Hour*time.Duration(j.config.ExpInHours)), time.Unix(int64((*claimsInToken)["exp"].(float64)), 0), time.Minute)
}

func TestJWTMiddleware(t *testing.T) {
	logger := zaptest.NewLogger(t)
	testClaims := jwt.MapClaims{"user_id": 123}

	app := fxtest.New(t,
		fx.Provide(func() *zap.Logger { return logger }),
		InitiateModule("test_echo_jwt"),
		fx.Invoke(func(jwt *JWT) {
			e := echo.New()
			e.GET("/test", func(c echo.Context) error {
				return c.String(http.StatusOK, "Authorized")
			}, jwt.Middleware())

			t.Run("TestUnauthorizedAccess", func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				rec := httptest.NewRecorder()
				e.ServeHTTP(rec, req)
				assert.NotEqual(t, http.StatusOK, rec.Code, "Request should fail without a valid token")
			})

			t.Run("TestAuthorizedAccess", func(t *testing.T) {
				token, err := jwt.GenerateToken(testClaims)
				assert.NoError(t, err)
				assert.NotNil(t, token)

				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				req.AddCookie(&http.Cookie{Name: "jwt", Value: *token})
				rec := httptest.NewRecorder()
				e.ServeHTTP(rec, req)
				assert.Equal(t, http.StatusOK, rec.Code, "Request with valid token should succeed")
				assert.Equal(t, "Authorized", rec.Body.String())
			})

			t.Run("TestInvalidToken", func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				req.AddCookie(&http.Cookie{Name: "jwt", Value: "invalid_token"})
				rec := httptest.NewRecorder()
				e.ServeHTTP(rec, req)
				assert.NotEqual(t, http.StatusOK, rec.Code, "Request with invalid token should fail")
			})
		}),
	)
	app.RequireStart()
	app.RequireStop()
}
