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
	Log      LogConfig
}

type ServerConfig struct {
	CreatePort   int           `envconfig:"SERVER_CREATE_PORT"`
	CallbackPort int           `envconfig:"SERVER_CALLBACK_PORT"`
	ReadTimeout  time.Duration `envconfig:"SERVER_READ_TIMEOUT"`
	WriteTimeout time.Duration `envconfig:"SERVER_WRITE_TIMEOUT"`
}

type DatabaseConfig struct {
	Host     string `envconfig:"DB_HOST"`
	Port     int    `envconfig:"DB_PORT"`
	User     string `envconfig:"DB_USER"`
	Password string `envconfig:"DB_PASSWORD"`
	DBName   string `envconfig:"DB_NAME"`
	SSLMode  string `envconfig:"DB_SSLMODE"`
	MaxConn  int    `envconfig:"DB_MAX_CONN"`
}

type SSHConfig struct {
	Timeout time.Duration `envconfig:"SSH_TIMEOUT"`
	Port    int           `envconfig:"SSH_PORT"`
}

type CallbackConfig struct {
	Token string `envconfig:"CALLBACK_TOKEN" required:"true"`
}

type LogConfig struct {
	Level  string `envconfig:"LOG_LEVEL"`
	Format string `envconfig:"LOG_FORMAT"`
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
		return fmt.Errorf("invalid LOG_LEVEL: %s (must be debug, info, warn, error)", c.Log.Level)
	}

	if c.Log.Format != "json" && c.Log.Format != "text" {
		return fmt.Errorf("invalid LOG_FORMAT: %s (must be json or text)", c.Log.Format)
	}

	return nil
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}
