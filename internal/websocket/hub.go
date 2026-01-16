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
			if h.users[c.UserID] == nil {
				h.users[c.UserID] = make(map[*Client]bool)
			}
			h.users[c.UserID][c] = true

		case c := <-h.unregister:
			if conns, ok := h.users[c.UserID]; ok {
				delete(conns, c)
				close(c.Send)
				if len(conns) == 0 {
					delete(h.users, c.UserID)
				}
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
		Type:      "direct_message",
		From:      sender.UserID,
		Content:   saved.Content,
		Timestamp: saved.CreatedAt.Unix(),
	}

	data, _ := json.Marshal(out)

	//Deliver if receiver is online
	if receivers, ok := h.users[msg.To]; ok {
		for c := range receivers {
			c.Send <- data
		}
	}
}
