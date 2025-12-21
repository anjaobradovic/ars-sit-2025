package handlers

import (
	"encoding/json"
	"log"
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

// CreateConfig creates a new configuration
// swagger:route POST /configs configurations createConfiguration
//
// # Create a new configuration
//
// This endpoint creates a new configuration with the provided data.
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Parameters:
//   - name: config
//     in: body
//     description: Configuration object to create
//     required: true
//     schema:
//     $ref: '#/definitions/Config'
//
// Responses:
//
//	201: body:Config
//	400: body:ErrorResponse
//	409: body:ErrorResponse
//	500: body:ErrorResponse
func (h *ConfigHandler) CreateConfig(w http.ResponseWriter, r *http.Request) {
	var config model.Config
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		log.Printf("JSON decode error: %+v\n", err)
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := h.service.Create(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(config)
}

// GetConfigByVersion retrieves a configuration by name and version
// swagger:route GET /configs/{name}/versions/{version} configurations getConfigurationByNameAndVersion
//
// # Get configuration by name and version
//
// This endpoint retrieves a specific configuration by its name and version.
//
// Parameters:
//   - name: name
//     in: path
//     required: true
//     type: string
//     description: The name of the configuration
//   - name: version
//     in: path
//     required: true
//     type: string
//     description: The version of the configuration
//
// Produces:
// - application/json
//
// Responses:
//
//	200: body:Config
//	404: body:ErrorResponse
//	500: body:ErrorResponse
func (h *ConfigHandler) GetConfigByVersion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	version := vars["version"]

	config, err := h.service.Get(name, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(w).Encode(config)
}

// DeleteConfigByVersion removes a configuration by name and version
// swagger:route DELETE /configs/{name}/versions/{version} configurations deleteConfigurationByNameAndVersion
//
// # Delete configuration by name and version
//
// This endpoint deletes a specific configuration by its name and version.
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
//
// Responses:
//
//	204: body:NoContentResponse
//	404: body:ErrorResponse
func (h *ConfigHandler) DeleteConfigByVersion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	version := vars["version"]

	if err := h.service.Delete(name, version); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
