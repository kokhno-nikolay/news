package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

	// News routes
	router.GET("/news/:id", h.Get)
	router.GET("/news/", h.List)
	router.POST("/news", h.Create)
	router.PUT("/news/:id", h.Update)
	router.DELETE("/news/:id", h.Delete)

	return router
}
