package connect

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"
	"github.com/pannpers/go-backend-scaffold/pkg/logging"
	api "github.com/pannpers/protobuf-scaffold/gen/go/proto/api/v1"
	entity "github.com/pannpers/protobuf-scaffold/gen/go/proto/entity/v1"
)

// UserHandler implements the UserService Connect interface.
type UserHandler struct {
	logger *logging.Logger
	// Add your use case dependencies here.
	// For example:
	// userUseCase usecase.UserUseCase.
}

// NewUserHandler creates a new user handler.
func NewUserHandler(logger *logging.Logger) *UserHandler {
	return &UserHandler{
		logger: logger,
		// Initialize your dependencies here.
	}
}

// GetUser retrieves a user by ID.
func (h *UserHandler) GetUser(ctx context.Context, req *connect.Request[api.GetUserRequest]) (*connect.Response[api.GetUserResponse], error) {
	if req == nil || req.Msg == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("request cannot be nil"))
	}

	if req.Msg.UserId == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("user_id is required"))
	}

	// TODO: Implement actual business logic using use case layer.
	// For now, return a mock user.
	user := &entity.User{
		Id: &entity.UserId{
			Value: req.Msg.UserId,
		},
		Name: &entity.UserName{
			Value: "Example User",
		},
		Email: &entity.UserEmail{
			Value: "example@example.com",
		},
	}

	return connect.NewResponse(&api.GetUserResponse{
		User: user,
	}), nil
}

// CreateUser creates a new user.
func (h *UserHandler) CreateUser(ctx context.Context, req *connect.Request[api.CreateUserRequest]) (*connect.Response[api.CreateUserResponse], error) {
	if req == nil || req.Msg == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("request cannot be nil"))
	}

	if req.Msg.User == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("user is required"))
	}

	// TODO: Implement actual business logic using use case layer.
	// For now, return the same user with a generated ID.
	createdUser := &entity.User{
		Id: &entity.UserId{
			Value: "generated-user-id-" + req.Msg.User.GetId().GetValue(),
		},
		Name:  req.Msg.User.Name,
		Email: req.Msg.User.Email,
	}

	return connect.NewResponse(&api.CreateUserResponse{
		User: createdUser,
	}), nil
}