package models

import (
	"app/pkg/config"
	"time"

	"github.com/jinzhu/gorm"
)

type Comment struct {
	gorm.Model
	PostID     uint      `json:"post_id"`
	AuthorName string    `json:"author_name"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&Comment{})
}
