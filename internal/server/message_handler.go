package server

import (
	"encoding/json"
	"net/http"

	"github.com/dakshcodez/gdg_chat_app_backend_task/internal/middleware"
	"github.com/dakshcodez/gdg_chat_app_backend_task/internal/websocket"
	"github.com/google/uuid"
)

type MessageHandler struct {
	Service *websocket.MessageService
	Hub		*websocket.Hub
}

func (h *MessageHandler) Edit(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)

	messageID, err := uuid.Parse(r.PathValue("messageId"))
	if err != nil {
		http.Error(w, "invalid message id", http.StatusBadRequest)
		return
	}

	var body struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Content == "" {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	// Update DB
	msg, err := h.Service.EditMessage(messageID, userID, body.Content)
	if err != nil {
		http.Error(w, "not allowed", http.StatusForbidden)
		return
	}

	// Build WS event
	event := websocket.OutgoingMessage{
		Type:    "message_edited",
		ID:      msg.ID.String(),
		Content: msg.Content,
	}

	data, _ := json.Marshal(event)

	// Broadcast to both users
	h.Hub.BroadcastToUsers(
		[]string{
			msg.SenderID.String(),
			msg.ReceiverID.String(),
		},
		data,
	)

	// REST response
	json.NewEncoder(w).Encode(msg)
}

func (h *MessageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)

	messageID, err := uuid.Parse(r.PathValue("messageId"))
	if err != nil {
		http.Error(w, "invalid message id", http.StatusBadRequest)
		return
	}

	// Soft delete in DB
	msg, err := h.Service.DeleteMessage(messageID, userID)
	if err != nil {
		http.Error(w, "not allowed", http.StatusForbidden)
		return
	}

	// Build delete event
	event := websocket.OutgoingMessage{
		Type: "message_deleted",
		ID:   msg.ID.String(),
	}

	data, _ := json.Marshal(event)

	// Broadcast
	h.Hub.BroadcastToUsers(
		[]string{
			msg.SenderID.String(),
			msg.ReceiverID.String(),
		},
		data,
	)

	// REST response
	json.NewEncoder(w).Encode(msg)
}
