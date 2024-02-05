package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/kokhno-nikolay/news/docs"
	"github.com/kokhno-nikolay/news/internal/repository"
)

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) InitRoutes() http.Handler {
	router := gin.Default()

	// Swagger documentation route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Posts routes
	router.GET("/posts/:id", h.Get)
	router.GET("/posts/", h.List)
	router.POST("/posts", h.Create)
	router.PUT("/posts/:id", h.Update)
	router.DELETE("/posts/:id", h.Delete)

	return router
}
