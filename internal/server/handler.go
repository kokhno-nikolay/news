package server

import (
	"context"

	proto "github.com/kokhno-nikolay/news/api/proto"
	"github.com/kokhno-nikolay/news/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// @Summary		Get post
// @Description	Getting post entity by id.
// @Tags		posts
// @Accept		json
// @Produce		json
// @Param		id           path       int true  "Post ID"
// @Success		200          {object}   domain.Post
// @Failure		400,404      {object}   errorResponse
// @Failure		500          {object}   errorResponse
// @Failure		default      {object}   errorResponse
// @Router		/posts/{id}  [get]
func (s *Server) Get(ctx context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
	post, err := s.postService.Get(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}

	return &proto.GetResponse{
		Post: convertPostToProto(post),
	}, nil
}

// @Summary		Get post list
// @Description	Getting post list
// @Tags		posts
// @Accept		json
// @Produce		json
// @Success		200      {object}   domain.Post
// @Failure		400,404  {object}   errorResponse
// @Failure		500      {object}   errorResponse
// @Failure		default  {object}   errorResponse
// @Router		/posts [get]
func (s *Server) List(ctx context.Context, req *proto.ListRequest) (*proto.ListResponse, error) {
	res, err := s.postService.List(ctx, int(req.Limit))
	if err != nil {
		return nil, err
	}

	posts := make([]*proto.Post, 0)
	for _, item := range res {
		posts = append(posts, convertPostToProto(item))
	}

	return &proto.ListResponse{
		Posts: posts,
	}, nil
}

// @Summary		Create post
// @Description	Creates a new post entity
// @Tags		posts
// @Accept		json
// @Produce		json
// @Param		input   body        domain.PostInput	true	"Post content"
// @Success		201     {object}    domain.Post
// @Failure		400,404 {object}    errorResponse
// @Failure		500     {object}    errorResponse
// @Failure		default {object}    errorResponse
// @Router		/posts [post]
func (s *Server) Create(ctx context.Context, req *proto.CreateRequest) (*proto.Post, error) {
	res, err := s.postService.Create(ctx, domain.PostInput{
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		return nil, err
	}

	return convertPostToProto(res), nil
}

// @Summary		Update post
// @Description	Updatting post entity
// @Tags		posts
// @Accept		json
// @Produce		json
// @Param		id      path        int                 true    "Post ID"
// @Param		input   body        domain.PostInput    true    "Post content"
// @Success		201     {object}    domain.Post
// @Failure		400,404 {object}    errorResponse
// @Failure		500     {object}    errorResponse
// @Failure		default {object}    errorResponse
// @Router		/posts/{id} [put]
func (s *Server) Update(ctx context.Context, req *proto.UpdateRequest) (*proto.Post, error) {
	res, err := s.postService.Update(ctx, int(req.Id), domain.PostInput{
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		return nil, err
	}

	return convertPostToProto(res), nil
}

// @Summary     Delete post
// @Description Deletting post entity
// @Tags		posts
// @Accept		json
// @Produce		json
// @Param		id      path        int true   "Post ID"
// @Success		200     {bool}      true
// @Failure		400,404 {object}    errorResponse
// @Failure		500     {object}    errorResponse
// @Failure		default {object}    errorResponse
// @Router		/posts/{id}  [delete]
func (s *Server) Delete(ctx context.Context, req *proto.DeleteRequest) (*proto.DeleteResponse, error) {
	success, err := s.postService.Delete(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}

	return &proto.DeleteResponse{
		Success: success,
	}, nil
}

func convertPostToProto(post *domain.Post) *proto.Post {
	return &proto.Post{
		Id:        int64(post.ID),
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: timestamppb.New(post.CreatedAt),
		UpdatedAt: timestamppb.New(post.UpdatedAt),
	}
}
