package handler

import "github.com/gin-gonic/gin"

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

}
