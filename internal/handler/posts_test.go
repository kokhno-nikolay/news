package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"

	"github.com/kokhno-nikolay/news/domain"
	"github.com/kokhno-nikolay/news/internal/handler"
	"github.com/kokhno-nikolay/news/internal/repository"
	mock_repository "github.com/kokhno-nikolay/news/internal/repository/mocks"
)

func TestHandler_Get(t *testing.T) {
	tests := []struct {
		name                 string
		paramID              string
		expectedStatusCode   int
		expectedResponseBody *domain.Post
		errorMessage         string
	}{
		{
			name:               "ok",
			paramID:            "1",
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: &domain.Post{
				ID:        1,
				Title:     "Test title",
				Content:   "Test content",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			name:                 "invalid parameter format",
			paramID:              "qwerty",
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: nil,
			errorMessage:         `{"message":"invalid parameter format"}`,
		},
		{
			name:                 "id is less than zero",
			paramID:              "-11",
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: nil,
			errorMessage:         `{"message":"ID must be greater than or equal to zero"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			postRepo := mock_repository.NewMockPosts(c)
			handlers := handler.NewHandler(&repository.Repository{Posts: postRepo})

			postRepo.EXPECT().Get(gomock.Any(), 1).Return(test.expectedResponseBody, nil).AnyTimes()

			// Init Endpoint
			r := gin.New()
			r.GET("/posts/:id", handlers.Get)

			// Create Request
			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", fmt.Sprintf("/posts/%s", test.paramID), nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)

			if test.errorMessage != "" {
				assert.Equal(t, w.Body.String(), test.errorMessage)
			}

			if test.expectedResponseBody != nil {
				var responseStruct domain.Post
				err := json.Unmarshal(w.Body.Bytes(), &responseStruct)
				if err != nil {
					t.Error(err)
				}

				assert.Equal(t, test.expectedResponseBody.ID, responseStruct.ID)
				assert.Equal(t, test.expectedResponseBody.Title, responseStruct.Title)
				assert.Equal(t, test.expectedResponseBody.Content, responseStruct.Content)
			}
		})
	}
}

func TestHandler_Create(t *testing.T) {
	tests := []struct {
		name                 string
		inputBody            *domain.PostInput
		expectedStatusCode   int
		expectedResponseBody *domain.Post
		errorMessage         string
	}{
		{
			name: "ok",
			inputBody: &domain.PostInput{
				Title:   "Test title",
				Content: "Test content",
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponseBody: &domain.Post{
				ID:        1,
				Title:     "Test title",
				Content:   "Test content",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			name: "invalid title",
			inputBody: &domain.PostInput{
				Title:   "",
				Content: "Test content",
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: nil,
			errorMessage:         `{"message":"title must be at least 3 characters long"}`,
		},
		{
			name: "invalid content",
			inputBody: &domain.PostInput{
				Title:   "Test title",
				Content: "",
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: nil,
			errorMessage:         `{"message":"content must be at least 3 characters long"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			postRepo := mock_repository.NewMockPosts(c)
			handlers := handler.NewHandler(&repository.Repository{Posts: postRepo})

			postRepo.EXPECT().Create(gomock.Any(), test.inputBody).Return(test.expectedResponseBody, nil).AnyTimes()

			// Init Endpoint
			r := gin.New()
			r.POST("/news", handlers.Create)

			// Create Request
			w := httptest.NewRecorder()

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(test.inputBody)
			if err != nil {
				t.Error(err)
			}
			req := httptest.NewRequest("POST", "/news", &buf)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)

			if test.errorMessage != "" {
				assert.Equal(t, w.Body.String(), test.errorMessage)
			}

			if test.expectedResponseBody != nil {
				var responseStruct domain.Post
				err := json.Unmarshal(w.Body.Bytes(), &responseStruct)
				if err != nil {
					t.Error(err)
				}

				assert.Equal(t, test.expectedResponseBody.ID, responseStruct.ID)
				assert.Equal(t, test.expectedResponseBody.Title, responseStruct.Title)
				assert.Equal(t, test.expectedResponseBody.Content, responseStruct.Content)
			}
		})
	}
}

func TestHandler_List(t *testing.T) {
	tests := []struct {
		name                 string
		expectedResponseBody []*domain.Post
		expectedStatusCode   int
		errorMessage         string
	}{
		{
			name: "ok",
			expectedResponseBody: []*domain.Post{
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
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			postRepo := mock_repository.NewMockPosts(c)
			handlers := handler.NewHandler(&repository.Repository{Posts: postRepo})

			postRepo.EXPECT().List(gomock.Any(), nil).Return(test.expectedResponseBody, nil).AnyTimes()

			// Init Endpoint
			r := gin.New()
			r.GET("/posts", handlers.List)

			// Create Request
			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", "/posts", nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)

			var responseStruct []*domain.Post
			err := json.Unmarshal(w.Body.Bytes(), &responseStruct)
			if err != nil {
				t.Error(err)
			}

			for i, post := range test.expectedResponseBody {
				assert.Equal(t, post.ID, responseStruct[i].ID)
				assert.Equal(t, post.Title, responseStruct[i].Title)
				assert.Equal(t, post.Content, responseStruct[i].Content)
			}
		})
	}
}

