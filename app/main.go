package main

import (
	"app/pkg/config"
	"app/pkg/routes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	db := config.Connect()

	router := mux.NewRouter().StrictSlash(true)
	routes.BlogRoutes(router, db)

	srv := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		fmt.Println("Server is running on http://localhost:8080")
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown due to err : %s", err)
	}

	log.Println("Interrupted")

}
