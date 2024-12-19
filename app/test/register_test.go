package test

import (
	"app/pkg/models"
	"app/pkg/routes"
	"app/pkg/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	dsn := "root:abc123@tcp(127.0.0.1:3333)/appdb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
	}
	db.AutoMigrate(&models.User{})
	return db
}

func Test_Register(t *testing.T) {
	db := setupTestDB()
	controller := routes.NewController(db)
	hashpasssword, _ := utils.HashPassword("password123")
	user := models.RegisterUser{
		Name:         "John Doe",
		Email:        "john@example.com",
		PasswordHash: hashpasssword,
	}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	controller.Register(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var respBody map[string]string
	json.NewDecoder(resp.Body).Decode(&respBody)
	if respBody["message"] != "User registered successfully" {
		t.Fatalf("Expected message %q, got %q", "User registered successfully", respBody["message"])
	}
}
