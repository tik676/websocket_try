package router

import (
	"chat_service/internal/delivery/middleware"
	"chat_service/internal/delivery/rest"
	"chat_service/internal/delivery/ws"
	"chat_service/internal/infrastructure"
	"chat_service/internal/usecase"

	"github.com/gin-gonic/gin"
)

func SetupRouter(uc *usecase.UseCase, token *infrastructure.JWTmaker, ws *ws.WsHandler) *gin.Engine {
	r := gin.Default()

	httpRepo := rest.NewHTTPHandler(uc)

	r.GET("/messages", httpRepo.GetMessages)
	r.DELETE("/message", httpRepo.DeleteMessage)
	middleware := middleware.NewMiddleware(token)
	protected := r.Group("/")
	protected.Use(middleware.RequireAuth())
	{
		protected.DELETE("/messages", httpRepo.DeleteMessage)
	}

	r.GET("/ws", middleware.RequireAuthWS(), ws.HandleWS)

	return r
}
