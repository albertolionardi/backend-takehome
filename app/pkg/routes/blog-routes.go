package routes

import (
	"app/pkg/controllers"

	"github.com/gorilla/mux"
)

var BlogRoutes = func(router *mux.Router) {
	router.HandleFunc("/register/", controllers.Register).Methods("POST")
	router.HandleFunc("/login/", controllers.Login).Methods("POST")

	router.HandleFunc("/posts/", controllers.CreatePost).Methods("POST")
	router.HandleFunc("/posts/{id}", controllers.GetPostById).Methods("GET")
	router.HandleFunc("/posts/", controllers.GetPosts).Methods("GET")
	router.HandleFunc("/posts/{id}", controllers.DeletePostById).Methods("DELETE")

	router.HandleFunc("/posts/{id}/comments/", controllers.PostComment).Methods("POST")
	router.HandleFunc("/posts/{id}/comments/", controllers.ListAllComments).Methods("GET")

}
