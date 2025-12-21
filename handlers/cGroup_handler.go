package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

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

// CreateGroup creates a new configuration group
// swagger:route POST /groups groups createGroup
//
// Create a new configuration group.
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
//
// Responses:
//   201: body:ConfigurationGroup
//   400: body:ErrorResponse

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

// GetGroup retrieves a configuration group by name and version
// swagger:route GET /groups/{name}/versions/{version} groups getGroup
//
// Get a configuration group.
//
// This endpoint retrieves a specific configuration group by name and version.
//
// Produces:
// - application/json
//
// Responses:
//
//	200: body:ConfigurationGroup
//	404: body:ErrorResponse
func (h *GroupHandler) GetGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	group, err := h.service.Get(vars["name"], vars["version"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(group)
}

// DeleteGroup deletes a configuration group by name and version
// swagger:route DELETE /groups/{name}/versions/{version} groups deleteGroup
//
// Delete a configuration group.
//
// This endpoint deletes a configuration group by name and version.
//
//
// Responses:
//   204: body:NoContentResponse
//   404: body:ErrorResponse

func (h *GroupHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if err := h.service.Delete(vars["name"], vars["version"]); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddConfig adds a configuration to a group
// swagger:route POST /groups/{name}/versions/{version}/add-config groups addConfig
//
// Add configuration to group.
//
// This endpoint adds a labeled configuration to an existing group.
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Responses:
//
//	200: body:LabeledConfiguration
//	400: body:ErrorResponse
func (h *GroupHandler) AddConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var cfg model.LabeledConfiguration

	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	log.Println("Handler: AddConfig called")
	log.Printf("Group: %s %s\n", vars["name"], vars["version"])
	log.Printf("Config payload: %+v\n", cfg)

	if err := h.service.AddConfig(vars["name"], vars["version"], cfg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// RemoveConfig removes a configuration from a group
// swagger:route POST /groups/{name}/versions/{version}/remove-config groups removeConfig
//
// Remove configuration from group.
//
// This endpoint removes a labeled configuration from a group by labeled configuration ID.
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Responses:
//   200: body:NoContentResponse
//   400: body:ErrorResponse

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

// GetConfigsByLabels gets configurations from a group filtered by labels
// swagger:route GET /groups/{name}/versions/{version}/configs groups getConfigsByLabels
//
// Get configurations by labels.
//
// This endpoint retrieves all configurations in a group that match the provided labels.
// All labels from the query must match (AND).
//
// Produces:
// - application/json
//
// Responses:
//   200: labeledConfigurationsResponse
//   400: body:ErrorResponse
//   404: body:ErrorResponse

func (h *GroupHandler) GetConfigsByLabels(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	group, err := h.service.Get(vars["name"], vars["version"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Očekujemo: ?labels=env:prod;region:eu
	raw := strings.TrimSpace(r.URL.Query().Get("labels"))

	// Ako nema labels parametra, vrati sve konfiguracije u grupi
	if raw == "" {
		_ = json.NewEncoder(w).Encode(group.Configurations)
		return
	}

	// Parsiranje: key:value;key2:value2
	queryLabels := map[string]string{}
	pairs := strings.Split(raw, ";")
	for _, p := range pairs {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		kv := strings.SplitN(p, ":", 2)
		if len(kv) != 2 {
			http.Error(w, "invalid labels format, expected key:value;key2:value2", http.StatusBadRequest)
			return
		}

		k := strings.TrimSpace(kv[0])
		v := strings.TrimSpace(kv[1])
		if k == "" || v == "" {
			http.Error(w, "invalid labels format, empty key or value", http.StatusBadRequest)
			return
		}

		queryLabels[k] = v
	}

	result := []*model.LabeledConfiguration{}

	// AND matching: sve labele iz upita moraju postojati u cfg.Labels i biti jednake
	for _, cfg := range group.Configurations {
		match := true
		for k, v := range queryLabels {
			if cfg.Labels == nil || cfg.Labels[k] != v {
				match = false
				break
			}
		}
		if match {
			result = append(result, cfg)
		}
	}

	_ = json.NewEncoder(w).Encode(result)
}
