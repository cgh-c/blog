package handler

import (
	"blog-backend/internal/repository"
	"blog-backend/pkg/response"
	"blog-backend/pkg/validator"

	"github.com/gin-gonic/gin"
)

// AdminSettingHandler handles admin setting endpoints.
type AdminSettingHandler struct {
	repo *repository.SettingRepo
}

// NewAdminSettingHandler creates a new AdminSettingHandler.
func NewAdminSettingHandler(repo *repository.SettingRepo) *AdminSettingHandler {
	return &AdminSettingHandler{repo: repo}
}

// SettingHandler handles public setting endpoints.
type SettingHandler struct {
	repo *repository.SettingRepo
}

// NewSettingHandler creates a new public SettingHandler.
func NewSettingHandler(repo *repository.SettingRepo) *SettingHandler {
	return &SettingHandler{repo: repo}
}

// GetPublic handles GET /api/settings — only public settings.
func (h *SettingHandler) GetPublic(c *gin.Context) {
	settings, err := h.repo.ListPublic()
	if err != nil {
		response.InternalError(c)
		return
	}

	// Convert to key-value map for easier frontend consumption
	result := make(map[string]string, len(settings))
	for _, s := range settings {
		result[s.Key] = s.Value
	}

	response.Success(c, result)
}

// GetAll handles GET /api/admin/settings — all settings.
func (h *AdminSettingHandler) GetAll(c *gin.Context) {
	settings, err := h.repo.ListAll()
	if err != nil {
		response.InternalError(c)
		return
	}
	response.Success(c, settings)
}

type updateSettingsRequest struct {
	Settings map[string]string `json:"settings" binding:"required"`
}

// Update handles PUT /api/admin/settings — batch update.
func (h *AdminSettingHandler) Update(c *gin.Context) {
	var req updateSettingsRequest
	if err := validator.Bind(c, &req); err != nil {
		return
	}

	for key, value := range req.Settings {
		if err := h.repo.Upsert(key, value); err != nil {
			response.InternalError(c)
			return
		}
	}

	// Return updated settings
	settings, err := h.repo.ListAll()
	if err != nil {
		response.InternalError(c)
		return
	}

	response.Success(c, settings)
}
