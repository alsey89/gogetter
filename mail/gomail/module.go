package mailer

import (
	"context"

	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"

	"github.com/alsey89/gogetter/common"
)

const (
	DefaultHost        = "0.0.0.0"
	DefaultPort        = 25
	DefaultUsername    = ""
	DefaultAppPassword = ""
	DefaultTLS         = false

	DefaultSubject = "From Gogetter Mail Module"
	DefaultBody    = "This is an email from Gogetter Mail Module."
	DefaultFrom    = "mail@gogetter.com"
	DefaultTo      = "mail@gogetter.com"
)

type Config struct {
	Host        string
	Port        int
	Username    string
	AppPassword string //! Important: Use app password, not the account password.
	TLS         bool
}

type Module struct {
	logger *zap.Logger
	config *Config

	scope  string
	dialer *gomail.Dialer
}

type Params struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *zap.Logger
}

func InitiateModule(scope string) fx.Option {

	var m *Module

	return fx.Module(
		scope,
		fx.Provide(func(p Params) *Module {
			logger := p.Logger.Named("[" + scope + "]")
			config := loadConfig(scope)
			dialer := gomail.NewDialer(config.Host, config.Port, config.Username, config.AppPassword)

			m := &Module{
				logger: logger,
				config: config,
				scope:  scope,
				dialer: dialer,
			}

			return m
		}),
		fx.Populate(&m),
		fx.Invoke(func(p Params) *Module {

			p.Lifecycle.Append(
				fx.Hook{
					OnStart: m.onStart,
					OnStop:  m.onStop,
				},
			)

			return m
		}),
	)

}

func loadConfig(scope string) *Config {

	//set defaults
	viper.SetDefault(common.GetConfigPath(scope, "host"), DefaultHost)
	viper.SetDefault(common.GetConfigPath(scope, "port"), DefaultPort)
	viper.SetDefault(common.GetConfigPath(scope, "username"), DefaultUsername)
	viper.SetDefault(common.GetConfigPath(scope, "password"), DefaultAppPassword)
	viper.SetDefault(common.GetConfigPath(scope, "tls"), DefaultTLS)
	//populate config
	return &Config{
		Host:        viper.GetString(common.GetConfigPath(scope, "host")),
		Port:        viper.GetInt(common.GetConfigPath(scope, "port")),
		Username:    viper.GetString(common.GetConfigPath(scope, "username")),
		AppPassword: viper.GetString(common.GetConfigPath(scope, "app_password")),
		TLS:         viper.GetBool(common.GetConfigPath(scope, "tls")),
	}
}

func (m *Module) onStart(ctx context.Context) error {
	m.logger.Info("Mailer initiated")

	err := m.TestSMTPConnection()
	if err != nil {
		m.logger.Error("Failed to connect to the SMTP server", zap.Error(err))
	}

	m.PrintDebugLogs()

	return nil
}

func (m *Module) onStop(ctx context.Context) error {

	m.logger.Info("Mailer module stopped")

	return nil
}

// ----------------------------------------------------------

func (m *Module) NewMessage() *gomail.Message {
	return gomail.NewMessage()
}

func (m *Module) Send(msg *gomail.Message) error {
	return m.dialer.DialAndSend(msg)
}

// creates and sends email
func (m *Module) SendTransactionalMail(from string, to string, subject string, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	err := m.dialer.DialAndSend(msg)
	if err != nil {
		m.logger.Error("Failed to send email", zap.Error(err))
		return err
	}

	m.logger.Info("Email sent successfully.")
	return nil
}

// ----------------------------------------------------------

func (m *Module) PrintDebugLogs() {
	//* Debug logs
	m.logger.Debug("----- Mailer Configuration -----")
	m.logger.Debug("Host", zap.String("Host", m.config.Host))
	m.logger.Debug("Port", zap.Int("Port", m.config.Port))
	m.logger.Debug("Username", zap.String("Username", m.config.Username))
	m.logger.Debug("AppPassword", zap.String("AppPassword", m.config.AppPassword))
	m.logger.Debug("TLS", zap.Bool("TLS", m.config.TLS))
}

func (m *Module) TestSMTPConnection() error {
	m.logger.Info("Testing SMTP connection...")
	s, err := m.dialer.Dial()
	if err != nil {
		m.logger.Error("Failed to connect to the SMTP server", zap.Error(err))
		return err
	}
	defer s.Close()

	m.logger.Info("Successfully connected to the SMTP server.")
	return nil
}

func (m *Module) SendTestMail(to, subject, body string) error {
	msg := m.NewMessage()
	msg.SetHeader("From", DefaultFrom)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	err := m.Send(msg)
	if err != nil {
		m.logger.Error("Failed to send test email", zap.Error(err))
		return err
	}

	m.logger.Info("Test email sent successfully.")
	return nil
}
