package usecase

import (
	"chat_service/internal/domain"
	"fmt"
	"log"
)

type UseCase struct {
	repo domain.MessageRepository
}

func NewUseCase(repo domain.MessageRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) SendMessage(msg domain.Message) (domain.Message, error) {
	if msg.Content == "" {
		return domain.Message{}, fmt.Errorf("error message can't be empty")
	}
	if msg.UserID == 0 {
		return domain.Message{}, fmt.Errorf("user id cannot be zero")
	}
	if msg.Role == "anon" {
		msg.IsAnon = true
	}

	saveMsg, err := uc.repo.SaveMessage(msg)
	if err != nil {
		log.Printf("Failed to send message:%v", err)
	}

	return saveMsg, nil
}

func (uc *UseCase) GetMessages(limit, offset int64) ([]domain.Message, error) {
	if limit <= 0 {
		limit = 50
	}

	if offset < 0 {
		offset = 0
	}

	return uc.repo.MessageHistory(limit, offset)
}

func (uc *UseCase) DeleteMessage(id int64) error {
	return uc.repo.DeleteMessage(id)
}
