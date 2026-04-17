package handler

import (
	"blog-backend/internal/model"
	"blog-backend/internal/service"
	"blog-backend/pkg/response"
	"blog-backend/pkg/validator"

	"github.com/gin-gonic/gin"
)

// AdminArticleHandler handles admin article endpoints.
type AdminArticleHandler struct {
	svc *service.ArticleService
}

// NewAdminArticleHandler creates a new AdminArticleHandler.
func NewAdminArticleHandler(svc *service.ArticleService) *AdminArticleHandler {
	return &AdminArticleHandler{svc: svc}
}

type adminArticleListQuery struct {
	validator.Pagination
	CategoryID *uint             `form:"category_id"`
	TagID      *uint             `form:"tag_id"`
	Visibility *model.Visibility `form:"visibility"`
}

// List handles GET /api/admin/articles — all articles for admin.
func (h *AdminArticleHandler) List(c *gin.Context) {
	var q adminArticleListQuery
	if err := validator.BindQuery(c, &q); err != nil {
		return
	}

	articles, total, err := h.svc.ListAll(q.GetOffset(), q.GetSize(), q.CategoryID, q.TagID, q.Visibility)
	if err != nil {
		response.InternalError(c)
		return
	}

	response.Paginate(c, articles, total, q.GetPage(), q.GetSize())
}

// GetByID handles GET /api/admin/articles/:id
func (h *AdminArticleHandler) GetByID(c *gin.Context) {
	id, err := validator.ParamID(c, "id")
	if err != nil {
		return
	}

	article, err := h.svc.GetByID(id)
	if err != nil {
		response.NotFound(c, "article not found")
		return
	}

	response.Success(c, article)
}

// Create handles POST /api/admin/articles
func (h *AdminArticleHandler) Create(c *gin.Context) {
	var input service.CreateArticleInput
	if err := validator.Bind(c, &input); err != nil {
		return
	}

	article, err := h.svc.Create(input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, article)
}

// Update handles PUT /api/admin/articles/:id
func (h *AdminArticleHandler) Update(c *gin.Context) {
	id, err := validator.ParamID(c, "id")
	if err != nil {
		return
	}

	var input service.UpdateArticleInput
	if err := validator.Bind(c, &input); err != nil {
		return
	}

	article, err := h.svc.Update(id, input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, article)
}

// Delete handles DELETE /api/admin/articles/:id
func (h *AdminArticleHandler) Delete(c *gin.Context) {
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
