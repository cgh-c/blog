package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is the unified API response format.
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// PaginatedData wraps a list with pagination info.
type PaginatedData struct {
	List  any   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Size  int   `json:"size"`
}

// Success sends a 200 response with data.
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Created sends a 201 response with data.
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Response{
		Code:    0,
		Message: "created",
		Data:    data,
	})
}

// Error sends an error response with a given HTTP status code.
func Error(c *gin.Context, httpStatus int, message string) {
	c.JSON(httpStatus, Response{
		Code:    httpStatus,
		Message: message,
	})
}

// BadRequest sends a 400 error.
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

// Unauthorized sends a 401 error.
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message)
}

// Forbidden sends a 403 error.
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message)
}

// NotFound sends a 404 error.
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

// InternalError sends a 500 error with a generic message (never expose internals).
func InternalError(c *gin.Context) {
	Error(c, http.StatusInternalServerError, "internal server error")
}

// TooManyRequests sends a 429 error.
func TooManyRequests(c *gin.Context) {
	Error(c, http.StatusTooManyRequests, "too many requests, please try again later")
}

// Paginate sends a paginated success response.
func Paginate(c *gin.Context, list any, total int64, page, size int) {
	Success(c, PaginatedData{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}
