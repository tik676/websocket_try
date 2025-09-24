package repositories

import (
	"context"
	"notification-service/internal/domain/entities"
)

type FileRepository interface {
	SaveMessage(ctx context.Context, notification *entities.Notification) error
}
