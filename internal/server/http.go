package server

import (
	"net/http"
	"time"

	"github.com/dakshcodez/gdg_chat_app_backend_task/internal/auth"
	"github.com/dakshcodez/gdg_chat_app_backend_task/internal/middleware"
	"github.com/dakshcodez/gdg_chat_app_backend_task/internal/ratelimit"
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
		rateLimit(protected(http.HandlerFunc(userHandler.Me))),
	)

	mux.Handle(
		"/users/me/update",
		rateLimit(protected(http.HandlerFunc(userHandler.UpdateMe))),
	)

	mux.Handle(
		"/chats/{userId}",
		rateLimit(protected(http.HandlerFunc(chatHandler.History))),
	)

	mux.Handle(
		"/messages/{messageId}",
		rateLimit(protected(http.HandlerFunc(messageHandler.Edit))),
	)

	mux.Handle(
		"/messages/{messageId}/delete",
		rateLimit(protected(http.HandlerFunc(messageHandler.Delete))),
	)

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWS(hub, jwtSecret, w, r)
	})
}