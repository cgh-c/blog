package handler

import (
	"blog-backend/internal/service"
	"blog-backend/pkg/response"
	"blog-backend/pkg/validator"

	"github.com/gin-gonic/gin"
)

// PublicArticleHandler handles public (unauthenticated) article endpoints.
type PublicArticleHandler struct {
	svc *service.ArticleService
}

// NewPublicArticleHandler creates a new PublicArticleHandler.
func NewPublicArticleHandler(svc *service.ArticleService) *PublicArticleHandler {
	return &PublicArticleHandler{svc: svc}
}

type articleListQuery struct {
	validator.Pagination
	CategoryID *uint `form:"category_id"`
	TagID      *uint `form:"tag_id"`
}

// List handles GET /api/articles — public article list with pagination.
func (h *PublicArticleHandler) List(c *gin.Context) {
	var q articleListQuery
	if err := validator.BindQuery(c, &q); err != nil {
		return
	}

	articles, total, err := h.svc.ListPublic(q.GetOffset(), q.GetSize(), q.CategoryID, q.TagID)
	if err != nil {
		response.InternalError(c)
		return
	}

	response.Paginate(c, articles, total, q.GetPage(), q.GetSize())
}

// GetBySlug handles GET /api/articles/:slug — public article detail.
func (h *PublicArticleHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		response.BadRequest(c, "slug is required")
		return
	}

	article, err := h.svc.GetPublicBySlug(slug)
	if err != nil {
		response.NotFound(c, "article not found")
		return
	}

	response.Success(c, article)
}

// Archives handles GET /api/archives — public article archives.
func (h *PublicArticleHandler) Archives(c *gin.Context) {
	archives, err := h.svc.GetPublicArchives()
	if err != nil {
		response.InternalError(c)
		return
	}

	response.Success(c, archives)
}
