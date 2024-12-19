package test

import (
	"app/middleware"
	"app/pkg/models"
	"app/pkg/routes"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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

func helperLoginAndGetSessionID(t *testing.T, db *gorm.DB) string {
	//This is a valid test account, no need to change anything
	loginRequest := models.User{
		Name:         "Alberto",
		Email:        "albertolionardi1@gmail.com",
		PasswordHash: "testpassword",
	}
	requestBody, _ := json.Marshal(loginRequest)
	fmt.Printf("%s", requestBody)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	controller := routes.NewController(db)
	handler := http.HandlerFunc(controller.Login)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	sessionID := response["data"].(map[string]interface{})["ID"].(string)
	fmt.Printf("%s", sessionID)

	if sessionID == "" {
		t.Fatalf("Session ID is empty")
	}
	return sessionID
}

func Test_Register(t *testing.T) {
	db := setupTestDB()
	controller := routes.NewController(db)
	user := models.RegisterUser{
		Name:         "PutAnyNameHere",
		Email:        "PutAnyEmailhere",
		PasswordHash: "PutAnyPasswordhere",
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

func Test_Login(t *testing.T) {
	db := setupTestDB()
	controller := routes.NewController(db)
	loginUser := models.User{
		Name:         "Alberto",
		Email:        "albertolionardi1@gmail.com",
		PasswordHash: "testpassword",
	}
	requestBody, _ := json.Marshal(loginUser)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.Login)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	assert.Equal(t, "User login successfully", response["message"])
}

func Test_CreatePost(t *testing.T) {
	db := setupTestDB()
	controller := routes.NewController(db)

	mux := http.NewServeMux()
	mux.Handle("/createPost", middleware.SessionMiddleware(db)(http.HandlerFunc(controller.CreatePost)))

	newPost := models.CreatePost{
		Title:   "Test Post",
		Content: "This is a test post content",
	}

	sessionID := helperLoginAndGetSessionID(t, db)

	postRequestBody, _ := json.Marshal(newPost)
	req := httptest.NewRequest("POST", "/createPost", bytes.NewBuffer(postRequestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("SessionID", sessionID)

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]string
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	assert.Equal(t, "Post successfully created", response["message"])
}
