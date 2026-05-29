package server

import (
	"net/http"
	"time"

	"github.com/dakshcodez/real_time_chat_application_backend/internal/auth"
	"github.com/dakshcodez/real_time_chat_application_backend/internal/middleware"
	"github.com/dakshcodez/real_time_chat_application_backend/internal/ratelimit"
	"github.com/dakshcodez/real_time_chat_application_backend/internal/websocket"
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

	hub := websocket.NewHub(msgService)
	go hub.Run()

	messageHandler := &MessageHandler{
		Service: msgService,
		Hub:	 hub,
	}

	mux.HandleFunc("/auth/register", authHandler.Register)
	mux.HandleFunc("/auth/login", authHandler.Login)

	protected := middleware.JWTAuth(jwtSecret)
	restLimiter := ratelimit.New(60, time.Minute)
	rateLimit := middleware.RateLimit(restLimiter)

	mux.Handle(
		"/users/me",
		protected(rateLimit(http.HandlerFunc(userHandler.Me))),
	)

	mux.Handle(
		"/users/me/update",
		protected(rateLimit(http.HandlerFunc(userHandler.UpdateMe))),
	)

	mux.Handle(
		"/chats/{userId}",
		protected(rateLimit(http.HandlerFunc(chatHandler.History))),
	)

	mux.Handle(
		"/messages/{messageId}",
		protected(rateLimit(http.HandlerFunc(messageHandler.Edit))),
	)

	mux.Handle(
		"/messages/{messageId}/delete",
		protected(rateLimit(http.HandlerFunc(messageHandler.Delete))),
	)

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWS(hub, jwtSecret, w, r)
	})
}