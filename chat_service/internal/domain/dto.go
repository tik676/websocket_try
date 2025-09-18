package domain

import (
	"context"
	"time"
)

//go:generate mockgen -source=dto.go -destination=mocks/dto.go -package=mockschat

type User struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	IsAnon bool   `json:"is_anon"`
}

type Message struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Role      string    `json:"role"`
	IsAnon    bool      `json:"is_anon"`
}

type KafkaProducerDeleteEvent struct {
	MessageID int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}

type ProducerEvents interface {
	SendDeleteMessageEvent(ctx context.Context, userID, message int64) error
	Close() error
}

type MessageRepository interface {
	SaveMessage(msg Message) (Message, error)
	MessageHistory(limit, offset int64) ([]Message, error)
	DeleteMessage(id int64) error
}

type TokenManager interface {
	VerifyToken(token string) (userID int64, name, role string, err error)
}
