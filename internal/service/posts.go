package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/kokhno-nikolay/news/domain"
	"github.com/kokhno-nikolay/news/internal/repository"
)

type PostService struct {
	repo repository.Posts
}

func NewPostsService(repo repository.Posts) *PostService {
	return &PostService{
		repo: repo,
	}
}

func (s *PostService) Get(ctx context.Context, id int) (*domain.Post, error) {
	if id < 0 {
		return nil, errors.New("<id> must be greater than or equal to zero")
	}

	post, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) List(ctx context.Context, limit int) ([]*domain.Post, error) {
	// ось тут все залежить від бізнес-завдання, звичайно ж у реальному проекті швидше за все використовувалася б пагінація

	posts, err := s.repo.List(ctx, nil)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *PostService) Create(ctx context.Context, input domain.PostInput) (*domain.Post, error) {
	if err := validatePostInput(&input); err != nil {
		return nil, err
	}

	post, err := s.repo.Create(ctx, &input)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) Update(ctx context.Context, id int, input domain.PostInput) (*domain.Post, error) {
	if id < 0 {
		return nil, errors.New("<id> must be greater than or equal to zero")
	}

	if err := validatePostInput(&input); err != nil {
		return nil, err
	}

	post, err := s.repo.Update(ctx, id, &input)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) Delete(ctx context.Context, id int) (bool, error) {
	success, err := s.repo.Delete(ctx, id)
	if err != nil {
		return false, err
	}

	return success, nil
}

// ось тут також все залежить від бізнес завдання, просто додав пару кейсів для тестів
// реальну валідацію обговорював би з продуктами
func validatePostInput(input *domain.PostInput) error {
	if len(input.Title) < 3 {
		return fmt.Errorf("title must be at least 3 characters long")
	}

	if len(input.Title) > 100 {
		return fmt.Errorf("title must be at most 100 characters long")
	}

	if len(input.Content) < 3 {
		return fmt.Errorf("content must be at least 3 characters long")
	}

	if len(input.Content) > 500 {
		return fmt.Errorf("content must be at most 500 characters long")
	}

	return nil
}
