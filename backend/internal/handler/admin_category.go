package handler

import (
	"blog-backend/internal/service"
	"blog-backend/pkg/response"
	"blog-backend/pkg/validator"

	"github.com/gin-gonic/gin"
)

// AdminCategoryHandler handles admin category endpoints.
type AdminCategoryHandler struct {
	svc *service.CategoryService
}

// NewAdminCategoryHandler creates a new AdminCategoryHandler.
func NewAdminCategoryHandler(svc *service.CategoryService) *AdminCategoryHandler {
	return &AdminCategoryHandler{svc: svc}
}

// List handles GET /api/admin/categories
func (h *AdminCategoryHandler) List(c *gin.Context) {
	categories, err := h.svc.ListAll()
	if err != nil {
		response.InternalError(c)
		return
	}
	response.Success(c, categories)
}

// Create handles POST /api/admin/categories
func (h *AdminCategoryHandler) Create(c *gin.Context) {
	var input service.CreateCategoryInput
	if err := validator.Bind(c, &input); err != nil {
		return
	}

	cat, err := h.svc.Create(input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, cat)
}

// Update handles PUT /api/admin/categories/:id
func (h *AdminCategoryHandler) Update(c *gin.Context) {
	id, err := validator.ParamID(c, "id")
	if err != nil {
		return
	}

	var input service.UpdateCategoryInput
	if err := validator.Bind(c, &input); err != nil {
		return
	}

	cat, err := h.svc.Update(id, input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, cat)
}

// Delete handles DELETE /api/admin/categories/:id
func (h *AdminCategoryHandler) Delete(c *gin.Context) {
	id, err := validator.ParamID(c, "id")
	if err != nil {
		return
	}

	if err := h.svc.Delete(id); err != nil {
		response.InternalError(c)
		return
	}

	response.Success(c, nil)
}
