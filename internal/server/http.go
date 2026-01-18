package server

import (
	"net/http"

	"github.com/dakshcodez/gdg_chat_app_backend_task/internal/auth"
	"github.com/dakshcodez/gdg_chat_app_backend_task/internal/middleware"
	"github.com/dakshcodez/gdg_chat_app_backend_task/internal/websocket"
	"gorm.io/gorm"
)

func RegisterRoutes(mux *http.ServeMux, db *gorm.DB, jwtSecret string) {
	msgService := &websocket.MessageService{
		DB: db,
	}

	authHandler := &auth.Handler{
		DB: db,
		Secret: jwtSecret,
	}

	userHandler := &UserHandler{
		DB: db,
	}

	chatHandler := &ChatHandler{
		DB: db,
	}

	messageHandler := &MessageHandler{
		Service: msgService,
	}

	mux.HandleFunc("/auth/register", authHandler.Register)
	mux.HandleFunc("/auth/login", authHandler.Login)

	protected := middleware.JWTAuth(jwtSecret)

	mux.Handle(
		"/users/me",
		protected(http.HandlerFunc(userHandler.Me)),
	)

	mux.Handle(
		"/users/me/update",
		protected(http.HandlerFunc(userHandler.UpdateMe)),
	)

	mux.Handle(
		"/chats/{userId}",
		protected(http.HandlerFunc(chatHandler.History)),
	)

	mux.Handle(
		"/messages/{messageId}",
		protected(http.HandlerFunc(messageHandler.Edit)),
	)

	mux.Handle(
		"/messages/{messageId}/delete",
		protected(http.HandlerFunc(messageHandler.Delete)),
	)

	hub := websocket.NewHub(msgService)
	go hub.Run()

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWS(hub, jwtSecret, w, r)
	})
}