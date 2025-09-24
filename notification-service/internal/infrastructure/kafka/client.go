package kafka

import (
	"notification-service/internal/usecase"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	processor *usecase.NotificationProcessor
	reader    *kafka.Reader
}

func NewConsumer(processor *usecase.NotificationProcessor) *Consumer {
	reader := &kafka.ReaderConfig{}
	return &Consumer{
		processor: processor,
		reader:    reader,
	}
}
