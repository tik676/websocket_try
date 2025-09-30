package kafka

import (
	"context"
	"log"
	"notification-service/internal/usecase"
	"time"

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
	defer func() {
		if err := c.reader.Close(); err != nil {
			log.Printf("Failed to close kafka reader %s", err)
		}
	}()
	for {
		select {
		case <-ctx.Done():
			log.Println("Consumer shutdown")
			return
		default:
			msgCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			msg, err := c.reader.ReadMessage(msgCtx)
			cancel()

			if err != nil {
				if err == context.DeadlineExceeded {
					continue
				}
				log.Printf("Error reading message: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			eventType := string(msg.Key)
			err = c.processor.ProcessMessage(ctx, eventType, msg.Value, msg.Topic)
			if err != nil {
				log.Printf("Error processing message: %v", err)
			}
		}
	}
}
