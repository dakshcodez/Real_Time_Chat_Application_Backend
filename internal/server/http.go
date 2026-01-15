package server

import (
	"net/http"

	"github.com/dakshcodez/gdg_chat_app_backend_task/internal/auth"
	"gorm.io/gorm"
)

func RegisterRoutes(mux *http.ServeMux, db *gorm.DB, jwtSecret string) {
	authHandler := &auth.Handler{
		DB: db,
		Secret: jwtSecret,
	}

	mux.HandleFunc("/auth/register", authHandler.Register)
	mux.HandleFunc("/auth/login", authHandler.Login)
}