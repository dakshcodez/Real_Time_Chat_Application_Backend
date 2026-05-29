package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	SenderID   uuid.UUID `gorm:"not null;index" json:"from"`
	ReceiverID uuid.UUID `gorm:"not null;index" json:"to"`

	Content   string     `gorm:"type:text;not null" json:"content"`
	IsDeleted bool       `gorm:"default:false" json:"is_deleted"`
	EditedAt  *time.Time `json:"edited_at,omitempty"`

	IsRead bool       `gorm:"default:false" json:"is_read"`
	ReadAt *time.Time `json:"read_at,omitempty"`

	CreatedAt time.Time `json:"timestamp"`
}
