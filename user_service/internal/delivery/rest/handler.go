package rest

import (
	"log"
	"strings"
	"time"
	"user_service/internal/domain"
	"user_service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	usecase *usecase.UseCase
}

func NewHTTPHandler(usecase *usecase.UseCase) *HTTPHandler {
	return &HTTPHandler{usecase: usecase}
}

type UserResponse struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Role         string    `json:"role"`
	RegisteredAt time.Time `json:"registered_at"`
}

func (h *HTTPHandler) RegisterHandler(c *gin.Context) {
	var userInput domain.AuthorizationInput
	if err := c.BindJSON(&userInput); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON format"})
		return
	}

	if userInput.Name == "" || userInput.Password == "" {
		c.JSON(400, gin.H{"error": "Name and password are required"})
		return
	}

	if len(userInput.Password) < 3 {
		c.JSON(400, gin.H{"error": "Password must be at least 3 characters"})
		return
	}

	user, err := h.usecase.RegisterUser(userInput)
	if err != nil {
		log.Printf("Registration error: %v", err)
		c.JSON(500, gin.H{"error": "User already exists"})
		return
	}

	response := UserResponse{
		ID:           user.ID,
		Name:         user.Name,
		Role:         user.Role,
		RegisteredAt: user.RegisteredAt,
	}

	c.JSON(201, response)
}

func (h *HTTPHandler) LoginHandler(c *gin.Context) {
	var userInput domain.AuthorizationInput
	if err := c.BindJSON(&userInput); err != nil {
		c.JSON(400, gin.H{"error": "failed to parse json"})
		return
	}

	token, err := h.usecase.LoginUser(userInput)
	if err != nil {
		if strings.Contains(err.Error(), "invalid password") ||
			strings.Contains(err.Error(), "user not found") {
			c.JSON(401, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(200, token)
}

func (h *HTTPHandler) LoginAnonHandler(c *gin.Context) {
	token, err := h.usecase.LoginAnonUser()
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(200, token)
}

func (h *HTTPHandler) RefreshTokenHandler(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON format"})
		return
	}

	if req.RefreshToken == "" {
		c.JSON(400, gin.H{"error": "refresh_token is required"})
		return
	}

	token, err := h.usecase.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	c.JSON(200, token)
}

func (h *HTTPHandler) LogoutHandler(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON format"})
		return
	}

	if req.RefreshToken == "" {
		c.JSON(400, gin.H{"error": "refresh_token is required"})
		return
	}

	err := h.usecase.LogoutUser(req.RefreshToken)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	c.JSON(200, gin.H{"message": "successful"})
}
