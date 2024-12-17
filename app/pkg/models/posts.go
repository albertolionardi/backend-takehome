package models

import (
	"app/pkg/config"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

type Post struct {
	gorm.Model
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorID  int       `json:"author_id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Comments  []Comment `json:"comments"`
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&Post{})
}

func (p *Post) CreatePost() *Post {
	db.NewRecord(p)
	db.Create(&p)
	return p
}

func GetPosts() []Post {
	var Posts []Post
	db.Find(&Posts)
	return Posts
}

func GetPostById(ID int64) (*Post, *gorm.DB) {
	var getPost Post
	db := db.Where("ID=?", ID).Find(&getPost)
	return &getPost, db
}

func DeletePost(ID int64) Post {
	var post Post
	db.Where("ID=?", ID).Delete(post)
	return post
}

// Update Post ?
func UpdatePost(ID int64) Post {
	var post Post
	db.Where("ID=?", ID).Update(post)
	return post
}
