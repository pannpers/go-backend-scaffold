package connect

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"
	api "github.com/pannpers/protobuf-scaffold/gen/go/proto/api/v1"
	entity "github.com/pannpers/protobuf-scaffold/gen/go/proto/entity/v1"
)

// UserHandler implements the UserService Connect interface.
type UserHandler struct {
	// Add your use case dependencies here.
	// For example:
	// userUseCase usecase.UserUseCase.
}

// NewUserHandler creates a new user handler.
func NewUserHandler() *UserHandler {
	return &UserHandler{
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

// PostHandler implements the PostService Connect interface.
type PostHandler struct {
	// Add your use case dependencies here.
	// For example:
	// postUseCase usecase.PostUseCase.
}

// NewPostHandler creates a new post handler.
func NewPostHandler() *PostHandler {
	return &PostHandler{
		// Initialize your dependencies here.
	}
}

// GetPost retrieves a post by ID.
func (h *PostHandler) GetPost(ctx context.Context, req *connect.Request[api.GetPostRequest]) (*connect.Response[api.GetPostResponse], error) {
	if req == nil || req.Msg == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("request cannot be nil"))
	}

	if req.Msg.PostId == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("post_id is required"))
	}

	// TODO: Implement actual business logic using use case layer.
	// For now, return a mock post.
	post := &entity.Post{
		Id: &entity.PostId{
			Value: req.Msg.PostId,
		},
		Title: &entity.PostTitle{
			Value: "Example Post Title",
		},
	}

	return connect.NewResponse(&api.GetPostResponse{
		Post: post,
	}), nil
}

// CreatePost creates a new post.
func (h *PostHandler) CreatePost(ctx context.Context, req *connect.Request[api.CreatePostRequest]) (*connect.Response[api.CreatePostResponse], error) {
	if req == nil || req.Msg == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("request cannot be nil"))
	}

	if req.Msg.Post == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("post is required"))
	}

	// TODO: Implement actual business logic using use case layer.
	// For now, return the same post with a generated ID.
	createdPost := &entity.Post{
		Id: &entity.PostId{
			Value: "generated-post-id-" + req.Msg.Post.GetId().GetValue(),
		},
		Title: req.Msg.Post.Title,
	}

	return connect.NewResponse(&api.CreatePostResponse{
		Post: createdPost,
	}), nil
}
