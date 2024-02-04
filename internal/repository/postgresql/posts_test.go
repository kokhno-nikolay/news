package postgresql_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/kokhno-nikolay/news/domain"
	"github.com/kokhno-nikolay/news/internal/repository/postgresql"
)

func TestPostsRepo_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	repo := postgresql.NewPostRepo(sqlx.NewDb(db, "sqlmock"))

	// Input data for testing
	id := 1

	// Expected data after retrieval
	expectedPost := &domain.Post{
		ID:        id,
		Title:     "Test Title",
		Content:   "Test Content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectQuery("SELECT id, title, content, created_at, updated_at FROM posts WHERE id = ?").
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "created_at", "updated_at"}).
			AddRow(
				expectedPost.ID,
				expectedPost.Title,
				expectedPost.Content,
				expectedPost.CreatedAt,
				expectedPost.UpdatedAt,
			),
		)

	retrievedPost, err := repo.Get(context.Background(), id)

	assert.NoError(t, err)
	assert.NotNil(t, retrievedPost)
	assert.Equal(t, expectedPost.ID, retrievedPost.ID)
	assert.Equal(t, expectedPost.Title, retrievedPost.Title)
	assert.Equal(t, expectedPost.Content, retrievedPost.Content)
	assert.WithinDuration(t, expectedPost.CreatedAt, retrievedPost.CreatedAt, time.Second)
	assert.WithinDuration(t, expectedPost.UpdatedAt, retrievedPost.UpdatedAt, time.Second)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestPostRepo_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	repo := postgresql.NewPostRepo(sqlx.NewDb(db, "sqlmock"))

	// Input data for testing
	limit := 5

	// Expected data after retrieval
	expectedPost := []*domain.Post{
		{
			ID:        1,
			Title:     "Test Title 1",
			Content:   "Test Content 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Title:     "Test Title 2",
			Content:   "Test Content 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        3,
			Title:     "Test Title 3",
			Content:   "Test Content 3",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mock.ExpectQuery("SELECT id, title, content, created_at, updated_at FROM posts LIMIT \\$1").
		WithArgs(limit).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "created_at", "updated_at"}).
			AddRow(expectedPost[0].ID, expectedPost[0].Title, expectedPost[0].Content,
				expectedPost[0].CreatedAt, expectedPost[0].UpdatedAt).
			AddRow(expectedPost[1].ID, expectedPost[1].Title, expectedPost[1].Content,
				expectedPost[1].CreatedAt, expectedPost[1].UpdatedAt).
			AddRow(expectedPost[2].ID, expectedPost[2].Title, expectedPost[2].Content,
				expectedPost[2].CreatedAt, expectedPost[2].UpdatedAt))

	postList, err := repo.List(context.Background(), &limit)

	assert.NoError(t, err)
	assert.NotNil(t, postList)
	assert.Len(t, postList, len(expectedPost))

	for i := range expectedPost {
		assert.Equal(t, expectedPost[i].ID, postList[i].ID)
		assert.Equal(t, expectedPost[i].Title, postList[i].Title)
		assert.Equal(t, expectedPost[i].Content, postList[i].Content)
		assert.WithinDuration(t, expectedPost[i].CreatedAt, postList[i].CreatedAt, time.Second)
		assert.WithinDuration(t, expectedPost[i].UpdatedAt, postList[i].UpdatedAt, time.Second)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestPostRepo_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	repo := postgresql.NewPostRepo(sqlx.NewDb(db, "sqlmock"))

	// Input data for testing
	input := &domain.PostInput{
		Title:   "Test Title",
		Content: "Test Content",
	}

	// Expected data after retrieval
	expectedPost := &domain.Post{
		ID:        1,
		Title:     input.Title,
		Content:   input.Content,
		CreatedAt: time.Now(),
	}

	mock.ExpectQuery("INSERT INTO posts").
		WithArgs(input.Title, input.Content).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "created_at"}).
			AddRow(
				expectedPost.ID,
				expectedPost.Title,
				expectedPost.Content,
				expectedPost.CreatedAt,
			),
		)

	createdPost, err := repo.Create(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, createdPost)
	assert.Equal(t, expectedPost.ID, createdPost.ID)
	assert.Equal(t, expectedPost.Title, createdPost.Title)
	assert.Equal(t, expectedPost.Content, createdPost.Content)
	assert.WithinDuration(t, expectedPost.CreatedAt, createdPost.CreatedAt, time.Second)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestPostRepo_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	repo := postgresql.NewPostRepo(sqlx.NewDb(db, "sqlmock"))

	// Input data for testing
	id := 1

	updateInput := &domain.PostInput{
		Title:   "Updated Title",
		Content: "Updated Content",
	}

	// Expected data after update
	expectedPost := &domain.Post{
		ID:        id,
		Title:     updateInput.Title,
		Content:   updateInput.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectQuery("UPDATE posts SET title = \\?, content = \\? WHERE id = \\? RETURNING id, title, content, created_at, updated_at").
		WithArgs(updateInput.Title, updateInput.Content, id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "created_at", "updated_at"}).
			AddRow(expectedPost.ID, expectedPost.Title, expectedPost.Content, expectedPost.CreatedAt, expectedPost.UpdatedAt))

	updatedPost, err := repo.Update(context.Background(), id, updateInput)

	assert.NoError(t, err)
	assert.NotNil(t, updatedPost)
	assert.Equal(t, expectedPost.ID, updatedPost.ID)
	assert.Equal(t, expectedPost.Title, updatedPost.Title)
	assert.Equal(t, expectedPost.Content, updatedPost.Content)
	assert.WithinDuration(t, expectedPost.CreatedAt, updatedPost.CreatedAt, time.Second)
	assert.WithinDuration(t, expectedPost.UpdatedAt, updatedPost.UpdatedAt, time.Second)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestPostRepo_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	repo := postgresql.NewPostRepo(sqlx.NewDb(db, "sqlmock"))

	// Input data for testing
	id := 1

	mock.ExpectExec("^DELETE FROM posts WHERE id = ?").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ok, err := repo.Delete(context.Background(), id)

	assert.NoError(t, err)
	assert.True(t, ok)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}
