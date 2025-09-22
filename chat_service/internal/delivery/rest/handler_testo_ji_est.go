package rest

/*
import (
	"chat_service/internal/domain"
	"chat_service/internal/usecase"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockMessageRepo struct {
	Saved      []domain.Message
	SaveErr    error
	History    []domain.Message
	HistoryErr error
	DeleteErr  error
}

func (m *mockMessageRepo) SaveMessage(msg domain.Message) (domain.Message, error) {
	m.Saved = append(m.Saved, msg)
	if m.SaveErr != nil {
		return domain.Message{}, m.SaveErr
	}
	return msg, nil
}

func (m *mockMessageRepo) MessageHistory(limit, offset int64) ([]domain.Message, error) {
	if m.HistoryErr != nil {
		return []domain.Message{}, m.HistoryErr
	}
	return m.History, nil
}

func (m *mockMessageRepo) DeleteMessage(id int64) error {
	return m.DeleteErr
}

func TestGetMessages(t *testing.T) {
	gin.SetMode(gin.TestMode)

	messages := []domain.Message{
		{UserID: 1, Username: "lol", Content: "sad", Role: "user"},
		{UserID: 2, Username: "kek", Content: "das", Role: "user"},
	}

	mockRepo := &mockMessageRepo{
		History: messages,
	}

	realUseCase := usecase.NewUseCase(mockRepo)
	httpHandler := NewHTTPHandler(realUseCase)

	router := gin.New()
	router.GET("/messages", httpHandler.GetMessages)

	req1 := httptest.NewRequest("GET", "/messages?limit=10&offset=0", nil)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	if w1.Code != 200 {
		t.Errorf("expected 200, got %d", w1.Code)
	}

	var response map[string][]domain.Message
	err := json.Unmarshal(w1.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if len(response["messages"]) != 2 {
		t.Errorf("expected 2 messages, got %d", len(response["messages"]))
	}

	req2 := httptest.NewRequest("GET", "/messages?limit=abc&offset=0", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != 400 {
		t.Errorf("expected 400, got %d", w2.Code)
	}

	req3 := httptest.NewRequest("GET", "/messages?limit=10&offset=abc", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	if w3.Code != 400 {
		t.Errorf("expected 400, got %d", w3.Code)
	}

	req4 := httptest.NewRequest("GET", "/messages", nil)
	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req4)

	if w4.Code != 400 {
		t.Errorf("expected 400, got %d", w4.Code)
	}

	req5 := httptest.NewRequest("GET", "/messages?limit=10", nil)
	w5 := httptest.NewRecorder()
	router.ServeHTTP(w5, req5)

	if w5.Code != 400 {
		t.Errorf("expected 400, got %d", w5.Code)
	}

	req6 := httptest.NewRequest("GET", "/messages?offset=10", nil)
	w6 := httptest.NewRecorder()
	router.ServeHTTP(w6, req6)

	if w6.Code != 400 {
		t.Errorf("expected 400, got %d", w6.Code)
	}

}

func TestGetMessages_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &mockMessageRepo{
		HistoryErr: errors.New("database connection failed"),
	}

	uc := usecase.NewUseCase(mock)
	httpHandler := NewHTTPHandler(uc)

	router := gin.New()
	router.GET("/messages", httpHandler.GetMessages)
	req := httptest.NewRequest("GET", "/messages?limit=10&offset=0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Errorf("expected 500, got %d", w.Code)
	}

}

func TestDeleteMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &mockMessageRepo{}
	uc := usecase.NewUseCase(mock)
	httpHandler := NewHTTPHandler(uc)

	router := gin.New()
	router.DELETE("/messages", httpHandler.DeleteMessage)
	body := `{"id": 1}`
	req := httptest.NewRequest("DELETE", "/messages", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}

	body = `{"id": "abc"}`
	req1 := httptest.NewRequest("DELETE", "/messages", strings.NewReader(body))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	if w1.Code != 400 {
		t.Errorf("expected 400, got %d", w1.Code)
	}

	req2 := httptest.NewRequest("DELETE", "/messages", strings.NewReader(""))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != 400 {
		t.Errorf("expected 400, got %d", w2.Code)
	}

	body = `{"id": 1000}`
	mock.DeleteErr = errors.New("invalid id")
	req3 := httptest.NewRequest("DELETE", "/messages", strings.NewReader(body))
	req3.Header.Set("Content-Type", "application/json")
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	if w3.Code != 400 {
		t.Errorf("exppected 400, got %d", w3.Code)
	}
}
*/
