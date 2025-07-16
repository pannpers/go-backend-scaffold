package connect

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"
	"github.com/pannpers/go-backend-scaffold/pkg/logging"
	api "github.com/pannpers/protobuf-scaffold/gen/go/proto/api/v1"
	entity "github.com/pannpers/protobuf-scaffold/gen/go/proto/entity/v1"
)

// PostHandler implements the PostService Connect interface.
type PostHandler struct {
	logger *logging.Logger
	// Add your use case dependencies here.
	// For example:
	// postUseCase usecase.PostUseCase.
}

// NewPostHandler creates a new post handler.
func NewPostHandler(logger *logging.Logger) *PostHandler {
	return &PostHandler{
		logger: logger,
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