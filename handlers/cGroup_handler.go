package handlers

import (
	"encoding/json"
	"log"
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

// CreateGroup creates a new configuration group
// swagger:route POST /groups groups createGroup
//
// # Create a new configuration group
//
// This endpoint creates a configuration group with the provided data.
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Parameters:
// - name: group
//   in: body
//   description: ConfigurationGroup object to create
//   required: true
//   schema:
//     $ref: '#/definitions/ConfigurationGroup'
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
// # Get a configuration group
//
// This endpoint retrieves a specific configuration group by name and version.
//
// Parameters:
// - name: name
//   in: path
//   required: true
//   type: string
//   description: Name of the configuration group
// - name: version
//   in: path
//   required: true
//   type: string
//   description: Version of the configuration group
//
// Responses:
//   200: body:ConfigurationGroup
//   404: body:ErrorResponse

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
// # Delete a configuration group
//
// This endpoint deletes a configuration group by name and version.
//
// Parameters:
// - name: name
//   in: path
//   required: true
//   type: string
// - name: version
//   in: path
//   required: true
//   type: string
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
// # Add configuration to group
//
// This endpoint adds a labeled configuration to an existing group.
//
// Consumes:
// - application/json
//
// Parameters:
//   - name: name
//     in: path
//     required: true
//     type: string
//   - name: version
//     in: path
//     required: true
//     type: string
//   - name: config
//     in: body
//     required: true
//     description: LabeledConfiguration to add
//     schema:
//     $ref: '#/definitions/LabeledConfiguration'
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
// # Remove configuration from group
//
// This endpoint removes a labeled configuration from a group.
//
// Consumes:
// - application/json
//
// Parameters:
// - name: name
//   in: path
//   required: true
//   type: string
// - name: version
//   in: path
//   required: true
//   type: string
// - name: configId
//   in: body
//   required: true
//   description: ID of the configuration to remove
//   schema:
//     type: object
//     properties:
//       configId:
//         type: string
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
// # Get configurations by labels
//
// This endpoint retrieves all configurations in a group that match the provided labels.
//
// Parameters:
//   - name: name
//     in: path
//     required: true
//     type: string
//   - name: version
//     in: path
//     required: true
//     type: string
//   - name: query
//     in: query
//     required: false
//     type: string
//     description: Labels to filter configurations (key=value)
//
// Responses:
//
//	200: body:[LabeledConfiguration]
//	404: body:ErrorResponse
func (h *GroupHandler) GetConfigsByLabels(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	group, err := h.service.Get(vars["name"], vars["version"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	queryLabels := r.URL.Query()
	result := []*model.LabeledConfiguration{}

	for _, cfg := range group.Configurations {
		match := true
		for key, values := range queryLabels {
			if cfg.Labels == nil || cfg.Labels[key] != values[0] {
				match = false
				break
			}
		}
		if match {
			result = append(result, cfg)
		}
	}

	json.NewEncoder(w).Encode(result)
}