func TestHandler_Update(t *testing.T) {
	tests := []struct {
		name                 string
		inputBody            *domain.PostInput
		paramID              string
		expectedStatusCode   int
		expectedResponseBody *domain.Post
		errorMessage         string
	}{
		{
			name: "ok",
			inputBody: &domain.PostInput{
				Title:   "Updated title",
				Content: "Updated content",
			},
			paramID:            "1",
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: &domain.Post{
				ID:        1,
				Title:     "Updated title",
				Content:   "Updated content",
				UpdatedAt: time.Now(),
			},
		},
		{
			name:    "invalid parameter format",
			paramID: "qwerty",
			inputBody: &domain.PostInput{
				Title:   "Updated title",
				Content: "Updated content",
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: nil,
			errorMessage:         `{"message":"invalid parameter format"}`,
		},
		{
			name:    "id is less than zero",
			paramID: "-11",
			inputBody: &domain.PostInput{
				Title:   "Updated title",
				Content: "Updated content",
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: nil,
			errorMessage:         `{"message":"ID must be greater than or equal to zero"}`,
		},
		{
			name:    "invalid title",
			paramID: "1",
			inputBody: &domain.PostInput{
				Title:   "",
				Content: "Test content",
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: nil,
			errorMessage:         `{"message":"title must be at least 3 characters long"}`,
		},
		{
			name:    "invalid content",
			paramID: "1",
			inputBody: &domain.PostInput{
				Title:   "Test title",
				Content: "",
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: nil,
			errorMessage:         `{"message":"content must be at least 3 characters long"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			postRepo := mock_repository.NewMockPosts(c)
			handlers := handler.NewHandler(&repository.Repository{Posts: postRepo})

			postRepo.EXPECT().Update(gomock.Any(), 1, test.inputBody).Return(test.expectedResponseBody, nil).AnyTimes()

			// Init Endpoint
			r := gin.New()
			r.PUT("/posts/:id", handlers.Update)

			// Create Request
			w := httptest.NewRecorder()

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(test.inputBody)
			if err != nil {
				t.Error(err)
			}
			req := httptest.NewRequest("PUT", fmt.Sprintf("/posts/%s", test.paramID), &buf)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)

			if test.errorMessage != "" {
				assert.Equal(t, w.Body.String(), test.errorMessage)
			}

			if test.expectedResponseBody != nil {
				var responseStruct domain.Post
				err := json.Unmarshal(w.Body.Bytes(), &responseStruct)
				if err != nil {
					t.Error(err)
				}

				assert.Equal(t, test.expectedResponseBody.ID, responseStruct.ID)
				assert.Equal(t, test.expectedResponseBody.Title, responseStruct.Title)
				assert.Equal(t, test.expectedResponseBody.Content, responseStruct.Content)
			}
		})
	}
}

func TestHandler_Delete(t *testing.T) {
	tests := []struct {
		name                 string
		paramID              string
		expectedStatusCode   int
		expectedResponseBody string
		errorMessage         string
	}{
		{
			name:                 "ok",
			paramID:              "1",
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "true",
			errorMessage:         "",
		},
		{
			name:                 "invalid parameter format",
			paramID:              "qwerty",
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "false",
			errorMessage:         `{"message":"invalid parameter format"}`,
		},
		{
			name:                 "id is less than zero",
			paramID:              "-11",
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "false",
			errorMessage:         `{"message":"ID must be greater than or equal to zero"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			postRepo := mock_repository.NewMockPosts(c)
			handlers := handler.NewHandler(&repository.Repository{Posts: postRepo})

			boolValue, err := strconv.ParseBool(test.expectedResponseBody)
			if err != nil {
				t.Error(err)
			}
			postRepo.EXPECT().Delete(gomock.Any(), 1).Return(boolValue, nil).AnyTimes()

			// Init Endpoint
			r := gin.New()
			r.DELETE("/posts/:id", handlers.Delete)

			// Create Request
			w := httptest.NewRecorder()

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/posts/%s", test.paramID), nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)

			if test.errorMessage != "" {
				assert.Equal(t, w.Body.String(), test.errorMessage)
			} else {
				assert.Equal(t, w.Body.String(), test.expectedResponseBody)
			}
		})
	}
}
