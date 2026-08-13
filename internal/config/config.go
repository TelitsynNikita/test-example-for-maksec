package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	SSH      SSHConfig
	Callback CallbackConfig
	Agent    AgentConfig
	Log      LogConfig
}

type ServerConfig struct {
	CreatePort   int           `envconfig:"SERVER_CREATE_PORT" default:"8080"`
	CallbackPort int           `envconfig:"SERVER_CALLBACK_PORT" default:"8081"`
	ReadTimeout  time.Duration `envconfig:"SERVER_READ_TIMEOUT" default:"10s"`
	WriteTimeout time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" default:"10s"`
}

type DatabaseConfig struct {
	Host     string `envconfig:"DB_HOST" default:"localhost"`
	Port     int    `envconfig:"DB_PORT" default:"5432"`
	User     string `envconfig:"DB_USER" default:"postgres"`
	Password string `envconfig:"DB_PASSWORD" required:"true"`
	DBName   string `envconfig:"DB_NAME" default:"script_monitor"`
	SSLMode  string `envconfig:"DB_SSLMODE" default:"disable"`
	MaxConn  int    `envconfig:"DB_MAX_CONN" default:"10"`
}

type SSHConfig struct {
	Timeout               time.Duration `envconfig:"SSH_TIMEOUT" default:"30s"`
	Port                  int           `envconfig:"SSH_PORT" default:"22"`
	StrictHostKeyChecking bool          `envconfig:"SSH_STRICT_HOST_KEY_CHECKING" default:"false"`
	KnownHostsFile        string        `envconfig:"SSH_KNOWN_HOSTS" default:""`
}

type CallbackConfig struct {
	Token string `envconfig:"CALLBACK_TOKEN" required:"true"`
}

type AgentConfig struct {
	CallbackURL   string `envconfig:"AGENT_CALLBACK_URL" default:"http://localhost:8081/callback"`
	CallbackToken string `envconfig:"AGENT_CALLBACK_TOKEN" default:"secret"`
}

type LogConfig struct {
	Level  string `envconfig:"LOG_LEVEL" default:"info"`
	Format string `envconfig:"LOG_FORMAT" default:"json"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to process env config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.Server.CreatePort < 1 || c.Server.CreatePort > 65535 {
		return fmt.Errorf("invalid SERVER_CREATE_PORT: %d", c.Server.CreatePort)
	}
	if c.Server.CallbackPort < 1 || c.Server.CallbackPort > 65535 {
		return fmt.Errorf("invalid SERVER_CALLBACK_PORT: %d", c.Server.CallbackPort)
	}
	if c.Server.CreatePort == c.Server.CallbackPort {
		return fmt.Errorf("SERVER_CREATE_PORT and SERVER_CALLBACK_PORT must be different")
	}

	if c.Callback.Token == "" {
		return fmt.Errorf("CALLBACK_TOKEN is required")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[c.Log.Level] {
		return fmt.Errorf("invalid LOG_LEVEL: %s", c.Log.Level)
	}

	if c.Log.Format != "json" && c.Log.Format != "text" {
		return fmt.Errorf("invalid LOG_FORMAT: %s", c.Log.Format)
	}

	return nil
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}
