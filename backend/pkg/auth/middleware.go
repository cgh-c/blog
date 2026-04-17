package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	// CookieName is the name of the JWT cookie.
	CookieName = "blog_token"
	// ContextKeyUserID is the gin context key for the authenticated user ID.
	ContextKeyUserID = "user_id"
	// ContextKeyUsername is the gin context key for the authenticated username.
	ContextKeyUsername = "username"
)

// RequireAuth returns a Gin middleware that rejects unauthenticated requests.
// On success, it sets user_id and username in the gin context.
// It also handles sliding token refresh automatically.
func RequireAuth(jwtMgr *JWTManager, secureCookie bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(CookieName)
		if err != nil || tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "authentication required",
			})
			return
		}

		claims, err := jwtMgr.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "invalid or expired token",
			})
			return
		}

		// Sliding refresh: if token is close to expiry, issue a new one
		if jwtMgr.ShouldRefresh(claims) {
			if newToken, err := jwtMgr.GenerateToken(claims.UserID, claims.Username); err == nil {
				SetTokenCookie(c, newToken, secureCookie)
			}
		}

		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyUsername, claims.Username)
		c.Next()
	}
}

// OptionalAuth returns a Gin middleware that does NOT reject unauthenticated requests.
// If a valid token is present, it sets user info in the context; otherwise it continues.
func OptionalAuth(jwtMgr *JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(CookieName)
		if err != nil || tokenStr == "" {
			c.Next()
			return
		}

		claims, err := jwtMgr.ParseToken(tokenStr)
		if err != nil {
			c.Next()
			return
		}

		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyUsername, claims.Username)
		c.Next()
	}
}

// SetTokenCookie sets the JWT token as an httpOnly cookie.
func SetTokenCookie(c *gin.Context, token string, secure bool) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		CookieName,
		token,
		86400, // 24 hours in seconds
		"/",
		"",
		secure, // Secure flag (true in production with HTTPS)
		true,   // httpOnly — not accessible via JavaScript
	)
}

// ClearTokenCookie removes the JWT cookie.
func ClearTokenCookie(c *gin.Context, secure bool) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(CookieName, "", -1, "/", "", secure, true)
}

// GetUserID extracts the authenticated user ID from gin context.
// Returns 0 if not authenticated.
func GetUserID(c *gin.Context) uint {
	if v, exists := c.Get(ContextKeyUserID); exists {
		if id, ok := v.(uint); ok {
			return id
		}
	}
	return 0
}

// GetUsername extracts the authenticated username from gin context.
func GetUsername(c *gin.Context) string {
	if v, exists := c.Get(ContextKeyUsername); exists {
		if name, ok := v.(string); ok {
			return name
		}
	}
	return ""
}
