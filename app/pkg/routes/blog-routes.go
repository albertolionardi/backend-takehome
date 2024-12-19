package routes

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func BlogRoutes(router *mux.Router, db *gorm.DB) {
	controllers := NewController(db)
	router.HandleFunc("/register", controllers.Register).Methods("POST")

	router.HandleFunc("/login", controllers.Login).Methods("POST")

	router.HandleFunc("/posts", controllers.CreatePost).Methods("POST")
	router.HandleFunc("/posts/{id}", controllers.GetPostById).Methods("GET")
	router.HandleFunc("/posts", controllers.GetPosts).Methods("GET")
	router.HandleFunc("/posts/{id}", controllers.DeletePostById).Methods("DELETE")
	router.HandleFunc("/posts/{id}", controllers.UpdatePost).Methods("PUT")

	router.HandleFunc("/posts/{id}/comments", controllers.CreateComment).Methods("POST")
	router.HandleFunc("/posts/{id}/comments", controllers.ListAllComments).Methods("GET")

}
