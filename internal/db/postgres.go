package db

import (
	"log"
	"os"
	"time"

	"github.com/dakshcodez/real_time_chat_application_backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(dsn string) *gorm.DB {
	// Configure custom GORM logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             2 * time.Second, // Raise threshold to 2s to accommodate remote DB latency
			LogLevel:                  logger.Warn,     // Only log warnings and errors
			IgnoreRecordNotFoundError: true,            // Ignore normal record-not-found warnings
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
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