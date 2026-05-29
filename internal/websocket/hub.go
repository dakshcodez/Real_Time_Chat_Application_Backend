package websocket

import (
	"encoding/json"
	//"time"
)

type Hub struct {
	users      map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client

	messageService *MessageService
}

func NewHub(messageService *MessageService) *Hub {
	return &Hub{
		users:           make(map[string]map[*Client]bool),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		messageService:  messageService,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			firstConnection := false
			if h.users[c.UserID] == nil {
				h.users[c.UserID] = make(map[*Client]bool)
				firstConnection = true
			}
			h.users[c.UserID][c] = true

			if firstConnection {
				h.broadcastPresence(c.UserID, true)
			}
			h.sendOnlineUsersList(c)

		case c := <-h.unregister:
			if conns, ok := h.users[c.UserID]; ok {
				delete(conns, c)
				close(c.Send)
				if len(conns) == 0 {
					delete(h.users, c.UserID)
					h.broadcastPresence(c.UserID, false)
				}
			}
		}
	}
}

func (h *Hub) broadcastPresence(userID string, online bool) {
	event := map[string]any{
		"type":    "presence_change",
		"user_id": userID,
		"online":  online,
	}
	data, _ := json.Marshal(event)
	h.BroadcastToAll(data, userID)
}

func (h *Hub) sendOnlineUsersList(client *Client) {
	var list []string
	for uid := range h.users {
		list = append(list, uid)
	}
	event := map[string]any{
		"type":  "online_users",
		"users": list,
	}
	data, _ := json.Marshal(event)
	client.Send <- data
}

func (h *Hub) BroadcastToAll(data []byte, exceptUserID string) {
	for uid, conns := range h.users {
		if uid == exceptUserID {
			continue
		}
		for c := range conns {
			select {
			case c.Send <- data:
			default:
				// Avoid blocking
			}
		}
	}
}

func (h *Hub) routeMessage(sender *Client, raw []byte) {
	var msg IncomingMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		return
	}

	if msg.Type != "direct_message" {
		return
	}

	//Persist message
	saved, err := h.messageService.SaveMessage(
		sender.UserID,
		msg.To,
		msg.Content,
	)
	if err != nil {
		return
	}

	//Build outgoing message
	out := OutgoingMessage{
		Type:           "direct_message",
		ID:             saved.ID.String(),
		From:           sender.UserID,
		To:             msg.To,
		Content:        saved.Content,
		Timestamp:      saved.CreatedAt.Unix(),
		SenderUsername: sender.Username,
	}

	data, _ := json.Marshal(out)

	//Deliver to both sender and receiver
	h.BroadcastToUsers([]string{sender.UserID, msg.To}, data)
}

func (h *Hub) BroadcastToUsers(userIDs []string, data []byte) {
	for _, uid := range userIDs {
		if conns, ok := h.users[uid]; ok {
			for c := range conns {
				select {
				case c.Send <- data:
				default:
					// Avoid blocking
				}
			}
		}
	}
}
