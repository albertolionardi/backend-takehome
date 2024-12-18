package main

import (
	"app/pkg/config"
	"app/pkg/models"
	"app/pkg/routes"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	db := config.Connect()
	models.MigrateComment(db)
	models.MigratePost(db)
	models.MigrateUser(db)
	models.MigrateSession(db)

	router := mux.NewRouter().StrictSlash(true)
	routes.BlogRoutes(router, db)

	fmt.Println("Server is running on222 http://localhost:8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
