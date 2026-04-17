package handler

import (
	"blog-backend/internal/service"
	"blog-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// PublicCategoryHandler handles public category endpoints.
type PublicCategoryHandler struct {
	svc *service.CategoryService
}

// NewPublicCategoryHandler creates a new PublicCategoryHandler.
func NewPublicCategoryHandler(svc *service.CategoryService) *PublicCategoryHandler {
	return &PublicCategoryHandler{svc: svc}
}

// List handles GET /api/categories — categories with public article counts.
func (h *PublicCategoryHandler) List(c *gin.Context) {
	categories, err := h.svc.ListWithPublicCount()
	if err != nil {
		response.InternalError(c)
		return
	}

	response.Success(c, categories)
}
