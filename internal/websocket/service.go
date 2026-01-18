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

func (s *MessageService) EditMessage(
	messageID uuid.UUID,
	userID uuid.UUID,
	newContent string,
) (*models.Message, error) {

	var msg models.Message

	err := s.DB.First(&msg, "id = ? AND sender_id = ?", messageID, userID).Error
	if err != nil {
		return nil, err
	}

	now := time.Now()
	msg.Content = newContent
	msg.EditedAt = &now

	if err := s.DB.Save(&msg).Error; err != nil {
		return nil, err
	}

	return &msg, nil
}

func (s *MessageService) DeleteMessage(
	messageID uuid.UUID,
	userID uuid.UUID,
) (*models.Message, error) {

	var msg models.Message

	err := s.DB.First(&msg, "id = ? AND sender_id = ?", messageID, userID).Error
	if err != nil {
		return nil, err
	}

	msg.IsDeleted = true

	if err := s.DB.Save(&msg).Error; err != nil {
		return nil, err
	}

	return &msg, nil
}
