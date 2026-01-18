package websocket

import (
	"net/http"
	"time"

	"github.com/dakshcodez/gdg_chat_app_backend_task/internal/auth"
	"github.com/dakshcodez/gdg_chat_app_backend_task/internal/ratelimit"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // tighten later
	},
}

func ServeWS(hub *Hub, secret string, w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	userID, err := auth.ParseJWT(token, secret)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	msgLimiter := ratelimit.New(10, time.Second)

	client := &Client{
		UserID: userID.String(),
		Conn:   conn,
		Send:   make(chan []byte, 256),
		Hub:    hub,
		Limiter: msgLimiter,
	}

	hub.register <- client

	go client.writePump()
	go client.readPump()
}