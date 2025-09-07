package main

import (
	"chat_service/internal/delivery/router"
	"chat_service/internal/delivery/ws"
	"chat_service/internal/infrastructure"
	"chat_service/internal/usecase"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_NAME"),
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Failed to connect to the DB: ", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Error ping DB: ", err)
	}

	log.Println("Connection to the DB was successful")
	repo := infrastructure.NewDBrepo(db)
	jwtSecret := os.Getenv("JWT_SECRET")
	tokenManager := infrastructure.NewJWTmaker(jwtSecret)
	uc := usecase.NewUseCase(repo)

	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	wsHandler := ws.NewWsHandler(uc, upgrader)

	r := router.SetupRouter(uc, tokenManager, wsHandler)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
