package postgres_db

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"
)

const (
	DefaultHost        = "0.0.0.0"
	DefaultPort        = 5432
	DefaultDbName      = "postgres"
	DefaultUser        = "postgres"
	DefaultPassword    = "password"
	DefaultSSLMode     = "allow"
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

type Database struct {
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
		fx.Provide(func(p Params) (*Database, error) {
			logger := p.Logger.Named("[" + scope + "]")
			config := loadConfig(scope)

			database := &Database{
				logger: logger,
				config: config,
				scope:  scope,
				schema: schema,
			}

			db := database.setUpDBConnectionOrFatal()

			database.db = db

			return database, nil
		}),
		fx.Invoke(func(d *Database, p Params) {
			p.Lifecycle.Append(
				fx.Hook{
					OnStart: d.onStart,
					OnStop:  d.onStop,
				},
			)
		}),
	)
}

func (d *Database) onStart(context.Context) error {
	d.logger.Info("Database initiated")
	if d.config.AutoMigrate {
		d.Migrate(d.schema...)
	}

	d.printDebugLogs()
	return nil
}

func (d *Database) onStop(context.Context) error {
	dbSql, err := d.db.DB()
	if err != nil {
		d.logger.Error("Error getting DB from GORM", zap.Error(err))
		return err
	}

	err = dbSql.Close()
	if err != nil {
		d.logger.Error("Error closing DB", zap.Error(err))
	}

	d.logger.Info("Database connection stopped")
	return nil
}

// ---------------------------------------------------------

func loadConfig(scope string) *Config {
	getConfigPath := func(key string) string {
		return fmt.Sprintf("%s.%s", scope, key)
	}

	//set default values
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

func (d *Database) setUpDBConnectionOrFatal() *gorm.DB {
	dsn := d.getConnectionStringFromConfig()
	loglevel := d.getLogLevelFromConfig()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gorm_logger.Default.LogMode(loglevel),
	})
	if err != nil {
		d.logger.Fatal("Error connecting to database", zap.Error(err))
	}

	return db
}
func (d *Database) getConnectionStringFromConfig() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.config.Host, d.config.Port, d.config.User, d.config.Password, d.config.DBName, d.config.SSLMode)
}
func (d *Database) getLogLevelFromConfig() gorm_logger.LogLevel {
	switch d.config.LogLevel {
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

func (d *Database) printDebugLogs() {
	//* Debug Logs
	d.logger.Debug("----- Database Configuration -----")
	d.logger.Debug("Host", zap.String("host", d.config.Host))
	d.logger.Debug("Port", zap.Int("port", d.config.Port))
	d.logger.Debug("DBName", zap.String("dbname", d.config.DBName))
	d.logger.Debug("User", zap.String("user", d.config.User))
	d.logger.Debug("SSLMode", zap.String("sslmode", d.config.SSLMode))
}

func (d *Database) Migrate(schema ...interface{}) {
	d.logger.Info("Auto migrating database")
	err := d.db.AutoMigrate(schema...)
	if err != nil {
		d.logger.Error("Error auto migrating database", zap.Error(err))
	}
}

// GETTERS ---------------------------------------------------------
func (d *Database) GetDB() *gorm.DB {
	return d.db
}
