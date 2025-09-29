package main

import (
	"context"
	"log"
	"notification-service/config"
	"notification-service/internal/infrastructure/file"
	"notification-service/internal/infrastructure/kafka"
	"notification-service/internal/usecase"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fileRepo, err := file.NewFileWriter(cfg.File.LogPath)
	if err != nil {
		log.Fatalf("Failed to create file writer: %v", err)
	}

	processor := usecase.NewUseCase(fileRepo)

	consumer := kafka.NewConsumer(processor, cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.GroupID)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Starting notification service for topic %s", cfg.Kafka.Topic)

	go func() {
		consumer.Start(ctx)
	}()

	<-signChan
	log.Println("Shutting down gracefully...")
	cancel()

	log.Println("Service stopped")
}
