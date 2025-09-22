package usecase

import (
	"chat_service/internal/domain"
	"context"
	"fmt"
	"log"
	"time"
)

type UseCase struct {
	repo     domain.MessageRepository
	producer domain.ProducerEvents
}

func NewUseCase(repo domain.MessageRepository, producer domain.ProducerEvents) *UseCase {
	return &UseCase{repo: repo, producer: producer}
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
		return domain.Message{}, err
	}

	return saveMsg, nil
}

func (uc *UseCase) GetMessages(limit, offset int64) ([]domain.Message, error) {
	if limit < 0 {
		limit = 50
	}

	if offset < 0 {
		offset = 0
	}

	msgs, err := uc.repo.MessageHistory(limit, offset)
	if err != nil {
		log.Printf("error: %v", err)
		return []domain.Message{}, err
	}

	return msgs, nil
}

func (uc *UseCase) DeleteMessage(id int64, userID int64) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := uc.producer.SendDeleteMessageEvent(ctx, userID, id)
	if err != nil {
		return err
	}

	return uc.repo.DeleteMessage(id)
}
