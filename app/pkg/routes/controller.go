package routes

import (
	"app/pkg/models"
	"app/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type controller struct {
	db *gorm.DB
}

func NewController(db *gorm.DB) *controller {
	return &controller{db: db}
}
func getUserIDFromContext(r *http.Request) uint {
	userID := r.Context().Value("userID")
	if id, ok := userID.(uint); ok {
		return id
	}
	return 0
}
func (c *controller) Register(w http.ResponseWriter, r *http.Request) {
	var user models.RegisterUser

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	if user.Email == "" || user.PasswordHash == "" {
		http.Error(w, "Email and Password are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	newUser := models.User{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: hashedPassword,
		UpdatedAt:    time.Now(),
	}
	if err := c.db.Create(&newUser).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})

}
func (c *controller) Login(w http.ResponseWriter, r *http.Request) {
	var input models.LoginUser
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	var user models.User
	if err := c.db.Table("users").Where("email=?", input.Email).First(&user).Error; err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if err := utils.CheckPassword(user.PasswordHash, input.PasswordHash); err != nil {
		http.Error(w, "Invalid email or passworD", http.StatusUnauthorized)
		return
	}

	var session models.Session
	if err := c.db.Where("user_id = ?", user.ID).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			session = models.Session{
				ID:        uuid.NewString(),
				UserID:    user.ID,
				CreatedAt: time.Now(),
				ExpiredAt: time.Now().Add(1 * time.Hour),
			}
			if err := c.db.Create(&session).Error; err != nil {
				http.Error(w, "Failed to create session", http.StatusInternalServerError)
				fmt.Println("Create session error:", err)
				return
			}
		} else {
			http.Error(w, "Failed to query session", http.StatusInternalServerError)
			fmt.Println("Query session error:", err)
			return
		}

	} else {
		// IF SESSION EXIST BUT EXPIRED; DELETE PREVIOUS AND CREATE A NEW ONE
		if session.ExpiredAt.Before(time.Now()) {
			session.ID = uuid.NewString()
			session.CreatedAt = time.Now()
			session.ExpiredAt = time.Now().Add(1 * time.Hour)
			if err := c.db.Save(&session).Error; err != nil {
				http.Error(w, "Failed to update session", http.StatusInternalServerError)
				fmt.Println("Update session error:", err)
				return
			}
		}

	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User login successfully",
		"data":    session,
	})
}

func (c *controller) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post models.CreatePost
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	newPost := models.Post{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Title:     post.Title,
		Content:   post.Content,
		AuthorID:  userID,
	}
	if err := c.db.Create(&newPost).Error; err != nil {
		http.Error(w, "Fail to create Post", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post successfully created"})
}
func (c *controller) GetPostById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	postID, err := strconv.Atoi(id)
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	var post models.Post
	if err := c.db.First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}
func (c *controller) GetPosts(w http.ResponseWriter, r *http.Request) {
	var posts []models.Post
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if err := c.db.Find(&posts).Error; err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}
func (c *controller) DeletePostById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	postID, err := strconv.Atoi(id)
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	var post models.Post
	if err := c.db.First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if post.AuthorID != userID {
		http.Error(w, "User is not authorized to delete this post", http.StatusUnauthorized)
		return
	}
	if err := c.db.Delete(&post).Error; err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post successfully deleted"})

}

func (c *controller) UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	postID, err := strconv.Atoi(id)
	userID := getUserIDFromContext(r)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var post models.Post
	if err := c.db.First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if post.AuthorID != userID {
		http.Error(w, "User is not authorized to update this post", http.StatusUnauthorized)
		return
	}
	var updateData models.UpdatePost
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	post.Title = updateData.Title
	post.Content = updateData.Content
	post.UpdatedAt = time.Now()
	if err := c.db.Save(&post).Error; err != nil {
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post successfully updated"})

}
func (c *controller) CreateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	postID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
	var post models.Post
	if err := c.db.First(&post, postID).Error; err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	var user models.User
	if err := c.db.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	newComment := models.Comment{
		PostID:     post.ID,
		AuthorName: user.Name,
		Content:    comment.Content,
		CreatedAt:  time.Now(),
	}
	if err := c.db.Create(&newComment).Error; err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Comment successfully created"})
}
func (c *controller) ListAllComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	postID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var comments []models.Comment
	if err := c.db.Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		http.Error(w, "Failed to fetch comments for the post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
