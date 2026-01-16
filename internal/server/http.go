package server

import (
	"net/http"

	"github.com/dakshcodez/gdg_chat_app_backend_task/internal/auth"
	"github.com/dakshcodez/gdg_chat_app_backend_task/internal/middleware"
	"github.com/dakshcodez/gdg_chat_app_backend_task/internal/websocket"
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

	msgService := &websocket.MessageService{
	DB: db,
	}

	hub := websocket.NewHub(msgService)
	go hub.Run()

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWS(hub, jwtSecret, w, r)
	})
}