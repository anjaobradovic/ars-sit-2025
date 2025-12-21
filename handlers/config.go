package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"github.com/anjaobradovic/ars-sit-2025/model"
	"github.com/anjaobradovic/ars-sit-2025/services"
	"github.com/gorilla/mux"
)

var tracer = otel.Tracer("handlers/config")

type ConfigHandler struct {
	service *services.ConfigService
}

func NewConfigHandler(service *services.ConfigService) *ConfigHandler {
	return &ConfigHandler{service: service}
}

func (h *ConfigHandler) CreateConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "ConfigHandler.CreateConfig")
	defer span.End()

	var config model.Config
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid JSON body")
		log.Printf("JSON decode error: %+v\n", err)
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	span.SetAttributes(
		attribute.String("config.name", config.Name),
		attribute.String("config.version", config.Version),
	)

	if err := h.service.Create(ctx, &config); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "create failed")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(config)
}

func (h *ConfigHandler) GetConfigByVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	name := vars["name"]
	version := vars["version"]

	ctx, span := tracer.Start(ctx, "ConfigHandler.GetConfigByVersion")
	defer span.End()

	span.SetAttributes(
		attribute.String("config.name", name),
		attribute.String("config.version", version),
	)

	config, err := h.service.Get(ctx, name, version)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(w).Encode(config)
}

func (h *ConfigHandler) DeleteConfigByVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	name := vars["name"]
	version := vars["version"]

	ctx, span := tracer.Start(ctx, "ConfigHandler.DeleteConfigByVersion")
	defer span.End()

	span.SetAttributes(
		attribute.String("config.name", name),
		attribute.String("config.version", version),
	)

	if err := h.service.Delete(ctx, name, version); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "delete failed")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
