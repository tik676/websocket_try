package rest

import (
	"chat_service/internal/domain"
	"chat_service/internal/domain/mocks"
	"chat_service/internal/usecase"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockMessageRepository(ctrl)
	mockProducer := mocks.NewMockProducerEvents(ctrl)
	usecase := usecase.NewUseCase(mockRepo, mockProducer)
	handler := NewHTTPHandler(usecase)

	router := gin.New()
	router.GET("/messages", handler.GetMessages)
	router.DELETE("/messages", handler.DeleteMessage)

	t.Run("get_messages_success", func(t *testing.T) {
		testMessage := domain.Message{
			UserID:   1,
			Username: "testuser",
			Content:  "test message",
			Role:     "user",
		}

		expectedSaved := testMessage

		mockRepo.EXPECT().
			MessageHistory(gomock.Eq(int64(10)), gomock.Eq(int64(0))).
			Return([]domain.Message{expectedSaved}, nil).
			Times(1)

		req := httptest.NewRequest("GET", "/messages?limit=10&offset=0", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		var response map[string][]domain.Message
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		messages := response["messages"]

		require.Len(t, messages, 1)
		assert.Equal(t, expectedSaved, messages[0])
	})

	t.Run("get_messages_invalid_limit or offset", func(t *testing.T) {
		req1 := httptest.NewRequest("GET", "/messages?offset=10", nil)
		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)

		assert.Equal(t, 400, w1.Code)

		req2 := httptest.NewRequest("GET", "/messages?limit=10", nil)
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		assert.Equal(t, 400, w2.Code)

		req3 := httptest.NewRequest("GET", "/messages?limit=abc&offset=1", nil)
		w3 := httptest.NewRecorder()
		router.ServeHTTP(w3, req3)

		assert.Equal(t, 400, w3.Code)

		req4 := httptest.NewRequest("GET", "/messages?limit=1&offset=abc", nil)
		w4 := httptest.NewRecorder()
		router.ServeHTTP(w4, req4)

		assert.Equal(t, 400, w4.Code)
	})

	t.Run("get_messages_repository_error", func(t *testing.T) {

		mockRepo.EXPECT().
			MessageHistory(gomock.Eq(int64(10)), gomock.Eq(int64(0))).
			Return([]domain.Message{}, errors.New("repo_error")).
			Times(1)

		req := httptest.NewRequest("GET", "/messages?limit=10&offset=0", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, 500, w.Code)
	})

	t.Run("delete_message_success", func(t *testing.T) {
		var msgID int64 = 1
		var userID int64 = 1

		mockProducer.EXPECT().
			SendDeleteMessageEvent(gomock.Any(), gomock.Eq(userID), gomock.Eq(msgID)).
			Return(nil).
			Times(1)

		mockRepo.EXPECT().
			DeleteMessage(gomock.Eq(msgID)).
			Return(nil).
			Times(1)

		body := ` {"id":1,"userID":1}`
		req := httptest.NewRequest("DELETE", "/messages", strings.NewReader(body))
		req.Header.Set("Content-type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("delete_message_invalid_data", func(t *testing.T) {

		req := httptest.NewRequest("DELETE", "/messages", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)

		body := `{"id":1,"userID":1`
		req2 := httptest.NewRequest("DELETE", "/messages", strings.NewReader(body))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, 400, w2.Code)

		body3 := `{"id":"not_number","userID":1}`
		req3 := httptest.NewRequest("DELETE", "/messages", strings.NewReader(body3))
		req3.Header.Set("Content-Type", "application/json")
		w3 := httptest.NewRecorder()
		router.ServeHTTP(w3, req3)
		assert.Equal(t, 400, w3.Code)
	})

	t.Run("delete_message_fail_repo", func(t *testing.T) {
		var msgID int64 = 1
		var userID int64 = 1

		mockProducer.EXPECT().
			SendDeleteMessageEvent(gomock.Any(), gomock.Eq(userID), gomock.Eq(msgID)).
			Return(nil).
			Times(1)

		mockRepo.EXPECT().
			DeleteMessage(gomock.Eq(msgID)).
			Return(errors.New("failed_Repo")).
			Times(1)

		body := ` {"id":1,"userID":1}`
		req := httptest.NewRequest("DELETE", "/messages", strings.NewReader(body))
		req.Header.Set("Content-type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})
}
