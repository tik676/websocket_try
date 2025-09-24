package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"user_service/internal/delivery/rest"
	"user_service/internal/infrastructure"
	"user_service/internal/usecase"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		log.Fatal("Invalid POSTGRES_PORT:", err)
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		port,
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_NAME"),
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Failed to close database connection: %v", err)
		}
	}()

	for i := 0; i < 10; i++ {
		err := db.Ping()
		if err == nil {
			break
		}
		log.Printf("Waiting for database, attempt %d...", i+1)
		time.Sleep(2 * time.Second)
	}

	brokersStr := os.Getenv("KAFKA_BROKERS")
	if brokersStr == "" {
		brokersStr = "localhost:9092"
	}

	brokers := strings.Split(brokersStr, ",")

	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "user-events"
	}

	producer := infrastructure.NewKafkaProducer(brokers, topic)
	defer func() {
		if err := producer.Close(); err != nil {
			log.Printf("Error closing Kafka producer: %v", err)
		}
	}()
	repo := infrastructure.NewDB(db)
	jwtKey := infrastructure.NewJWTMaker(os.Getenv("JWT_SECRET"), db)
	usecase := usecase.NewUseCase(repo, jwtKey, producer)

	router := rest.SetupRouter(usecase, jwtKey)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Printf("Failed to up server:%v", err)
	}
}
