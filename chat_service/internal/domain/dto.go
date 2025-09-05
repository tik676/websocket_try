package domain

import "time"

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

type MessageRepository interface {
	SaveMessage(msg Message) (Message, error)
	MessageHistory(limit, offset int64) ([]Message, error)
	DeleteMessage(id int64) error
}

type TokenManager interface {
	VerifyToken(token string) (userID int64, name, role string, err error)
}

type Broadcaster interface {
	Broadcast(msg []byte)
}
