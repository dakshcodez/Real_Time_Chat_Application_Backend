package websocket

type IncomingMessage struct {
	Type    string `json:"type"`    // "direct_message"
	To      string `json:"to"`      // receiver user_id
	Content string `json:"content"` // message text
}

type OutgoingMessage struct {
	Type      string `json:"type"`               // event type
	ID        string `json:"id,omitempty"`       // message id
	From      string `json:"from,omitempty"`     // sender
	Content   string `json:"content,omitempty"`  // message text
	Timestamp int64  `json:"timestamp,omitempty"`
}
