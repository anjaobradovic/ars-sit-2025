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
	// --- Config repository + service ---
	configRepo, err := repositories.NewConfigRepository("consul:8500")
	if err != nil {
		log.Fatal(err)
	}
	configService := services.NewConfigService(configRepo)
	configHandler := handlers.NewConfigHandler(configService)

	// --- Configuration Groups ---
	groupRepo, err := repositories.NewGroupRepository("consul:8500")
	if err != nil {
		log.Fatal(err)
	}
	groupService := services.NewGroupService(groupRepo)
	groupHandler := handlers.NewGroupHandler(groupService)

	// --- Router ---
	r := mux.NewRouter()

	// Config endpoints
	r.HandleFunc("/configs", configHandler.CreateConfig).Methods("POST")
	r.HandleFunc("/configs/{id}/versions/{version}", configHandler.GetConfigByVersion).Methods("GET")
	r.HandleFunc("/configs/{id}/versions/{version}", configHandler.DeleteConfigByVersion).Methods("DELETE")

	// Configuration Groups
	r.HandleFunc("/groups", groupHandler.CreateGroup).Methods("POST")
	r.HandleFunc("/groups/{name}/versions/{version}", groupHandler.GetGroup).Methods("GET")
	r.HandleFunc("/groups/{name}/versions/{version}", groupHandler.DeleteGroup).Methods("DELETE")
	r.HandleFunc("/groups/{name}/versions/{version}/add-config", groupHandler.AddConfig).Methods("POST")
	r.HandleFunc("/groups/{name}/versions/{version}/remove-config", groupHandler.RemoveConfig).Methods("POST")
	r.HandleFunc("/groups/{name}/versions/{version}/configs", groupHandler.GetConfigsByLabels).Methods("GET")

	// --- HTTP server ---
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// OS signal channel
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Start server
	go func() {
		log.Println("Config service running on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Wait for shutdown signal
	<-quit
	log.Println("Shutting down Config service...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server stopped gracefully")
}
