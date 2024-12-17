package main

import (
	"app/pkg/config"
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	http.HandleFunc("/", handler)
	config.Connect()
	fmt.Println("Server is running on222 http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}
