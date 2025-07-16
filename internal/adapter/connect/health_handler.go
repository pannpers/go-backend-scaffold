package connect

import (
	"context"
	"log/slog"

	"connectrpc.com/grpchealth"
	"github.com/pannpers/go-backend-scaffold/internal/infrastructure/database/rdb"
	"github.com/pannpers/go-backend-scaffold/pkg/logging"
)

// HealthCheckHandler implements grpchealth.Checker interface with database ping.
type HealthCheckHandler struct {
	db     *rdb.Database
	logger *logging.Logger
}

// NewHealthCheckHandler creates a new health check handler.
func NewHealthCheckHandler(db *rdb.Database, logger *logging.Logger) *HealthCheckHandler {
	return &HealthCheckHandler{
		db:     db,
		logger: logger,
	}
}

// Check implements the grpchealth.Checker interface.
func (h *HealthCheckHandler) Check(ctx context.Context, req *grpchealth.CheckRequest) (*grpchealth.CheckResponse, error) {
	service := req.Service

	// For service-specific checks, you can add logic here
	// For now, we'll check the database connection for all services

	if err := h.db.Ping(); err != nil {
		h.logger.Error(ctx, "Health check failed: database ping failed", err, slog.String("service", service))

		return &grpchealth.CheckResponse{Status: grpchealth.StatusNotServing}, nil
	}

	h.logger.Debug(ctx, "Health check passed", slog.String("service", service))

	return &grpchealth.CheckResponse{Status: grpchealth.StatusServing}, nil
}
