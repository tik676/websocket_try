package file

import (
	"context"
	"fmt"
	"log"
	"notification-service/internal/domain/entities"
	"os"
	"path/filepath"
	"time"
)

type FileWriter struct {
	path string
}

func NewFileWriter(path string) (*FileWriter, error) {
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	return &FileWriter{path: path}, nil
}

func (fw *FileWriter) SaveMessage(ctx context.Context, notification *entities.Notification) error {
	fileName := fmt.Sprintf("%s_%s.log", notification.TopicName, time.Now().Format("2006-01-02"))
	filePath := filepath.Join(fw.path, fileName)

	logEntry := fmt.Sprintf("[%s] ID:%d Topic:%s,EventType:%s Message:%s\n",
		notification.Timestamp.Format("2006-01-02 15:04:05"),
		notification.ID,
		notification.TopicName,
		notification.EventType,
		string(notification.Message),
	)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Printf("Failed to close log file: %s", err)
		}
	}(file)

	_, err = file.WriteString(logEntry)

	return err
}
