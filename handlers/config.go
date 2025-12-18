package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/anjaobradovic/ars-sit-2025/model"
	"github.com/anjaobradovic/ars-sit-2025/services"
	"github.com/gorilla/mux"
)

type ConfigHandler struct {
	service *services.ConfigService
}

func NewConfigHandler(service *services.ConfigService) *ConfigHandler {
	return &ConfigHandler{service: service}
}

func (h *ConfigHandler) CreateConfig(w http.ResponseWriter, r *http.Request) {
	var config model.Config

	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := h.service.Create(config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(config)
}

func (h *ConfigHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	config, err := h.service.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(config)
}

func (h *ConfigHandler) DeleteConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.service.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
