package server

import (
	"net/http"

	"github.com/dakshcodez/gdg_chat_app_backend_task/internal/auth"
	"github.com/dakshcodez/gdg_chat_app_backend_task/internal/middleware"
	"gorm.io/gorm"
)

func RegisterRoutes(mux *http.ServeMux, db *gorm.DB, jwtSecret string) {
	authHandler := &auth.Handler{
		DB: db,
		Secret: jwtSecret,
	}

	userHandler := &UserHandler{
		DB: db,
	}

	mux.HandleFunc("/auth/register", authHandler.Register)
	mux.HandleFunc("/auth/login", authHandler.Login)

	protected := middleware.JWTAuth(jwtSecret)

	mux.Handle(
		"/users/me",
		protected(http.HandlerFunc(userHandler.Me)),
	)
}