package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/dakshcodez/gdg_chat_app_backend_task/internal/middleware"
	"github.com/dakshcodez/gdg_chat_app_backend_task/internal/websocket"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatHandler struct {
	DB *gorm.DB
}

func (h *ChatHandler) History(w http.ResponseWriter, r *http.Request) {
	// Authenticated user
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)

	// Extract other user ID from URL
	otherIDStr := r.PathValue("userId")
	otherID, err := uuid.Parse(otherIDStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	// Pagination params
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}

	var before *time.Time
	if b := r.URL.Query().Get("before"); b != "" {
		if ts, err := strconv.ParseInt(b, 10, 64); err == nil {
			t := time.Unix(ts, 0)
			before = &t
		}
	}

	// Fetch messages
	messages, err := websocket.FetchChatHistory(
		h.DB,
		userID,
		otherID,
		limit,
		before,
	)
	if err != nil {
		http.Error(w, "failed to fetch messages", http.StatusInternalServerError)
		return
	}

	// 5️⃣ Build response
	type MessageResponse struct {
		ID        uuid.UUID `json:"id"`
		From      uuid.UUID `json:"from"`
		To        uuid.UUID `json:"to"`
		Content   string    `json:"content"`
		Timestamp time.Time `json:"timestamp"`
		EditedAt  *time.Time `json:"edited_at,omitempty"`
	}

	resp := make([]MessageResponse, 0, len(messages))
	for _, m := range messages {
		resp = append(resp, MessageResponse{
			ID:        m.ID,
			From:      m.SenderID,
			To:        m.ReceiverID,
			Content:   m.Content,
			Timestamp: m.CreatedAt,
			EditedAt:  m.EditedAt,
		})
	}

	json.NewEncoder(w).Encode(resp)
}
