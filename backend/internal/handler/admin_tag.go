package handler

import (
	"blog-backend/internal/model"
	"blog-backend/pkg/response"
	"blog-backend/pkg/validator"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AdminTagHandler handles admin tag endpoints.
type AdminTagHandler struct {
	db *gorm.DB
}

// NewAdminTagHandler creates a new AdminTagHandler.
func NewAdminTagHandler(db *gorm.DB) *AdminTagHandler {
	return &AdminTagHandler{db: db}
}

// List handles GET /api/admin/tags
func (h *AdminTagHandler) List(c *gin.Context) {
	var tags []model.Tag
	if err := h.db.Order("name ASC").Find(&tags).Error; err != nil {
		response.InternalError(c)
		return
	}
	response.Success(c, tags)
}

type createTagRequest struct {
	Name string `json:"name" binding:"required,min=1,max=50"`
}

// Create handles POST /api/admin/tags
func (h *AdminTagHandler) Create(c *gin.Context) {
	var req createTagRequest
	if err := validator.Bind(c, &req); err != nil {
		return
	}

	tag := model.Tag{
		Name: req.Name,
		Slug: generateTagSlug(req.Name),
	}

	if err := h.db.Create(&tag).Error; err != nil {
		response.BadRequest(c, "tag may already exist")
		return
	}

	response.Created(c, tag)
}

// Delete handles DELETE /api/admin/tags/:id
func (h *AdminTagHandler) Delete(c *gin.Context) {
	id, err := validator.ParamID(c, "id")
	if err != nil {
		return
	}

	// Clean up article_tags
	h.db.Exec("DELETE FROM article_tags WHERE tag_id = ?", id)

	if err := h.db.Delete(&model.Tag{}, id).Error; err != nil {
		response.InternalError(c)
		return
	}

	response.Success(c, nil)
}

func generateTagSlug(name string) string {
	// Reuse the same slug logic from the service package
	// For tags, just lowercase and replace spaces
	slug := ""
	for _, r := range name {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-' {
			slug += string(r)
		} else if r >= 'A' && r <= 'Z' {
			slug += string(r + 32) // to lowercase
		} else if r == ' ' {
			slug += "-"
		}
	}
	if slug == "" {
		slug = "tag"
	}
	return slug
}
