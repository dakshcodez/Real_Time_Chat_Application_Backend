package db

import (
	"log"

	"github.com/dakshcodez/gdg_chat_app_backend_task/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Message{},
	)
	if err != nil {
		log.Fatal("migration failed:", err)
	}

	return db
}