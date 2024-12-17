package models

import (
	"time"

	"github.com/albertolionardi/app/pkg/config"

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

func DeleteComment(ID int64) Comment {
	var comment Comment
	db.Where("Comment=?", comment).Delete(comment)
	return comment
}

func (c *Comment) CreateComment() *Comment {
	db.NewRecord(c)
	db.Create(&c)
	return c
}
