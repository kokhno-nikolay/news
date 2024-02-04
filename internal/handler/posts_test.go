package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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

			postRepo.EXPECT().Get(gomock.Any(), test.paramID).Return(test.expectedResponseBody, nil).AnyTimes()

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
		})
	}
}
