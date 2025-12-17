package main

import (
	"log"
	"net/http"

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

	r := mux.NewRouter()
	r.HandleFunc("/configs", handler.CreateConfig).Methods("POST")

	log.Println("Config service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
