package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

// Config represents the application configuration loaded from environment variables.
type Config struct {
	// Server configuration
	Server ServerConfig `envconfig:"SERVER"`

	// Database configuration
	Database DatabaseConfig `envconfig:"DATABASE"`

	// Logging configuration
	Logging LoggingConfig `envconfig:"LOGGING"`

	// Environment
	Environment string `envconfig:"ENVIRONMENT" default:"development"`

	// Debug mode
	Debug bool `envconfig:"DEBUG" default:"false"`
}

// ServerConfig represents server-specific configuration.
type ServerConfig struct {
	// Port to listen on
	Port int `envconfig:"PORT" default:"8080"`

	// Host to bind to
	Host string `envconfig:"HOST" default:"localhost"`

	// Read timeout in seconds
	ReadTimeout int `envconfig:"READ_TIMEOUT" default:"30"`

	// Write timeout in seconds
	WriteTimeout int `envconfig:"WRITE_TIMEOUT" default:"30"`

	// Idle timeout in seconds
	IdleTimeout int `envconfig:"IDLE_TIMEOUT" default:"60"`

	// Shutdown timeout in seconds (for graceful shutdown)
	ShutdownTimeout int `envconfig:"SHUTDOWN_TIMEOUT" default:"30"`
}

// DatabaseConfig represents database-specific configuration.
type DatabaseConfig struct {
	// Database host
	Host string `envconfig:"HOST" default:"localhost"`

	// Database port
	Port int `envconfig:"PORT" default:"5432"`

	// Database name
	Name string `envconfig:"NAME" required:"true"`

	// Database user
	User string `envconfig:"USER" required:"true"`

	// Database password
	Password string `envconfig:"PASSWORD" required:"true"`

	// Database SSL mode
	SSLMode string `envconfig:"SSL_MODE" default:"disable"`

	// Connection pool settings
	MaxOpenConns    int `envconfig:"MAX_OPEN_CONNS" default:"25"`
	MaxIdleConns    int `envconfig:"MAX_IDLE_CONNS" default:"5"`
	ConnMaxLifetime int `envconfig:"CONN_MAX_LIFETIME" default:"300"`
}

// LoggingConfig represents logging-specific configuration.
type LoggingConfig struct {
	// Log level (debug, info, warn, error)
	Level string `envconfig:"LEVEL" default:"info"`

	// Log format (json, text)
	Format string `envconfig:"FORMAT" default:"json"`

	// Enable structured logging
	Structured bool `envconfig:"STRUCTURED" default:"true"`

	// Include caller information
	IncludeCaller bool `envconfig:"INCLUDE_CALLER" default:"false"`
}

// Load loads configuration from environment variables.
// The prefix parameter is used to namespace environment variables.
// For example, with prefix "APP", environment variables like APP_SERVER_PORT will be loaded.
func Load(prefix string) (*Config, error) {
	var cfg Config

	// Process environment variables with the given prefix
	err := envconfig.Process(prefix, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	return &cfg, nil
}

// LoadFromFile loads configuration from a file and environment variables.
// Environment variables take precedence over file values.
func LoadFromFile(prefix, filepath string) (*Config, error) {
	// First, load from file if it exists
	if _, err := os.Stat(filepath); err == nil {
		// Set environment variables from file
		envVars, err := loadEnvFile(filepath)
		if err != nil {
			return nil, fmt.Errorf("failed to load environment file: %w", err)
		}

		// Set environment variables
		for key, value := range envVars {
			os.Setenv(key, value)
		}
	}

	// Then load from environment variables
	return Load(prefix)
}

// loadEnvFile loads environment variables from a file.
// Expected format: KEY=value (one per line)
func loadEnvFile(filepath string) (map[string]string, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	envVars := make(map[string]string)
	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		if len(value) >= 2 && (value[0] == '"' && value[len(value)-1] == '"' || value[0] == '\'' && value[len(value)-1] == '\'') {
			value = value[1 : len(value)-1]
		}

		envVars[key] = value
	}

	return envVars, nil
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.Database.Port <= 0 || c.Database.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", c.Database.Port)
	}

	validEnvironments := []string{"development", "staging", "production"}
	valid := false
	for _, env := range validEnvironments {
		if c.Environment == env {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid environment: %s", c.Environment)
	}

	validLogLevels := []string{"debug", "info", "warn", "error"}
	valid = false
	for _, level := range validLogLevels {
		if c.Logging.Level == level {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid log level: %s", c.Logging.Level)
	}

	validLogFormats := []string{"json", "text"}
	valid = false
	for _, format := range validLogFormats {
		if c.Logging.Format == format {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid log format: %s", c.Logging.Format)
	}

	return nil
}

// GetDSN returns the database connection string.
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
}

// IsDevelopment returns true if the environment is development.
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if the environment is production.
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsStaging returns true if the environment is staging.
func (c *Config) IsStaging() bool {
	return c.Environment == "staging"
}
