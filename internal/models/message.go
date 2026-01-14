package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	SenderID   uuid.UUID `gorm:"not null;index"`
	ReceiverID uuid.UUID `gorm:"not null;index"`

	Content   string `gorm:"type:text;not null"`
	IsDeleted bool   `gorm:"default:false"`
	EditedAt  *time.Time

	CreatedAt time.Time
}
