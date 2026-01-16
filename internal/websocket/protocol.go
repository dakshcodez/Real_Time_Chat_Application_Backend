package websocket

type IncomingMessage struct {
	Type    string `json:"type"`    // "direct_message"
	To      string `json:"to"`      // receiver user_id
	Content string `json:"content"` // message text
}

type OutgoingMessage struct {
	Type      string `json:"type"` // "direct_message"
	From      string `json:"from"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}
