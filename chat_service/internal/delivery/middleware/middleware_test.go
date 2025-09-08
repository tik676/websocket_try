package middleware

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockMaker struct {
	ShouldFail bool
	UserID     int64
	Name       string
	Role       string
	Error      error
}

func (m *mockMaker) VerifyToken(token string) (userID int64, name, role string, err error) {
	if m.ShouldFail == true {
		return 0, "", "", m.Error
	}
	return m.UserID, m.Name, m.Role, nil
}

func TestRequireAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &mockMaker{
		ShouldFail: false,
		UserID:     1,
		Name:       "lol",
		Role:       "user",
	}

	middleware := NewMiddleware(mock)
	router := gin.New()
	router.Use(middleware.RequireAuth())
	router.GET("/protected", func(c *gin.Context) {
		userID, exists1 := c.Get("user_id")
		name, exists2 := c.Get("name")
		role, exists3 := c.Get("role")

		if !exists1 || !exists2 || !exists3 {
			c.JSON(500, gin.H{"error": "missing context values"})
			return
		}

		c.JSON(200, gin.H{
			"user_id": userID,
			"name":    name,
			"role":    role,
		})
	})

	req1 := httptest.NewRequest("GET", "/protected", nil)
	req1.Header.Set("Authorization", "Bearer valid_token")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	if w1.Code != 200 {
		t.Errorf("expected 200, got %v", w1.Code)
	}

	req2 := httptest.NewRequest("GET", "/protected", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != 401 {
		t.Errorf("expected 401, got %d", w2.Code)
	}

	req3 := httptest.NewRequest("GET", "/protected", nil)
	req3.Header.Set("Authorization", "InvalidFormat")
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	if w3.Code != 401 {
		t.Errorf("expected 401, got %d", w3.Code)
	}

	mock.ShouldFail = true
	mock.Error = errors.New("invalid token")

	req4 := httptest.NewRequest("GET", "/protected", nil)
	req4.Header.Set("Authorization", "Bearer bad_token")
	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req4)

	if w4.Code != 401 {
		t.Errorf("expected 401, got %d", w4.Code)
	}
}

func TestRequireAuthWS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &mockMaker{
		ShouldFail: false,
		UserID:     1,
		Name:       "lol",
		Role:       "user",
	}

	middleware := NewMiddleware(mock)
	router := gin.New()
	router.Use(middleware.RequireAuthWS())
	router.GET("/ws", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ws connected"})
	})

	req1 := httptest.NewRequest("GET", "/ws?token=valid_token", nil)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	if w1.Code != 200 {
		t.Errorf("expected 200, got %d", w1.Code)
	}

	req2 := httptest.NewRequest("GET", "/ws", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != 401 {
		t.Errorf("expected 401, got %d", w2.Code)
	}

	req3 := httptest.NewRequest("GET", "/ws?token=", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	if w3.Code != 401 {
		t.Errorf("expected 401,got %d", w3.Code)
	}

	mock.ShouldFail = true
	mock.Error = errors.New("invalid token")

	req4 := httptest.NewRequest("GET", "/ws?token=bad_token", nil)
	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req4)

	if w4.Code != 401 {
		t.Errorf("expected 401, got %d", w4.Code)
	}
}
