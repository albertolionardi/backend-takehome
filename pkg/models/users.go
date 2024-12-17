package models

import (
	"app/pkg/config"
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name         string    `json:"name"`
	Email        string    `json:"email" gorm:"unique"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&User{})
}
