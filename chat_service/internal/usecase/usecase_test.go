package usecase

import (
	"chat_service/internal/domain"
	"chat_service/internal/domain/mocks"
	"errors"
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

	t.Run("empty_content", func(t *testing.T) {
		input := domain.Message{
			ID:        1,
			UserID:    1,
			Username:  "testuser",
			Content:   "",
			CreatedAt: time.Now(),
			Role:      "user",
			IsAnon:    false,
		}

		msg, err := usecase.SendMessage(input)

		assert.Equal(t, domain.Message{}, msg)
		assert.Equal(t, "error message can't be empty", err.Error())
	})

	t.Run("UserID_by_zero", func(t *testing.T) {
		input := domain.Message{
			ID:        1,
			UserID:    0,
			Username:  "testuser",
			Content:   "lol",
			CreatedAt: time.Now(),
			Role:      "user",
			IsAnon:    false,
		}

		msg, err := usecase.SendMessage(input)
		assert.Equal(t, domain.Message{}, msg)
		assert.Equal(t, "user id cannot be zero", err.Error())
	})

	t.Run("role_isAnon", func(t *testing.T) {
		input := domain.Message{
			ID:        1,
			UserID:    1,
			Username:  "testuser",
			Content:   "lol",
			CreatedAt: time.Now(),
			Role:      "anon",
			IsAnon:    false,
		}

		expectedInput := input
		expectedInput.IsAnon = true

		expected := domain.Message{
			ID:        1,
			UserID:    1,
			Username:  "testuser",
			Content:   "lol",
			CreatedAt: time.Now(),
			Role:      "anon",
			IsAnon:    true,
		}

		mockRepo.EXPECT().
			SaveMessage(gomock.Eq(expectedInput)).
			Return(expected, nil).
			Times(1)

		msg, err := usecase.SendMessage(input)

		assert.NoError(t, err)
		assert.NotNil(t, msg)
		assert.Equal(t, true, msg.IsAnon)
	})

	t.Run("repository_error", func(t *testing.T) {
		input := domain.Message{
			ID:        1,
			UserID:    1,
			Username:  "testuser",
			Content:   "lol",
			CreatedAt: time.Now(),
			Role:      "user",
			IsAnon:    false,
		}

		mockRepo.EXPECT().
			SaveMessage(gomock.Eq(input)).
			Return(domain.Message{}, errors.New("repo_error")).
			Times(1)

		msg, err := usecase.SendMessage(input)

		assert.Error(t, err)
		assert.Equal(t, domain.Message{}, msg)

	})

}

func TestUseCase_GetMessages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMessageRepository(ctrl)
	mockProducer := mocks.NewMockProducerEvents(ctrl)
	usecase := NewUseCase(mockRepo, mockProducer)

	t.Run("success", func(t *testing.T) {
		var limit int64 = 0
		var offset int64 = 5

		expected := []domain.Message{
			{ID: 1, UserID: 1, Username: "testuser", Content: "lol", CreatedAt: time.Now(), Role: "user", IsAnon: false},
			{ID: 3, UserID: 1, Username: "testuser", Content: "lol", CreatedAt: time.Now(), Role: "user", IsAnon: false},
			{ID: 2, UserID: 1, Username: "testuser", Content: "lol", CreatedAt: time.Now(), Role: "user", IsAnon: false},
		}

		mockRepo.EXPECT().
			MessageHistory(gomock.Eq(limit), gomock.Eq(offset)).
			Return(expected, nil).
			Times(1)

		msgs, err := usecase.GetMessages(int64(limit), int64(offset))

		assert.NoError(t, err)
		assert.NotNil(t, msgs)
		assert.Equal(t, expected, msgs)
	})

	t.Run("repository_error", func(t *testing.T) {
		var limit int64 = 0
		var offset int64 = 5

		mockRepo.EXPECT().
			MessageHistory(gomock.Eq(limit), gomock.Eq(offset)).
			Return([]domain.Message{}, errors.New("repo_error")).
			Times(1)

		msgs, err := usecase.GetMessages(limit, offset)

		assert.Error(t, err)
		assert.Equal(t, []domain.Message{}, msgs)
	})
}

func TestUseCase_DeleteMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMessageRepository(ctrl)
	mockProducer := mocks.NewMockProducerEvents(ctrl)
	usecase := NewUseCase(mockRepo, mockProducer)

	t.Run("success", func(t *testing.T) {
		var id int64 = 1
		var userID int64 = 10

		mockProducer.EXPECT().
			SendDeleteMessageEvent(gomock.Any(), gomock.Eq(userID), gomock.Eq(id)).
			Return(nil).
			Times(1)

		mockRepo.EXPECT().
			DeleteMessage(gomock.Eq(id)).
			Return(nil).
			Times(1)

		err := usecase.DeleteMessage(id, userID)

		assert.NoError(t, err)
	})

	t.Run("producer_error", func(t *testing.T) {
		var id int64 = 1
		var userID int64 = 10

		mockProducer.EXPECT().
			SendDeleteMessageEvent(gomock.Any(), gomock.Eq(userID), gomock.Eq(id)).
			Return(errors.New("producer error")).
			Times(1)

		err := usecase.DeleteMessage(id, userID)

		assert.Error(t, err)
	})

	t.Run("repository_error", func(t *testing.T) {
		var id int64 = 1
		var userID int64 = 10

		mockProducer.EXPECT().
			SendDeleteMessageEvent(gomock.Any(), gomock.Eq(userID), gomock.Eq(id)).
			Return(nil).
			Times(1)

		mockRepo.EXPECT().
			DeleteMessage(gomock.Eq(id)).
			Return(errors.New("repo_error")).
			Times(1)

		err := usecase.DeleteMessage(id, userID)

		assert.Error(t, err)
	})
}
