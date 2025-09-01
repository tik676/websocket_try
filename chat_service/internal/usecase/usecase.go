package usecase

import (
	"chat_service/internal/domain"
	"encoding/json"
	"errors"
	"time"
)

type UseCase struct {
	repo        domain.MessageRepository
	broadcaster domain.Broadcaster
}

func NewUseCase(repo domain.MessageRepository, b domain.Broadcaster) *UseCase {
	return &UseCase{
		repo:        repo,
		broadcaster: b,
	}
}

func (uc *UseCase) SendMessage(user domain.User, content string) (domain.Message, error) {
	if content == "" {
		return domain.Message{}, errors.New("message content cannot be empty")
	}

	msg := domain.Message{
		UserID:    user.ID,
		UserName:  user.Name,
		Role:      user.Role,
		IsAnon:    user.IsAnon,
		Content:   content,
		CreatedAt: time.Now(),
	}

	saveMsg, err := uc.repo.SendMessage(msg)
	if err != nil {
		return domain.Message{}, err
	}

	if uc.broadcaster != nil {
		msgJSON, _ := json.Marshal(saveMsg)
		uc.broadcaster.Broadcast(msgJSON)
	}

	return saveMsg, nil
}

func (uc *UseCase) DeleteMessage(id int64) error {
	return uc.repo.DeleteMessage(id)
}

func (uc *UseCase) GetMessages(limit, offset int) ([]domain.Message, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	if offset < 0 {
		offset = 0
	}

	return uc.repo.GetMessages(limit, offset)
}
