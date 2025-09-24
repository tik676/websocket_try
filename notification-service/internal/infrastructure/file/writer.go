package file

import (
	"context"
	"fmt"
	"notification-service/internal/domain/entities"
	"os"
	"path/filepath"
	"time"
)

type FileWriter struct {
	path string
}

func (fw *FileWriter) SaveMessage(ctx context.Context, notification *entities.Notification) error {
	fileName := fmt.Sprintf("%s_%s.log", notification.TopicName, time.Now().Format("2025-09-24"))
	return os.WriteFile(filepath.Join(fw.path, fileName), notification.Message, 0644)
}
