package rest

import (
	"chat_service/internal/domain"
	"chat_service/internal/domain/mocks"
	"chat_service/internal/usecase"
	"encoding/json"
	"errors"
	"net/http/httptest"
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

	usecase := usecase.NewUseCase(mockRepo, nil)
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

	t.Run("get_messages_repository_error", func(t *testing.T) {

		mockRepo.EXPECT().
			MessageHistory(gomock.Eq(int64(10)), gomock.Eq(int64(0))).
			Return([]domain.Message{}, errors.New("repo_error")).
			Times(1)

		msgs, err := usecase.GetMessages(int64(10), int64(0))
		require.Error(t, err)
		assert.Equal(t, []domain.Message{}, msgs)
	})
}
