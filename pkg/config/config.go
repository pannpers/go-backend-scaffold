// Package config provides application configuration management using environment variables.
// It uses github.com/kelseyhightower/envconfig for loading configuration from environment variables
// with support for validation, default values, and environment-specific helpers.
//
// # Basic Usage
//
// Load configuration from environment variables:
//
//	cfg, err := config.Load("APP")
//	if err != nil {
//		log.Fatalf("Failed to load configuration: %v", err)
//	}
//
//	// Validate configuration
//	if err := cfg.Validate(); err != nil {
//		log.Fatalf("Invalid configuration: %v", err)
//	}
//
// # Environment Variables
//
// The following environment variables are supported (using "APP" prefix):
//
// Basic configuration:
//   - APP_ENVIRONMENT: Environment (development, staging, production)
//   - APP_DEBUG: Debug mode (true/false)
//
// Server configuration:
//   - APP_SERVER_PORT: Server port (default: 8080)
//   - APP_SERVER_HOST: Server host (default: localhost)
//   - APP_SERVER_READ_TIMEOUT: Read timeout in seconds (default: 30)
//   - APP_SERVER_WRITE_TIMEOUT: Write timeout in seconds (default: 30)
//   - APP_SERVER_IDLE_TIMEOUT: Idle timeout in seconds (default: 60)
//   - APP_SERVER_SHUTDOWN_TIMEOUT: Shutdown timeout in seconds (default: 30)
//
// Database configuration:
//   - APP_DATABASE_HOST: Database host (default: localhost)
//   - APP_DATABASE_PORT: Database port (default: 5432)
//   - APP_DATABASE_NAME: Database name (required)
//   - APP_DATABASE_USER: Database user (required)
//   - APP_DATABASE_PASSWORD: Database password (required)
//   - APP_DATABASE_SSL_MODE: SSL mode (default: disable)
//   - APP_DATABASE_MAX_OPEN_CONNS: Maximum open connections (default: 25)
//   - APP_DATABASE_MAX_IDLE_CONNS: Maximum idle connections (default: 5)
//   - APP_DATABASE_CONN_MAX_LIFETIME: Connection max lifetime in seconds (default: 300)
//
// Logging configuration:
//   - APP_LOGGING_LEVEL: Log level (debug, info, warn, error, default: info)
//   - APP_LOGGING_FORMAT: Log format (json, text, default: json)
//   - APP_LOGGING_STRUCTURED: Enable structured logging (default: true)
//   - APP_LOGGING_INCLUDE_CALLER: Include caller information (default: false)
//
// # Environment Helpers
//
// Use environment detection helpers:
//
//	if cfg.IsDevelopment() {
//		// Development-specific logic
//	}
//
//	if cfg.IsProduction() {
//		// Production-specific logic
//	}
//
// # Database Connection
//
// Get database connection string:
//
//	dsn := cfg.Database.GetDSN()
//	// Returns: "postgres://user:pass@host:port/dbname?sslmode=disable"
package config

import (
	"fmt"

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
//
// Example:
//
//	cfg, err := config.Load("APP")
//	if err != nil {
//		return fmt.Errorf("failed to load config: %w", err)
//	}
func Load(prefix string) (*Config, error) {
	var cfg Config

	// Process environment variables with the given prefix
	err := envconfig.Process(prefix, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	return &cfg, nil
}

// Validate validates the configuration according to the following rules:
//   - Server port: 1-65535 range
//   - Database port: 1-65535 range
//   - Environment: development, staging, or production
//   - Log level: debug, info, warn, or error
//   - Log format: json or text
//   - Required fields: Database name, user, and password
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

// GetDSN returns the PostgreSQL database connection string in the format:
// "postgres://user:password@host:port/dbname?sslmode=mode"
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
}

// IsDevelopment returns true if the environment is "development".
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if the environment is "production".
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsStaging returns true if the environment is "staging".
func (c *Config) IsStaging() bool {
	return c.Environment == "staging"
}
