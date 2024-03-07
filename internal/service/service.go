package service

import "github.com/kokhno-nikolay/news/internal/repository"

type Service struct {
	PostService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		PostService: *NewPostsService(repo.Posts),
	}
}
