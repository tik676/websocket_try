package infrastructure

import (
	"chat_service/internal/domain"
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{writer: writer}
}

func (p *Producer) SendDeleteMessageEvent(ctx context.Context, userID, messageID int64) error {
	event := domain.KafkaProducerDeleteEvent{
		UserID:    userID,
		MessageID: messageID,
		Timestamp: time.Now(),
	}

	return p.SendEvent(ctx, "message-deleted", event)
}

func (p *Producer) SendEvent(ctx context.Context, eventType string, event interface{}) error {
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

func (p *Producer) Close() error {
	return p.writer.Close()
}
