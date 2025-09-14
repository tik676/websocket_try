package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"user_service/internal/domain/mocks"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware_RequireAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockToken := mocks.NewMockTokenManager(ctrl)
	middleware := NewAuthMiddleware(mockToken)

	t.Run("success", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		mockToken.EXPECT().
			VerifyToken("valid token").
			Return(int64(123), "testuser", "user", nil).
			Times(1)

		router := gin.New()
		router.Use(middleware.RequireAuth())
		router.GET("/protected", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			name, _ := c.Get("name")
			role, _ := c.Get("role")

			c.JSON(http.StatusOK, gin.H{
				"user_id": userID,
				"name":    name,
				"role":    role,
			})

		})

		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer valid token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("without_Header", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		router := gin.New()
		router.Use(middleware.RequireAuth())
		router.GET("/protected", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			name, _ := c.Get("name")
			role, _ := c.Get("role")

			c.JSON(http.StatusOK, gin.H{
				"user_id": userID,
				"name":    name,
				"role":    role,
			})
		})
		req := httptest.NewRequest("GET", "/protected", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
	})

	t.Run("empty_Header", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		router := gin.New()
		router.Use(middleware.RequireAuth())
		router.GET("/protected", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			name, _ := c.Get("name")
			role, _ := c.Get("role")

			c.JSON(http.StatusOK, gin.H{
				"user_id": userID,
				"name":    name,
				"role":    role,
			})
		})
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
	})

	t.Run("repository_error", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		mockToken.EXPECT().
			VerifyToken("Bad_token").
			Return(int64(0), "", "", errors.New("database error")).
			Times(1)

		router := gin.New()
		router.Use(middleware.RequireAuth())
		router.GET("/protected", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			name, _ := c.Get("name")
			role, _ := c.Get("role")

			c.JSON(http.StatusOK, gin.H{
				"user_id": userID,
				"name":    name,
				"role":    role,
			})
		})
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer Bad_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
	})

	t.Run("missing_bearer_prefix", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		router := gin.New()
		router.Use(middleware.RequireAuth())
		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
	})
}
