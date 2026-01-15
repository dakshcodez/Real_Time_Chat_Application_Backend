package auth

import (
	"errors"
	"strings"

	"github.com/dakshcodez/gdg_chat_app_backend_task/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(db *gorm.DB, username, email, password string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	password = strings.TrimSpace(password)

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user := models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
	}

	return db.Create(&user).Error
}

func Login(db *gorm.DB, email, password string) (*models.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	password = strings.TrimSpace(password)

	var user models.User
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}

