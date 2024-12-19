package routes

import (
	"app/middleware"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func BlogRoutes(router *mux.Router, db *gorm.DB) {
	controllers := NewController(db)
	sessionMiddleware := middleware.SessionMiddleware(db)
	// Public routes
	router.HandleFunc("/register", controllers.Register).Methods("POST")
	router.HandleFunc("/login", controllers.Login).Methods("POST")

	// Routes for middleware
	proctedtedRouter := router.PathPrefix("/").Subrouter()
	proctedtedRouter.Use(sessionMiddleware)

	proctedtedRouter.HandleFunc("/posts", controllers.CreatePost).Methods("POST")
	proctedtedRouter.HandleFunc("/posts/{id}", controllers.GetPostById).Methods("GET")
	proctedtedRouter.HandleFunc("/posts", controllers.GetPosts).Methods("GET")
	proctedtedRouter.HandleFunc("/posts/{id}", controllers.DeletePostById).Methods("DELETE")
	proctedtedRouter.HandleFunc("/posts/{id}", controllers.UpdatePost).Methods("PUT")

	proctedtedRouter.HandleFunc("/posts/{id}/comments", controllers.CreateComment).Methods("POST")
	proctedtedRouter.HandleFunc("/posts/{id}/comments", controllers.ListAllComments).Methods("GET")

}
