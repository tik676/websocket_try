package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Kafka KafkaConfig
	File  FileConfig
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

type FileConfig struct {
	LogPath string
}

func Load() (*Config, error) {
	cfg := &Config{
		Kafka: KafkaConfig{
			Brokers: getStringSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
			Topic:   getEnvString("KAFKA_TOPIC", "user-events"),
			GroupID: getEnvString("KAFKA_GROUP_ID", "notification-consumer"),
		},
		File: FileConfig{
			LogPath: getEnvString("LOG_PATH", "./logs"),
		},
	}

	return cfg, cfg.validate()
}

func (c *Config) validate() error {
	if len(c.Kafka.Brokers) == 0 {
		return fmt.Errorf("kafka brokers cannot be empty")
	}
	if len(c.Kafka.Topic) == 0 {
		return fmt.Errorf("kafka topic cannot be empty")
	}
	if len(c.Kafka.GroupID) == 0 {
		return fmt.Errorf("kafka group ID cannot be empty")
	}
	return nil
}

func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getStringSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}
