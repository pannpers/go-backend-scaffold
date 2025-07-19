package rdb_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/pannpers/go-backend-scaffold/internal/infrastructure/database/rdb"
	"github.com/pannpers/go-backend-scaffold/pkg/config"
	"github.com/pannpers/go-backend-scaffold/pkg/logging"
	"github.com/uptrace/bun/extra/bundebug"
)

var testDB *rdb.Database

func TestMain(m *testing.M) {
	testDB = setupTestDatabase()
	testDB.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
	))

	// Run tests first to determine if we're in short mode
	code := m.Run()

	// Clean up database if it was initialized
	if testDB != nil {
		testDB.Close()
	}

	os.Exit(code)
}

func setupTestDatabase() *rdb.Database {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:            "localhost",
			Port:            5432,
			Name:            "scaffold_test",
			User:            "testuser",
			Password:        "testpassword",
			SSLMode:         "disable",
			MaxOpenConns:    10,
			MaxIdleConns:    5,
			ConnMaxLifetime: 300,
		},
	}

	logger := logging.New()
	ctx := context.Background()

	// Create database connection using rdb.New()
	db, err := rdb.New(ctx, cfg, logger)
	if err != nil {
		panic("Failed to connect to test database: " + err.Error())
	}

	// Enable UUID extension
	_, err = db.NewRaw("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Exec(ctx)
	if err != nil {
		panic("Failed to create uuid-ossp extension: " + err.Error())
	}

	// Reset test tables
	err = db.ResetModel(ctx,
		(*rdb.User)(nil),
		(*rdb.Post)(nil),
	)
	if err != nil {
		panic("Failed to reset models: " + err.Error())
	}

	// Add foreign key constraint
	_, err = db.NewRaw("ALTER TABLE posts ADD CONSTRAINT fk_posts_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE").Exec(ctx)
	if err != nil {
		// Ignore error if constraint already exists
		if !strings.Contains(err.Error(), "already exists") {
			panic("Failed to add foreign key constraint: " + err.Error())
		}
	}

	return db
}
