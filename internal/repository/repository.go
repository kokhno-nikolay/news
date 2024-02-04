package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/kokhno-nikolay/news/domain"
	"github.com/kokhno-nikolay/news/internal/repository/postgresql"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type Posts interface {
	Get(ctx context.Context, id int) (*domain.Post, error)
	List(ctx context.Context, limit *int) ([]*domain.Post, error)
	Create(ctx context.Context, input *domain.PostInput) (*domain.Post, error)
	Update(ctx context.Context, id int, input *domain.PostInput) (*domain.Post, error)
	Delete(ctx context.Context, id int) (bool, error)
}

type Repository struct {
	Posts
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Posts: postgresql.NewPostRepo(db),
	}
}
