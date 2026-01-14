package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username     string    `gorm:"unique;not null"`
	Email        string    `gorm:"unique;not null"`
	PasswordHash string    `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
