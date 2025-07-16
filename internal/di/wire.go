//go:build wireinject
// +build wireinject

package di

import (
	"context"
	"io"
	"log/slog"
	"time"

	"github.com/google/wire"
	"github.com/pannpers/go-backend-scaffold/internal/adapter/connect"
	"github.com/pannpers/go-backend-scaffold/internal/entity"
	"github.com/pannpers/go-backend-scaffold/internal/infrastructure/database/rdb"
	"github.com/pannpers/go-backend-scaffold/internal/infrastructure/server"
	"github.com/pannpers/go-backend-scaffold/internal/usecase"
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
		
		// Repository layer (temporary mock implementations)
		provideMockUserRepository,
		provideMockPostRepository,
		
		// Use case layer
		usecase.NewUserUseCase,
		usecase.NewPostUseCase,
		
		// Handler layer
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

// Mock implementations for development/testing
// TODO: Replace with actual database implementations

// MockUserRepository is a simple mock implementation for development
type MockUserRepository struct{}

func (m *MockUserRepository) Create(ctx context.Context, params *entity.NewUser) (*entity.User, error) {
	return &entity.User{
		ID:        "mock-user-id",
		Name:      params.Name,
		Email:     params.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *MockUserRepository) Get(ctx context.Context, id string) (*entity.User, error) {
	return &entity.User{
		ID:        id,
		Name:      "Mock User",
		Email:     "mock@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	return nil
}

// MockPostRepository is a simple mock implementation for development
type MockPostRepository struct{}

func (m *MockPostRepository) Create(ctx context.Context, params *entity.NewPost) (*entity.Post, error) {
	return &entity.Post{
		ID:        "mock-post-id",
		Title:     params.Title,
		UserID:    params.UserID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *MockPostRepository) Get(ctx context.Context, id string) (*entity.Post, error) {
	return &entity.Post{
		ID:        id,
		Title:     "Mock Post",
		UserID:    "mock-user-id",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *MockPostRepository) Delete(ctx context.Context, id string) error {
	return nil
}

// provideMockUserRepository creates a mock user repository implementation.
// TODO: Replace with actual database implementation
func provideMockUserRepository() entity.UserRepository {
	return &MockUserRepository{}
}

// provideMockPostRepository creates a mock post repository implementation.
// TODO: Replace with actual database implementation
func provideMockPostRepository() entity.PostRepository {
	return &MockPostRepository{}
}
