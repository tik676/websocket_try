package entities

import "time"

type Notification struct {
	ID        int64
	TopicName string
	Message   []byte
	Timestamp time.Time
}
