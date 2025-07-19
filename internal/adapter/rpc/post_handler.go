package rpc

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/pannpers/go-backend-scaffold/internal/adapter/rpc/mapper"
	"github.com/pannpers/go-backend-scaffold/internal/usecase"
	"github.com/pannpers/go-backend-scaffold/pkg/logging"
	api "github.com/pannpers/protobuf-scaffold/gen/go/proto/api/v1"
)

// PostHandler implements the PostService Connect interface.
type PostHandler struct {
	postUseCase *usecase.PostUseCase
	logger      *logging.Logger
}

// NewPostHandler creates a new post handler.
func NewPostHandler(postUseCase *usecase.PostUseCase, logger *logging.Logger) *PostHandler {
	return &PostHandler{
		postUseCase: postUseCase,
		logger:      logger,
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

	// Use the use case layer for business logic
	post, err := h.postUseCase.GetPost(ctx, req.Msg.PostId)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&api.GetPostResponse{
		Post: mapper.PostToProto(post),
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

	// TODO: Extract user ID from context/authentication
	userID := "default-user-id"

	// Convert protobuf to domain DTO
	newPost := mapper.NewPostFromProto(req.Msg.Post, userID)

	// Use the use case layer for business logic
	createdPost, err := h.postUseCase.CreatePost(ctx, newPost)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&api.CreatePostResponse{
		Post: mapper.PostToProto(createdPost),
	}), nil
}
