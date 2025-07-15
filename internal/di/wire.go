//go:build wireinject
// +build wireinject

package di

import (
	"io"
	"log/slog"

	"github.com/google/wire"
	"github.com/pannpers/go-backend-scaffold/internal/adapter/connect"
	"github.com/pannpers/go-backend-scaffold/internal/infrastructure/database/rdb"
	"github.com/pannpers/go-backend-scaffold/internal/infrastructure/server"
	"github.com/pannpers/go-backend-scaffold/pkg/config"
	"github.com/pannpers/go-backend-scaffold/pkg/logging"
	"github.com/pannpers/protobuf-scaffold/gen/go/proto/api/v1/v1connect"
)

// InitializeApp creates a new App with all dependencies wired up.
func InitializeApp() (*App, error) {
	wire.Build(
		newApp,
		server.NewConnectServer,
		provideDatabase,
		provideConfig,
		provideLogger,
		connect.NewUserHandler,
		connect.NewPostHandler,
		wire.Bind(new(v1connect.UserServiceHandler), new(*connect.UserHandler)),
		wire.Bind(new(v1connect.PostServiceHandler), new(*connect.PostHandler)),
	)
	return nil, nil
}

func newApp(server *server.ConnectServer, db *rdb.Database) *App {
	return &App{
		Server:  server,
		Closers: []io.Closer{db},
	}
}

// provideConfig creates a new config instance.
func provideConfig() (*config.Config, error) {
	return config.Load("")
}

// provideLogger creates a new logger instance based on config.
func provideLogger(cfg *config.Config) *logging.Logger {
	var opts []logging.Option

	// Set log level based on config
	switch cfg.Logging.Level {
	case "debug":
		opts = append(opts, logging.WithLevel(slog.LevelDebug))
	case "info":
		opts = append(opts, logging.WithLevel(slog.LevelInfo))
	case "warn":
		opts = append(opts, logging.WithLevel(slog.LevelWarn))
	case "error":
		opts = append(opts, logging.WithLevel(slog.LevelError))
	}

	// Set log format based on config
	switch cfg.Logging.Format {
	case "text":
		opts = append(opts, logging.WithFormat(logging.FormatText))
	case "json":
		opts = append(opts, logging.WithFormat(logging.FormatJSON))
	}

	return logging.New(opts...)
}

// provideDatabase creates a new database instance.
func provideDatabase(cfg *config.Config, logger *logging.Logger) (*rdb.Database, error) {
	return rdb.New(cfg, logger)
}
