# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a modern Go backend scaffold implementing clean architecture with gRPC/Connect-RPC, structured logging, and dependency injection. The project uses protobuf definitions from `github.com/pannpers/protobuf-scaffold` and follows domain-driven design principles.

## Development Commands

### Running the Application
```bash
# Start the server (HTTP on :8080, gRPC/Connect on :9090)
go run cmd/api/main.go
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./pkg/config
go test ./pkg/logging
go test ./pkg/apperr
```

### Code Generation
```bash
# Generate Wire dependency injection code (when wire.go is modified)
wire internal/di/

# Generate protobuf code (if working with proto files)
buf generate
```

### API Testing with buf curl
```bash
# Test GetUser endpoint
buf curl --schema ../protobuf-scaffold --protocol connect \
  -d '{"user_id": "123"}' \
  http://localhost:9090/api.UserService/GetUser

# Test CreateUser endpoint  
buf curl --schema ../protobuf-scaffold --protocol connect \
  -d '{"user": {"id": {"value": "test123"}, "name": {"value": "John Doe"}, "email": {"value": "john@example.com"}}}' \
  http://localhost:9090/api.UserService/CreateUser

# Test Health Check endpoint
buf curl --schema buf.build/grpc/health --protocol connect \
  -d '{"service": ""}' \
  http://localhost:9090/grpc.health.v1.Health/Check
```

## Architecture

### Clean Architecture Layers
- **cmd/**: Application entry points (main.go)
- **internal/adapter/**: Interface adapters (Connect-RPC handlers)
- **internal/di/**: Dependency injection with Google Wire
- **internal/entity/**: Domain entities and business objects
- **internal/infrastructure/**: Infrastructure concerns (servers, databases)
- **internal/usecase/**: Business logic and use cases
- **pkg/**: Reusable packages (config, logging, apperr)

### Key Dependencies
- **Connect-RPC**: `github.com/bufbuild/connect-go` for HTTP/gRPC-compatible APIs
- **Wire**: `github.com/google/wire` for compile-time dependency injection  
- **Protobuf**: Uses `github.com/pannpers/protobuf-scaffold` for shared definitions
- **Database**: Bun ORM with PostgreSQL support via `github.com/uptrace/bun`
- **Logging**: Custom structured logging with OpenTelemetry integration
- **Health Checks**: `connectrpc.com/grpchealth` for gRPC-compatible health monitoring

### Configuration Management
The project uses environment variables for configuration with prefix support:
- Default prefix: `APP_` (e.g., `APP_SERVER_PORT=8080`)
- Configuration is managed in `pkg/config/` with comprehensive validation
- Supports .env files and runtime environment variables
- See `pkg/config/README.md` for detailed configuration options

### Dependency Injection
- Uses Google Wire for compile-time DI
- Main DI configuration in `internal/di/wire.go`
- Generated code in `internal/di/wire_gen.go` (regenerate with `wire internal/di/`)
- App initialization creates server and manages resource lifecycle

### Error Handling
- Custom error package `pkg/apperr/` provides structured error handling
- Includes error codes, HTTP status mapping, and context preservation
- Use `apperr` for consistent error responses across the application

### Logging
- Custom logging package `pkg/logging/` with OpenTelemetry integration
- Supports both JSON and text formats
- Automatic trace_id and span_id injection when using context
- Configurable log levels (debug, info, warn, error)

## Service Implementation

### Connect-RPC Handlers
Handlers are in `internal/adapter/connect/` and implement the generated service interfaces:
- **User Service**: `user_handler.go` - User management endpoints (`/api.UserService/`)
- **Post Service**: `post_handler.go` - Post management endpoints (`/api.PostService/`)
- **Health Check**: `health_handler.go` - Database connectivity health checks (`/grpc.health.v1.Health/`)
- Use Connect protocol, not plain gRPC
- Handlers are bound to interfaces via Wire in `internal/di/wire.go`

### Health Monitoring
- gRPC-compatible health check endpoint at `/grpc.health.v1.Health/Check`
- Verifies database connectivity by pinging the PostgreSQL connection
- Returns `SERVING` when healthy, `NOT_SERVING` when database is unreachable
- Compatible with Kubernetes liveness/readiness probes and load balancers
- Structured logging of health check results with service context

### Database Integration
- Uses Bun ORM with PostgreSQL driver
- Database configuration via environment variables (see config package)
- Connection management handled in `internal/infrastructure/database/rdb/`

## Development Notes

- The project follows Go module conventions with `github.com/pannpers/go-backend-scaffold` as module name
- Wire dependency injection requires regeneration when `wire.go` is modified
- Connect-RPC handlers use HTTP/1.1 compatible protocol (no need for HTTP/2)
- Configuration supports multiple environments (development, staging, production)
- Graceful shutdown is implemented in main.go with proper resource cleanup