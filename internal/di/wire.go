//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/pannpers/go-backend-scaffold/internal/adapter/connect"
	"github.com/pannpers/go-backend-scaffold/internal/infrastructure/server"
	"github.com/pannpers/protobuf-scaffold/gen/go/proto/api/v1/v1connect"
)

// InitializeConnectServer creates a new Connect server with all dependencies wired up.
func InitializeConnectServer() (*server.ConnectServer, error) {
	wire.Build(
		server.NewConnectServer,
		connect.NewUserHandler,
		connect.NewPostHandler,
		wire.Bind(new(v1connect.UserServiceHandler), new(*connect.UserHandler)),
		wire.Bind(new(v1connect.PostServiceHandler), new(*connect.PostHandler)),
	)
	return nil, nil
}
