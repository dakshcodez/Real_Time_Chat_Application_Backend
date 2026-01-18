package websocket

import (
	"time"

	"github.com/dakshcodez/real_time_chat_application_backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func FetchChatHistory(
	db *gorm.DB,
	userA uuid.UUID,
	userB uuid.UUID,
	limit int,
	before *time.Time,
) ([]models.Message, error) {

	query := db.
		Where(
			"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userA, userB, userB, userA,
		).
		Where("is_deleted = FALSE").
		Order("created_at DESC").
		Limit(limit)

	if before != nil {
		query = query.Where("created_at < ?", *before)
	}

	var messages []models.Message
	err := query.Find(&messages).Error

	return messages, err
}
