package models

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint      `gorm:"primarykey"`
	Name         string    `json:"name"`
	Email        string    `json:"email" gorm:"unique"`
	PasswordHash string    `json:"password"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RegisterUser struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password"`
}
type LoginUser struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password"`
}

func MigrateUser(db *gorm.DB) {
	db.AutoMigrate(&User{})
	log.Println("User Table Migration Completed...")
}
