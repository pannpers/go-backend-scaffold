package rdb

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/pannpers/go-backend-scaffold/pkg/config"
	"github.com/pannpers/go-backend-scaffold/pkg/logging"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// Database represents the database instance.
type Database struct {
	db     *bun.DB
	logger *logging.Logger
	cfg    *config.Config
}

// New creates a new database instance with connection and ping verification.
func New(cfg *config.Config, logger *logging.Logger) (*Database, error) {
	// Create PostgreSQL driver
	dsn := cfg.Database.GetDSN()
	driver := pgdriver.NewConnector(pgdriver.WithDSN(dsn))

	// Create sql.DB instance
	sqldb := sql.OpenDB(driver)

	// Create Bun database instance
	db := bun.NewDB(sqldb, pgdialect.New())

	// Set connection pool settings
	sqldb.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqldb.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqldb.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)

	// Create database instance
	database := &Database{
		db:     db,
		logger: logger,
		cfg:    cfg,
	}

	// Verify connection with ping
	if err := database.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info(context.Background(), "Database connection established successfully",
		slog.String("host", cfg.Database.Host),
		slog.Int("port", cfg.Database.Port),
		slog.String("database", cfg.Database.Name),
		slog.Int("max_open_conns", cfg.Database.MaxOpenConns),
		slog.Int("max_idle_conns", cfg.Database.MaxIdleConns),
	)

	return database, nil
}

// Ping verifies the database connection.
func (d *Database) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return d.db.PingContext(ctx)
}

// Close closes the database connection.
func (d *Database) Close() error {
	if d.db != nil {
		d.logger.Info(context.Background(), "Closing database connection")
		return d.db.Close()
	}
	return nil
}

// GetDB returns the underlying Bun database instance.
func (d *Database) GetDB() *bun.DB {
	return d.db
}
