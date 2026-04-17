package handler

import (
	"blog-backend/internal/model"
	"blog-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DashboardHandler handles dashboard statistics.
type DashboardHandler struct {
	db *gorm.DB
}

// NewDashboardHandler creates a new DashboardHandler.
func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

// Stats handles GET /api/admin/dashboard
func (h *DashboardHandler) Stats(c *gin.Context) {
	var totalArticles, publicArticles, draftArticles, totalCategories, totalTags int64
	var totalViews int64

	h.db.Model(&model.Article{}).Count(&totalArticles)
	h.db.Model(&model.Article{}).Where("visibility = ?", "public").Count(&publicArticles)
	h.db.Model(&model.Article{}).Where("visibility = ?", "draft").Count(&draftArticles)
	h.db.Model(&model.Category{}).Count(&totalCategories)
	h.db.Model(&model.Tag{}).Count(&totalTags)
	h.db.Model(&model.Article{}).Select("COALESCE(SUM(view_count), 0)").Scan(&totalViews)

	response.Success(c, gin.H{
		"total_articles":  totalArticles,
		"public_articles": publicArticles,
		"draft_articles":  draftArticles,
		"total_categories": totalCategories,
		"total_tags":      totalTags,
		"total_views":     totalViews,
	})
}
