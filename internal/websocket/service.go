package websocket

import (
	"time"

	"github.com/dakshcodez/gdg_chat_app_backend_task/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageService struct {
	DB *gorm.DB
}

func (s *MessageService) SaveMessage(
	senderID string,
	receiverID string,
	content string,
) (*models.Message, error) {

	msg := &models.Message{
		SenderID:   uuid.MustParse(senderID),
		ReceiverID: uuid.MustParse(receiverID),
		Content:    content,
		CreatedAt:  time.Now(),
	}
	
	if err := s.DB.Create(msg).Error; err != nil {
		return nil, err
	}

	return msg, nil
}
