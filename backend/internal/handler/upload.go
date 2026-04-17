package handler

import (
	"blog-backend/pkg/response"
	"blog-backend/pkg/upload"

	"github.com/gin-gonic/gin"
)

// UploadHandler handles file upload endpoints.
type UploadHandler struct {
	uploader *upload.Uploader
}

// NewUploadHandler creates a new UploadHandler.
func NewUploadHandler(uploader *upload.Uploader) *UploadHandler {
	return &UploadHandler{uploader: uploader}
}

// Upload handles POST /api/admin/upload
func (h *UploadHandler) Upload(c *gin.Context) {
	urlPath, err := h.uploader.Upload(c, "file")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"url": urlPath,
	})
}
