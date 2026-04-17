package validator

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Bind binds the request body to the given struct and validates it.
// On validation failure, it writes a 400 JSON response and returns an error.
func Bind(c *gin.Context, obj any) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": formatError(err),
		})
		return err
	}
	return nil
}

// BindQuery binds query parameters to the given struct and validates.
func BindQuery(c *gin.Context, obj any) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": formatError(err),
		})
		return err
	}
	return nil
}

// formatError converts validation errors to a readable message.
func formatError(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		if len(validationErrors) > 0 {
			fe := validationErrors[0]
			switch fe.Tag() {
			case "required":
				return fe.Field() + " is required"
			case "min":
				return fe.Field() + " must be at least " + fe.Param() + " characters"
			case "max":
				return fe.Field() + " must be at most " + fe.Param() + " characters"
			case "email":
				return fe.Field() + " must be a valid email"
			case "oneof":
				return fe.Field() + " must be one of: " + fe.Param()
			default:
				return fe.Field() + " is invalid"
			}
		}
	}
	return "invalid request parameters"
}

// Pagination holds common pagination parameters.
type Pagination struct {
	Page int `form:"page" binding:"omitempty,min=1"`
	Size int `form:"size" binding:"omitempty,min=1,max=100"`
}

// GetPage returns the page number with a default of 1.
func (p Pagination) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

// GetSize returns the page size with a default of 10.
func (p Pagination) GetSize() int {
	if p.Size <= 0 {
		return 10
	}
	return p.Size
}

// GetOffset returns the SQL offset based on page and size.
func (p Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetSize()
}

// ParamID extracts an unsigned integer ID from the URL path parameter.
func ParamID(c *gin.Context, name string) (uint, error) {
	idStr := c.Param(name)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid " + name,
		})
		return 0, err
	}
	return uint(id), nil
}
