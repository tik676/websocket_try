package infrastructure

import (
	"chat_service/internal/domain"
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

type producer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(writer *kafka.Writer) *producer {
	return &producer{writer: writer}
}

func (p *producer) SendDeleteMessageEvent(ctx context.Context, userID, messageID int64) error {
	event := domain.KafkaProducerDeleteEvent{
		UserID:    userID,
		MessageID: messageID,
		Timestamp: time.Now(),
	}

	return p.SendEvent(ctx, "message-deleted", event)
}

func (p *producer) SendEvent(ctx context.Context, eventType string, event interface{}) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Key:   []byte(eventType),
		Value: data,
		Time:  time.Now(),
	}

	return p.writer.WriteMessages(ctx, message)
}

func (p *producer) Close() error {
	return p.writer.Close()
}
