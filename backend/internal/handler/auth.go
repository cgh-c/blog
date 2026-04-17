package handler

import (
	"blog-backend/internal/service"
	"blog-backend/pkg/auth"
	"blog-backend/pkg/response"
	"blog-backend/pkg/validator"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	authService  *service.AuthService
	jwtMgr       *auth.JWTManager
	secureCookie bool
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService *service.AuthService, jwtMgr *auth.JWTManager, secureCookie bool) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		jwtMgr:       jwtMgr,
		secureCookie: secureCookie,
	}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login handles POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := validator.Bind(c, &req); err != nil {
		return
	}

	user, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		response.Unauthorized(c, "invalid username or password")
		return
	}

	token, err := h.jwtMgr.GenerateToken(user.ID, user.Username)
	if err != nil {
		response.InternalError(c)
		return
	}

	auth.SetTokenCookie(c, token, h.secureCookie)

	response.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
	})
}

// Logout handles POST /api/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	auth.ClearTokenCookie(c, h.secureCookie)
	response.Success(c, nil)
}

// Me handles GET /api/auth/me — returns current auth status
func (h *AuthHandler) Me(c *gin.Context) {
	userID := auth.GetUserID(c)
	username := auth.GetUsername(c)

	if userID == 0 {
		response.Success(c, gin.H{
			"authenticated": false,
		})
		return
	}

	response.Success(c, gin.H{
		"authenticated": true,
		"id":            userID,
		"username":      username,
	})
}
