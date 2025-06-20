# Go Backend Scaffold

A modern Go backend scaffold following clean architecture principles with gRPC support, structured logging, and dependency injection.

## Features

- **gRPC Server**: Built-in gRPC server with protobuf support
- **Structured Logging**: Advanced logger with OpenTelemetry integration for distributed tracing
- **Dependency Injection**: Wire-based dependency injection for clean architecture
- **Clean Architecture**: Well-organized project structure following domain-driven design
- **Protobuf Integration**: Ready-to-use with `github.com/pannpers/protobuf-scaffold`

## Project Structure

```
.
├── cmd/
│   └── api/                    # Main application entry point
├── internal/
│   ├── adapter/               # External interface adapters
│   │   └── grpc/             # gRPC handlers and services
│   ├── di/                   # Dependency injection configuration
│   ├── entitiy/              # Domain entities
│   ├── infrastructure/       # Infrastructure concerns
│   │   └── server/           # HTTP and gRPC server implementations
│   └── usecase/              # Business logic and use cases
├── pkg/
│   └── logger/               # Structured logging with OpenTelemetry
└── go.mod                    # Go module definition
```

## Prerequisites

- Go 1.24 or later
- Protocol Buffers compiler (for gRPC development)

## Getting Started

### Installation

1. Clone the repository:

```bash
git clone https://github.com/pannpers/go-backend-scaffold.git
cd go-backend-scaffold
```

2. Install dependencies:

```bash
go mod download
```

### Running the Application

#### HTTP Server

```bash
go run cmd/api/main.go
```

The HTTP server will start on port 8080.

#### gRPC Server

The gRPC server is configured to run on port 9090 with reflection enabled for development tools like `grpcurl`.

### Development

#### Logger Usage

The scaffold includes a powerful structured logger with OpenTelemetry integration:

```go
import "github.com/pannpers/go-backend-scaffold/pkg/logger"

// Create a logger with default options (JSON format)
logger := logger.New()

// Create a logger with custom options
logger := logger.New(
    logger.WithLevel(slog.LevelDebug),
    logger.WithFormat(logger.FormatText), // Human-readable format
    logger.WithWriter(os.Stderr),
)

// Log with context (automatically includes trace_id and span_id)
ctx := context.Background()
logger.Info(ctx, "User logged in", slog.String("user_id", "123"))
```

#### gRPC Handler Implementation

The scaffold provides a foundation for implementing gRPC services:

```go
// Example handler implementation
type UserHandler struct {
    api.UnimplementedUserServiceServer
}

func (h *UserHandler) GetUser(ctx context.Context, req *api.GetUserRequest) (*api.GetUserResponse, error) {
    // Your business logic here
    return &api.GetUserResponse{
        User: &entity.User{
            Id: &entity.UserId{Value: req.UserId},
            Name: &entity.UserName{Value: "Example User"},
        },
    }, nil
}
```

#### Dependency Injection

The project uses Google Wire for dependency injection:

```go
// Initialize HTTP server
server, err := di.InitializeAPI()

// Initialize gRPC server
grpcServer, err := di.InitializeGRPCServer()
```

## Architecture

This scaffold follows clean architecture principles:

- **Entities**: Core business objects (`internal/entity`)
- **Use Cases**: Business logic and rules (`internal/usecase`)
- **Adapters**: External interface implementations (`internal/adapter`)
- **Infrastructure**: Technical concerns like servers and databases (`internal/infrastructure`)

## Dependencies

- **gRPC**: `google.golang.org/grpc` for RPC communication
- **Wire**: `github.com/google/wire` for dependency injection
- **OpenTelemetry**: `go.opentelemetry.io/otel` for distributed tracing
- **Protobuf Scaffold**: `github.com/pannpers/protobuf-scaffold` for shared protobuf definitions

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is licensed under the MIT License.
