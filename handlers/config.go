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

// POST /configs
func (h *ConfigHandler) CreateConfig(w http.ResponseWriter, r *http.Request) {
	var config model.Config
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := h.service.Create(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(config)
}

// GET /configs/{id}/versions/{version}
func (h *ConfigHandler) GetConfigByVersion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	version := vars["version"]

	config, err := h.service.Get(id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(config)
}

// DELETE /configs/{id}/versions/{version}
func (h *ConfigHandler) DeleteConfigByVersion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	version := vars["version"]

	if err := h.service.Delete(id, version); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
