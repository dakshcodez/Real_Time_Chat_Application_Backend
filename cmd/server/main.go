package main

import (
	"log"
	"net/http"

	"github.com/dakshcodez/real_time_chat_application_backend/internal/config"
	"github.com/dakshcodez/real_time_chat_application_backend/internal/db"
	"github.com/dakshcodez/real_time_chat_application_backend/internal/middleware"
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

	// Wrap the mux with the CORS middleware
	handler := middleware.CORS(mux)

	log.Printf("Server running on :%s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, handler))
}