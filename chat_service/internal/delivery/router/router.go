package router

import (
	"chat_service/internal/delivery/middleware"
	"chat_service/internal/delivery/rest"
	"chat_service/internal/delivery/ws"
	"chat_service/internal/infrastructure"
	"chat_service/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func SetupRouter(uc *usecase.UseCase, token *infrastructure.JWTmaker) *gin.Engine {
	r := gin.Default()

	httpRepo := rest.NewHTTPHandler(uc)

	r.GET("/messages", httpRepo.GetMessages)
	middleware := middleware.NewMiddleware(token)
	protected := r.Group("/")
	protected.Use(middleware.RequireAuth())
	{
		protected.DELETE("/messages", httpRepo.DeleteMessage)
	}

	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	wsHandler := ws.NewWsHandler(uc, upgrader)

	r.GET("/ws", middleware.RequireAuthWS(), wsHandler.HandleWS)

	return r
}
