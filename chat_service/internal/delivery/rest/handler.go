package rest

import (
	"chat_service/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	uc *usecase.UseCase
}

func NewHTTPHandler(uc *usecase.UseCase) *HTTPHandler {
	return &HTTPHandler{uc: uc}
}

func (h *HTTPHandler) GetMessages(c *gin.Context) {
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid limit"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid offset"})
		return
	}

	messages, err := h.uc.GetMessages(int64(limit), int64(offset))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"messages": messages})
}

func (h *HTTPHandler) DeleteMessage(c *gin.Context) {
	var req struct {
		ID_message int64 `json:"id"`
		UserID     int64 `json:"userID"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	if err := h.uc.DeleteMessage(req.ID_message, req.UserID); err != nil {
		c.JSON(400, gin.H{"error": "Failed to delete message"})
		return
	}

	c.JSON(200, gin.H{"messages": "succes"})
}
