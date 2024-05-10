package database

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"
)

const (
	DefaultHost        = "0.0.0.0"
	DefaultPort        = 3306
	DefaultDbName      = "mysql"
	DefaultUser        = "root"
	DefaultPassword    = "password"
	DefaultSSLMode     = "true"
	DefaultLogLevel    = "info"
	DefaultAutoMigrate = false
)

type Config struct {
	Host        string
	Port        int
	DBName      string
	User        string
	Password    string
	SSLMode     string
	LogLevel    string
	AutoMigrate bool
}

type Module struct {
	logger *zap.Logger
	config *Config
	schema []interface{}

	scope string
	db    *gorm.DB
}

type Params struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *zap.Logger
}

func InitiateModuleAndSchema(scope string, schema ...interface{}) fx.Option {
	return fx.Module(
		scope,
		fx.Provide(func(p Params) (*Module, error) {
			logger := p.Logger.Named("[" + scope + "]")
			config := loadConfig(scope)

			database := &Module{
				logger: logger,
				config: config,
				scope:  scope,
				schema: schema,
			}

			db := database.setUpDBConnectionOrFatal()

			database.db = db

			return database, nil
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

func (m *Module) onStart(context.Context) error {
	m.logger.Info("Database initiated")
	if m.config.AutoMigrate {
		m.Migrate(m.schema...)
	}

	m.printDebugLogs()
	return nil
}

func (m *Module) onStop(context.Context) error {
	dbSql, err := m.db.DB()
	if err != nil {
		m.logger.Error("Error getting DB from GORM", zap.Error(err))
		return err
	}

	err = dbSql.Close()
	if err != nil {
		m.logger.Error("Error closing DB", zap.Error(err))
	}

	m.logger.Info("Database connection stopped")
	return nil
}

func loadConfig(scope string) *Config {
	getConfigPath := func(key string) string {
		return fmt.Sprintf("%s.%s", scope, key)
	}

	viper.SetDefault(getConfigPath("host"), DefaultHost)
	viper.SetDefault(getConfigPath("port"), DefaultPort)
	viper.SetDefault(getConfigPath("dbname"), DefaultDbName)
	viper.SetDefault(getConfigPath("user"), DefaultUser)
	viper.SetDefault(getConfigPath("password"), DefaultPassword)
	viper.SetDefault(getConfigPath("sslmode"), DefaultSSLMode)
	viper.SetDefault(getConfigPath("log_level"), DefaultLogLevel)
	viper.SetDefault(getConfigPath("auto_migrate"), DefaultAutoMigrate)

	return &Config{
		Host:        viper.GetString(getConfigPath("host")),
		Port:        viper.GetInt(getConfigPath("port")),
		DBName:      viper.GetString(getConfigPath("dbname")),
		User:        viper.GetString(getConfigPath("user")),
		Password:    viper.GetString(getConfigPath("password")),
		SSLMode:     viper.GetString(getConfigPath("sslmode")),
		LogLevel:    viper.GetString(getConfigPath("log_level")),
		AutoMigrate: viper.GetBool(getConfigPath("auto_migrate")),
	}
}

func (m *Module) setUpDBConnectionOrFatal() *gorm.DB {
	dsn := m.getConnectionStringFromConfig()
	loglevel := m.getLogLevelFromConfig()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gorm_logger.Default.LogMode(loglevel),
	})
	if err != nil {
		m.logger.Fatal("Error connecting to database", zap.Error(err))
	}

	return db
}

func (m *Module) getConnectionStringFromConfig() string {
	return fmt.Sprintf("root:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.config.Password, m.config.Host, m.config.Port, m.config.DBName)
}

func (m *Module) getLogLevelFromConfig() gorm_logger.LogLevel {
	switch m.config.LogLevel {
	case "silent":
		return gorm_logger.Silent
	case "error":
		return gorm_logger.Error
	case "warn":
		return gorm_logger.Warn
	case "info":
		return gorm_logger.Info
	default:
		return gorm_logger.Info
	}
}

func (m *Module) printDebugLogs() {
	//* Debug Logs
	m.logger.Debug("----- Database Configuration -----")
	m.logger.Debug("Host", zap.String("host", m.config.Host))
	m.logger.Debug("Port", zap.Int("port", m.config.Port))
	m.logger.Debug("DBName", zap.String("dbname", m.config.DBName))
	m.logger.Debug("User", zap.String("user", m.config.User))
	m.logger.Debug("SSLMode", zap.String("sslmode", m.config.SSLMode))
}

func (m *Module) Migrate(schema ...interface{}) {
	m.logger.Info("Auto migrating database")
	err := m.db.AutoMigrate(schema...)
	if err != nil {
		m.logger.Error("Error auto migrating database", zap.Error(err))
	}
}

// GETTERS ---------------------------------------------------------
func (m *Module) GetDB() *gorm.DB {
	return m.db
}
