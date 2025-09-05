package rest

import (
	"chat_service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	uc *usecase.UseCase
}

func NewHTTPHandler(uc *usecase.UseCase) *HTTPHandler {
	return &HTTPHandler{uc: uc}
}

func (h *HTTPHandler) GetMessages(c *gin.Context) {

}
