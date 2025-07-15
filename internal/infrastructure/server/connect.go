package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"log/slog"

	"github.com/pannpers/go-backend-scaffold/pkg/config"
	"github.com/pannpers/go-backend-scaffold/pkg/logging"
	"github.com/pannpers/protobuf-scaffold/gen/go/proto/api/v1/v1connect"
)

// ConnectServer represents the Connect server.
type ConnectServer struct {
	server  *http.Server
	logger  *logging.Logger
	Cfg     *config.Config
	address string
}

// NewConnectServer creates a new Connect server instance.
func NewConnectServer(
	cfg *config.Config,
	logger *logging.Logger,
	userHandler v1connect.UserServiceHandler,
	postHandler v1connect.PostServiceHandler,
) *ConnectServer {
	mux := http.NewServeMux()

	// Register Connect handlers.
	path, handler := v1connect.NewUserServiceHandler(userHandler)
	mux.Handle(path, handler)

	path, handler = v1connect.NewPostServiceHandler(postHandler)
	mux.Handle(path, handler)

	address := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	server := &http.Server{
		Addr:         address,
		Handler:      mux,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	return &ConnectServer{
		server:  server,
		logger:  logger,
		Cfg:     cfg,
		address: address,
	}
}

// Start starts the Connect server.
func (s *ConnectServer) Start() error {
	s.logger.Info(context.Background(), fmt.Sprintf("Connect Server starting on %s", s.address))
	return s.server.ListenAndServe()
}

// Stop gracefully stops the Connect server.
func (s *ConnectServer) Stop() error {
	if s.server != nil {
		timeout := s.Cfg.Server.ShutdownTimeout
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()

		s.logger.Info(ctx, "Shutting down Connect server gracefully...", slog.Int("timeout_sec", timeout))

		return s.server.Shutdown(ctx)
	}

	return nil
}
