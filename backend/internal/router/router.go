package router

import (
	"time"

	"blog-backend/internal/config"
	"blog-backend/internal/handler"
	"blog-backend/internal/repository"
	"blog-backend/internal/service"
	"blog-backend/pkg/auth"
	"blog-backend/pkg/ratelimit"
	"blog-backend/pkg/security"
	"blog-backend/pkg/upload"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
)

// Setup creates and configures the Gin engine with all routes.
func Setup(cfg *config.Config, db *gorm.DB) *gin.Engine {
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Security middleware
	r.Use(security.SecureHeaders())
	r.Use(security.CORS(security.CORSConfig{
		AllowedOrigins: cfg.CORSOrigins,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		MaxAge:         86400,
	}))

	// JWT manager
	jwtMgr := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiry, cfg.JWTExpiry/4)
	secureCookie := cfg.IsProduction()

	// Rate limiter for auth endpoints: 5 requests per minute
	authLimiter := ratelimit.NewIPRateLimiter(rate.Every(12*time.Second), 5)

	// File uploader
	uploadCfg := upload.DefaultConfig()
	uploadCfg.Dir = cfg.UploadDir
	uploadCfg.MaxSize = cfg.UploadMaxSize
	uploader := upload.NewUploader(uploadCfg)

	// Serve uploaded files
	r.Static("/uploads", cfg.UploadDir)

	// --- Repositories ---
	articleRepo := repository.NewArticleRepo(db)
	categoryRepo := repository.NewCategoryRepo(db)
	tagRepo := repository.NewTagRepo(db)
	settingRepo := repository.NewSettingRepo(db)

	// --- Services ---
	authService := service.NewAuthService(db)
	articleService := service.NewArticleService(articleRepo, tagRepo)
	categoryService := service.NewCategoryService(categoryRepo)

	// --- Handlers ---
	authHandler := handler.NewAuthHandler(authService, jwtMgr, secureCookie)
	publicArticleHandler := handler.NewPublicArticleHandler(articleService)
	publicCategoryHandler := handler.NewPublicCategoryHandler(categoryService)
	adminArticleHandler := handler.NewAdminArticleHandler(articleService)
	adminCategoryHandler := handler.NewAdminCategoryHandler(categoryService)
	adminTagHandler := handler.NewAdminTagHandler(db)
	settingHandler := handler.NewSettingHandler(settingRepo)
	adminSettingHandler := handler.NewAdminSettingHandler(settingRepo)
	uploadHandler := handler.NewUploadHandler(uploader)
	dashboardHandler := handler.NewDashboardHandler(db)

	// ==================== Routes ====================

	// Public tag list (shares the simple tag handler logic)
	api := r.Group("/api")
	{
		// Articles (public only)
		api.GET("/articles", publicArticleHandler.List)
		api.GET("/articles/:slug", publicArticleHandler.GetBySlug)
		api.GET("/archives", publicArticleHandler.Archives)

		// Categories (with public article count)
		api.GET("/categories", publicCategoryHandler.List)

		// Tags
		api.GET("/tags", adminTagHandler.List) // reuse — tags are always public

		// Settings (public only)
		api.GET("/settings", settingHandler.GetPublic)
	}

	// Auth endpoints (rate limited)
	authGroup := r.Group("/api/auth")
	authGroup.Use(authLimiter.Middleware())
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/logout", authHandler.Logout)
	}
	// /api/auth/me uses optional auth (no rate limit needed)
	r.GET("/api/auth/me", auth.OptionalAuth(jwtMgr), authHandler.Me)

	// Admin endpoints (JWT required)
	admin := r.Group("/api/admin")
	admin.Use(auth.RequireAuth(jwtMgr, secureCookie))
	{
		// Articles
		admin.GET("/articles", adminArticleHandler.List)
		admin.GET("/articles/:id", adminArticleHandler.GetByID)
		admin.POST("/articles", adminArticleHandler.Create)
		admin.PUT("/articles/:id", adminArticleHandler.Update)
		admin.DELETE("/articles/:id", adminArticleHandler.Delete)

		// Categories
		admin.GET("/categories", adminCategoryHandler.List)
		admin.POST("/categories", adminCategoryHandler.Create)
		admin.PUT("/categories/:id", adminCategoryHandler.Update)
		admin.DELETE("/categories/:id", adminCategoryHandler.Delete)

		// Tags
		admin.GET("/tags", adminTagHandler.List)
		admin.POST("/tags", adminTagHandler.Create)
		admin.DELETE("/tags/:id", adminTagHandler.Delete)

		// Settings
		admin.GET("/settings", adminSettingHandler.GetAll)
		admin.PUT("/settings", adminSettingHandler.Update)

		// Upload
		admin.POST("/upload", uploadHandler.Upload)

		// Dashboard
		admin.GET("/dashboard", dashboardHandler.Stats)
	}

	return r
}
