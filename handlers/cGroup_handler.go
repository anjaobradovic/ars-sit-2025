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
	if err := h.service.Create(group); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(group)
}

func (h *GroupHandler) GetGroup(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	group, err := h.service.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(group)
}

func (h *GroupHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.service.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *GroupHandler) AddConfig(w http.ResponseWriter, r *http.Request) {
	groupID := mux.Vars(r)["id"]
	var labeledConfig model.LabeledConfiguration
	if err := json.NewDecoder(r.Body).Decode(&labeledConfig); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if err := h.service.AddConfigToGroup(groupID, labeledConfig); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *GroupHandler) RemoveConfig(w http.ResponseWriter, r *http.Request) {
	groupID := mux.Vars(r)["id"]
	var payload struct {
		ConfigID string `json:"configId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if err := h.service.RemoveConfigFromGroup(groupID, payload.ConfigID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
