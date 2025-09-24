package usecase

import (
	"context"
	"errors"
	"notification-service/internal/domain/entities"
	"notification-service/internal/domain/repositories"
	"sync/atomic"
	"time"
)

var idNotification int64

type NotificationProcessor struct {
	fileRepo  repositories.FileRepository
	idCounter int64
}

func NewUseCase(fileRepo repositories.FileRepository) *NotificationProcessor {
	return &NotificationProcessor{
		fileRepo:  fileRepo,
		idCounter: 0,
	}
}

func (np *NotificationProcessor) generateID() int64 {
	return atomic.AddInt64(&np.idCounter, 1)
}

func (np *NotificationProcessor) ProcessMessage(msg []byte, topic string) error {
	if len(msg) == 0 {
		return errors.New("empty message")
	}

	notification := &entities.Notification{
		ID:        np.generateID(),
		TopicName: topic,
		Message:   msg,
		Timestamp: time.Now(),
	}

	return np.fileRepo.SaveMessage(context.Background(), notification)
}
