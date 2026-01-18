package main

import (
	"log"
	"net/http"

	"github.com/dakshcodez/real_time_chat_application_backend/internal/config"
	"github.com/dakshcodez/real_time_chat_application_backend/internal/db"
	"github.com/dakshcodez/real_time_chat_application_backend/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	cfg := config.Load()
	dbConn := db.Connect(cfg.DBUrl)

	mux := http.NewServeMux()
	server.RegisterRoutes(mux, dbConn, cfg.JWTSecret)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}