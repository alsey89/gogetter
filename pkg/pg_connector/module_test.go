package pg_connector

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	gorm_logger "gorm.io/gorm/logger"
)

func TestConfigs(t *testing.T) {
	t.Run("TestLoadConfig", func(t *testing.T) {
		// set config values
		viper.Set("test_postgres_db.host", "testhost")
		viper.Set("test_postgres_db.port", 1234)
		viper.Set("test_postgres_db.dbname", "testdb")
		viper.Set("test_postgres_db.user", "testuser")
		viper.Set("test_postgres_db.password", "testpassword")
		viper.Set("test_postgres_db.sslmode", "testssl")
		viper.Set("test_postgres_db.log_level", "debug")
		viper.Set("test_postgres_db.auto_migrate", true)

		config := loadConfig("test_postgres_db")

		assert.Equal(t, "testhost", config.Host)
		assert.Equal(t, 1234, config.Port)
		assert.Equal(t, "testdb", config.DBName)
		assert.Equal(t, "testuser", config.User)
		assert.Equal(t, "testpassword", config.Password)
		assert.Equal(t, "testssl", config.SSLMode)
		assert.Equal(t, "debug", config.LogLevel)
		assert.True(t, config.AutoMigrate)

		viper.Reset()
	})

	t.Run("TestLoadConfigDefaults", func(t *testing.T) {
		//load config without config values
		config := loadConfig("test_postgres_db")

		assert.Equal(t, DefaultHost, config.Host)
		assert.Equal(t, DefaultPort, config.Port)
		assert.Equal(t, DefaultDbName, config.DBName)
		assert.Equal(t, DefaultUser, config.User)
		assert.Equal(t, DefaultPassword, config.Password)
		assert.Equal(t, DefaultSSLMode, config.SSLMode)
		assert.Equal(t, DefaultLogLevel, config.LogLevel)
		assert.Equal(t, DefaultAutoMigrate, config.AutoMigrate)

		viper.Reset()
	})
}

func TestGetConnectionStringFromConfig(t *testing.T) {
	d := &Module{
		config: &Config{
			Host:     "testhost",
			Port:     1234,
			User:     "testuser",
			Password: "testpassword",
			DBName:   "testdb",
			SSLMode:  "testssl",
		},
	}

	expected := "host=testhost port=1234 user=testuser password=testpassword dbname=testdb sslmode=testssl"
	actual := d.getConnectionStringFromConfig()

	assert.Equal(t, expected, actual)
}

func TestGetLogLevelFromConfig(t *testing.T) {
	d := &Module{
		config: &Config{
			LogLevel: "debug",
		},
	}
	t.Run("Silent", func(t *testing.T) {
		d.config.LogLevel = "silent"
		expected := gorm_logger.Silent
		actual := d.getLogLevelFromConfig()
		assert.Equal(t, expected, actual)
	})

	t.Run("Error", func(t *testing.T) {
		d.config.LogLevel = "error"
		expected := gorm_logger.Error
		actual := d.getLogLevelFromConfig()
		assert.Equal(t, expected, actual)
	})

	t.Run("Warn", func(t *testing.T) {
		d.config.LogLevel = "warn"
		expected := gorm_logger.Warn
		actual := d.getLogLevelFromConfig()
		assert.Equal(t, expected, actual)
	})

	t.Run("Info", func(t *testing.T) {
		d.config.LogLevel = "info"
		expected := gorm_logger.Info
		actual := d.getLogLevelFromConfig()
		assert.Equal(t, expected, actual)
	})

	t.Run("Default", func(t *testing.T) {
		d.config.LogLevel = "invalid"
		expected := gorm_logger.Info
		actual := d.getLogLevelFromConfig()
		assert.Equal(t, expected, actual)
	})
}
