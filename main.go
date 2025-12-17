package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/anjaobradovic/ars-sit-2025/handlers"
	"github.com/anjaobradovic/ars-sit-2025/repositories"
	"github.com/anjaobradovic/ars-sit-2025/services"
)

func main() {
	repo, err := repositories.NewConfigRepository("consul:8500")
	if err != nil {
		log.Fatal(err)
	}

	service := services.NewConfigService(repo)
	handler := handlers.NewConfigHandler(service)

	//routers
	r := mux.NewRouter()
	r.HandleFunc("/configs", handler.CreateConfig).Methods("POST")
	r.HandleFunc("/configs/{id}", handler.GetConfig).Methods("GET")
	r.HandleFunc("/configs/{id}", handler.DeleteConfig).Methods("DELETE")

	//server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Channel - OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// server start
	go func() {
		log.Println("Config service running on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-quit
	log.Println("Shutting down Config service...")

	// Graceful shutdown with timer
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server stopped gracefully")
}
