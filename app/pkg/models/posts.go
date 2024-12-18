package models

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string `json:"title"`
	Content   string `json:"content"`
	AuthorID  uint   `json:"author_id"`
}

type CreatePost struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdatePost struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	UpdatedAt time.Time
}

func MigratePost(db *gorm.DB) {
	db.AutoMigrate(&Post{})
	log.Println("Post Table Migration Completed...")
}
