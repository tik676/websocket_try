package infrastructure

import (
	"context"
	"encoding/json"
	"time"
	"user_service/internal/domain"

	"github.com/segmentio/kafka-go"
)

type kafkaWriter struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) domain.KafkaProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &kafkaWriter{writer: writer}
}

func (k *kafkaWriter) SendUserRegistered(ctx context.Context, userID int64, name string) error {
	event := domain.UserActionEvent{
		UserID:    userID,
		UserName:  name,
		Timestamp: time.Now(),
	}

	return k.sendEvent(ctx, "user-registered", event)
}

func (k *kafkaWriter) SendUserLoggedIn(ctx context.Context, userID int64, name string) error {
	event := domain.UserActionEvent{
		UserID:    userID,
		UserName:  name,
		Timestamp: time.Now(),
	}

	return k.sendEvent(ctx, "user-logged-in", event)
}

func (k *kafkaWriter) Close() error {
	return k.writer.Close()
}

func (k *kafkaWriter) sendEvent(ctx context.Context, eventType string, event interface{}) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Key:   []byte(eventType),
		Value: data,
		Time:  time.Now(),
	}

	return k.writer.WriteMessages(ctx, message)
}
