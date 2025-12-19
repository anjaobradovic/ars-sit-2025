package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/anjaobradovic/ars-sit-2025/model"
	"github.com/anjaobradovic/ars-sit-2025/services"
	"github.com/gorilla/mux"
)

type GroupHandler struct {
	service *services.GroupService
}

func NewGroupHandler(service *services.GroupService) *GroupHandler {
	return &GroupHandler{service: service}
}

func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group model.ConfigurationGroup
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.service.Create(&group); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(group)
}

func (h *GroupHandler) GetGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	group, err := h.service.Get(vars["name"], vars["version"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(group)
}

func (h *GroupHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if err := h.service.Delete(vars["name"], vars["version"]); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *GroupHandler) AddConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var cfg model.LabeledConfiguration

	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.service.AddConfig(vars["name"], vars["version"], cfg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *GroupHandler) RemoveConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var payload struct {
		ConfigID string `json:"configId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.service.RemoveConfig(vars["name"], vars["version"], payload.ConfigID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
