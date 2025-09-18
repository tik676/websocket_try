package usecase

import (
	"chat_service/internal/domain"
	"chat_service/internal/domain/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUsecase_SendMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMessageRepository(ctrl)
	mockProducer := mocks.NewMockProducerEvents(ctrl)

	usecase := NewUseCase(mockRepo, mockProducer)

	t.Run("succes", func(t *testing.T) {
		input := domain.Message{
			UserID:   1,
			Username: "testuser",
			Content:  "testmessage",
			Role:     "user",
		}

		expectedMessage := domain.Message{
			ID:        1,
			UserID:    1,
			Username:  "testuser",
			Content:   "testmessage",
			CreatedAt: time.Now(),
			Role:      "user",
			IsAnon:    false,
		}

		mockRepo.EXPECT().
			SaveMessage(gomock.Eq(input)).
			Return(expectedMessage, nil).
			Times(1)

		msg, err := usecase.SendMessage(input)

		assert.NoError(t, err)
		assert.NotNil(t, msg)
		assert.Equal(t, expectedMessage, msg)
	})

}
