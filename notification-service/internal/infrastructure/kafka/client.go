package kafka

import (
	"context"
	"log"
	"notification-service/internal/usecase"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	processor *usecase.NotificationProcessor
	reader    *kafka.Reader
}

func NewConsumer(processor *usecase.NotificationProcessor, brokers []string, topic string, groupID string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MaxBytes: 10e6,
	})

	return &Consumer{
		processor: processor,
		reader:    reader,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	defer c.reader.Close()
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		err = c.processor.ProcessMessage(ctx, msg.Value, msg.Topic)
		if err != nil {
			log.Printf("Error processing message: %v", err)
		}
	}
}
