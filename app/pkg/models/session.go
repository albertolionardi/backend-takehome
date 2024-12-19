package models

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type Session struct {
	ID        string    `gorm:"primarykey"`
	UserID    uint      `json:"user_id"`
	ExpiredAt time.Time `json:"expired_at"`
	CreatedAt time.Time `json:"created_at"`
}

func MigrateSession(db *gorm.DB) {
	db.AutoMigrate(&Session{})
	log.Println("Session Table Migration Completed...")
}
