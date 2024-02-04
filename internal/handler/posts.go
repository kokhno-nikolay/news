package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kokhno-nikolay/news/domain"
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
func (h *Handler) Get(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid parameter format")
		return
	}

	if id < 0 {
		newErrorResponse(c, http.StatusBadRequest, "ID must be greater than or equal to zero")
		return
	}

	post, err := h.repo.Get(c, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, post)
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
func (h *Handler) List(c *gin.Context) {
	posts, err := h.repo.List(c, nil)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, posts)
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
func (h *Handler) Create(c *gin.Context) {
	var input domain.PostInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	if err := validatePostInput(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	post, err := h.repo.Posts.Create(c, &input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, post)
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
func (h *Handler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "parameter format")
		return
	}

	if id < 0 {
		newErrorResponse(c, http.StatusBadRequest, "ID must be greater than or equal to zero")
		return
	}

	var input domain.PostInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	if err := validatePostInput(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	post, err := h.repo.Posts.Update(c, id, &input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, post)
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
func (h *Handler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "parameter format")
		return
	}

	if id < 0 {
		newErrorResponse(c, http.StatusBadRequest, "ID must be greater than or equal to zero")
		return
	}

	ok, err := h.repo.Delete(c, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, ok)
}

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
