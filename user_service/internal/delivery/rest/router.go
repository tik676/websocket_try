package rest

import (
	"user_service/internal/domain"
	"user_service/internal/middleware"
	"user_service/internal/usecase"

	"github.com/gin-gonic/gin"
)

func SetupRouter(usecase *usecase.UseCase, tokenManager domain.TokenManager) *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	})

	handler := NewHTTPHandler(usecase)
	middleware := middleware.NewAuthMiddleware(tokenManager)

	r.POST("/register", handler.RegisterHandler)
	r.POST("/login/anon", handler.LoginAnonHandler)
	r.POST("/login", handler.LoginHandler)
	r.POST("/refresh", handler.RefreshTokenHandler)
	r.POST("/logout", handler.LogoutHandler)
	protected := r.Group("/")
	protected.Use(middleware.RequireAuth())
	{

	}

	return r
}
