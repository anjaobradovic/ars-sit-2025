package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/anjaobradovic/ars-sit-2025/model"
	"github.com/anjaobradovic/ars-sit-2025/services"
	"github.com/google/uuid"
)

func CreateConfigHandler(w http.ResponseWriter, r *http.Request) { // ✅ prvo slovo veliko
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var config model.Config
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config.ID = uuid.NewString()

	if err := services.CreateConfig(config); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(config)
}

func AddConfig(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var cfg model.Config
	err := json.NewDecoder(r.Body).Decode(&cfg)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = services.AddConfig(cfg)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cfg)
}

func GetConfig(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	version := r.URL.Query().Get("version")

	cfg, err := services.GetConfig(name, version)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cfg)
}

func Config(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		AddConfig(w, r)
	case http.MethodGet:
		GetConfig(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
