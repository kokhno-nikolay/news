package postgresql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/kokhno-nikolay/news/domain"
)

const defaultLimit = 200

type PostRepo struct {
	db *sqlx.DB
}

func NewPostRepo(db *sqlx.DB) *PostRepo {
	return &PostRepo{
		db: db,
	}
}

func (r *PostRepo) Get(ctx context.Context, id int) (*domain.Post, error) {
	query := `
		SELECT id, title, content, created_at, updated_at 
		FROM posts
		WHERE id = $1
	`

	row := r.db.QueryRow(query, &id)

	var post domain.Post
	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(sql.ErrNoRows.Error())
		}

		return nil, err
	}

	return &post, nil
}

func (r *PostRepo) List(ctx context.Context, limit *int) ([]*domain.Post, error) {
	queryLimit := getQueryLimit(limit)

	query := `
		SELECT id, title, content, created_at, updated_at
		FROM posts
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, queryLimit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var postList []*domain.Post

	for rows.Next() {
		var post domain.Post
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		postList = append(postList, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return postList, nil
}

func (r *PostRepo) Create(ctx context.Context, input *domain.PostInput) (*domain.Post, error) {
	query := `
        INSERT INTO posts (title, content) 
        VALUES ($1, $2) 
		RETURNING id, title, content, created_at
    `

	var post domain.Post
	err := r.db.QueryRowContext(ctx, query, input.Title, input.Content).
		Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepo) Update(ctx context.Context, id int, input *domain.PostInput) (*domain.Post, error) {
	query := `
		UPDATE posts
		SET title = $1, content = $2
		WHERE id = $3
		RETURNING id, title, content, created_at, updated_at
	`

	var updatedPost domain.Post
	err := r.db.QueryRowContext(ctx, query, input.Title, input.Content, id).
		Scan(&updatedPost.ID, &updatedPost.Title, &updatedPost.Content, &updatedPost.CreatedAt, &updatedPost.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &updatedPost, nil
}

func (r *PostRepo) Delete(ctx context.Context, id int) (bool, error) {
	query := `
		DELETE FROM posts 
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func getQueryLimit(limit *int) int {
	if limit != nil && *limit <= defaultLimit {
		return *limit
	}

	return defaultLimit
}
