package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/dakshcodez/real_time_chat_application_backend/internal/middleware"
	"github.com/dakshcodez/real_time_chat_application_backend/internal/models"
	"github.com/dakshcodez/real_time_chat_application_backend/internal/websocket"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatHandler struct {
	DB  *gorm.DB
	Hub *websocket.Hub
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

	// Build response
	type MessageResponse struct {
		ID        uuid.UUID  `json:"id"`
		From      uuid.UUID  `json:"from"`
		To        uuid.UUID  `json:"to"`
		Content   string     `json:"content"`
		Timestamp time.Time  `json:"timestamp"`
		EditedAt  *time.Time `json:"edited_at,omitempty"`
		IsRead    bool       `json:"is_read"`
		ReadAt    *time.Time `json:"read_at,omitempty"`
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
			IsRead:    m.IsRead,
			ReadAt:    m.ReadAt,
		})
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *ChatHandler) Conversations(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)

	type ConvoRow struct {
		OtherID      uuid.UUID `gorm:"column:other_id"`
		LastActivity time.Time `gorm:"column:last_activity"`
	}

	var rows []ConvoRow
	err := h.DB.Raw(`
		SELECT other_id, MAX(created_at) as last_activity
		FROM (
			SELECT receiver_id as other_id, created_at FROM messages WHERE sender_id = ? AND is_deleted = FALSE
			UNION ALL
			SELECT sender_id as other_id, created_at FROM messages WHERE receiver_id = ? AND is_deleted = FALSE
		) sub
		GROUP BY other_id
		ORDER BY last_activity DESC
	`, userID, userID).Scan(&rows).Error

	if err != nil {
		http.Error(w, "failed to fetch conversations", http.StatusInternalServerError)
		return
	}

	type ConvoResponse struct {
		ID          uuid.UUID      `json:"id"`
		OtherUser   map[string]any `json:"other_user"`
		LastMessage map[string]any `json:"last_message,omitempty"`
		UnreadCount int            `json:"unread_count"`
	}

	response := make([]ConvoResponse, 0, len(rows))

	for _, row := range rows {
		var otherUser models.User
		if err := h.DB.First(&otherUser, "id = ?", row.OtherID).Error; err != nil {
			continue
		}

		var lastMsg models.Message
		h.DB.Order("created_at DESC").First(&lastMsg,
			"((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)) AND is_deleted = FALSE",
			userID, row.OtherID, row.OtherID, userID,
		)

		var unreadCount int64
		h.DB.Model(&models.Message{}).Where(
			"sender_id = ? AND receiver_id = ? AND is_read = ? AND is_deleted = FALSE",
			row.OtherID, userID, false,
		).Count(&unreadCount)

		var lastMsgMap map[string]any
		if lastMsg.ID != uuid.Nil {
			lastMsgMap = map[string]any{
				"id":         lastMsg.ID,
				"from":       lastMsg.SenderID,
				"to":         lastMsg.ReceiverID,
				"content":    lastMsg.Content,
				"timestamp":  lastMsg.CreatedAt,
				"edited_at":  lastMsg.EditedAt,
				"is_deleted": lastMsg.IsDeleted,
				"is_read":    lastMsg.IsRead,
			}
		}

		response = append(response, ConvoResponse{
			ID: row.OtherID,
			OtherUser: map[string]any{
				"id":       otherUser.ID,
				"username": otherUser.Username,
				"email":    otherUser.Email,
			},
			LastMessage: lastMsgMap,
			UnreadCount: int(unreadCount),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ChatHandler) CreateConversation(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)

	var body struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	otherID, err := uuid.Parse(body.UserID)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	var otherUser models.User
	if err := h.DB.First(&otherUser, "id = ?", otherID).Error; err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	var lastMsg models.Message
	h.DB.Order("created_at DESC").First(&lastMsg,
		"((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)) AND is_deleted = FALSE",
		userID, otherID, otherID, userID,
	)

	var unreadCount int64
	h.DB.Model(&models.Message{}).Where(
		"sender_id = ? AND receiver_id = ? AND is_read = ? AND is_deleted = FALSE",
		otherID, userID, false,
	).Count(&unreadCount)

	var lastMsgMap map[string]any
	if lastMsg.ID != uuid.Nil {
		lastMsgMap = map[string]any{
			"id":         lastMsg.ID,
			"from":       lastMsg.SenderID,
			"to":         lastMsg.ReceiverID,
			"content":    lastMsg.Content,
			"timestamp":  lastMsg.CreatedAt,
			"edited_at":  lastMsg.EditedAt,
			"is_deleted": lastMsg.IsDeleted,
			"is_read":    lastMsg.IsRead,
		}
	}

	response := map[string]any{
		"id": otherID,
		"other_user": map[string]any{
			"id":       otherUser.ID,
			"username": otherUser.Username,
			"email":    otherUser.Email,
		},
		"last_message": lastMsgMap,
		"unread_count": unreadCount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ChatHandler) MarkConversationRead(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)

	otherIDStr := r.PathValue("userId")
	otherID, err := uuid.Parse(otherIDStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	now := time.Now()
	err = h.DB.Model(&models.Message{}).
		Where("sender_id = ? AND receiver_id = ? AND is_read = FALSE", otherID, userID).
		Updates(map[string]any{
			"is_read": true,
			"read_at": &now,
		}).Error

	if err != nil {
		http.Error(w, "failed to mark messages as read", http.StatusInternalServerError)
		return
	}

	// Broadcast conversation_read event to both users
	event := map[string]any{
		"type":      "conversation_read",
		"reader_id": userID.String(),
		"user_id":   otherID.String(),
	}
	data, _ := json.Marshal(event)
	h.Hub.BroadcastToUsers([]string{userID.String(), otherID.String()}, data)

	w.WriteHeader(http.StatusOK)
}
