package entities

import "time"

type Notification struct {
	ID        int64
	TopicName string
	EventType string
	Message   []byte
	Timestamp time.Time
}
