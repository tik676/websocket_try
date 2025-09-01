package main

import (
	"chat_service/internal/delivery/router"
	websocket "chat_service/internal/delivery/webSocket"
	"chat_service/internal/infrastructure/postgres"
	"chat_service/internal/middleware"
	"chat_service/internal/usecase"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_NAME")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}
	log.Println("Successfully connected to database")

	cm := websocket.NewClientManager()
	postgresRepo := postgres.NewDB(db)
	usecase := usecase.NewUseCase(postgresRepo, cm.GetBroadcaster())
	jwtKey := middleware.NewJWTMaker(os.Getenv("JWT_SECRET"))

	router := router.SetupRouter(usecase, jwtKey)
	router.Run(":8080")
}
