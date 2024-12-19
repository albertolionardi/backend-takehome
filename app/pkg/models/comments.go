package models

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID         uint      `gorm:"primarykey"`
	PostID     uint      `json:"post_id"`
	AuthorName string    `json:"author_name"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

func MigrateComment(db *gorm.DB) {
	db.AutoMigrate(&Comment{})
	log.Println("Comment Table Migration Completed...")
}
