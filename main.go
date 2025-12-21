// Package classification Configuration Service API.
//
// Documentation of our Configuration Service API.
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
// Title: Configuration Service API
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/hashicorp/consul/api"

	"github.com/anjaobradovic/ars-sit-2025/handlers"
	"github.com/anjaobradovic/ars-sit-2025/metrics"
	"github.com/anjaobradovic/ars-sit-2025/middleware"
	"github.com/anjaobradovic/ars-sit-2025/repositories"
	"github.com/anjaobradovic/ars-sit-2025/services"
	"github.com/anjaobradovic/ars-sit-2025/tracing"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

func main() {
	// ---- Tracing init ----
	rootCtx := context.Background()

	shutdownTracer, err := tracing.InitTracer(rootCtx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = shutdownTracer(rootCtx) }()

	// ---- Consul / repos / services / handlers ----
	// Izaberi JEDNU adresu u zavisnosti kako pokrećeš.
	consulAddr := "consul:8500" // docker-compose varijanta

	configRepo, err := repositories.NewConfigRepository(consulAddr)
	if err != nil {
		log.Fatal(err)
	}
	configService := services.NewConfigService(configRepo)
	configHandler := handlers.NewConfigHandler(configService)

	groupRepo, err := repositories.NewGroupRepository(consulAddr)
	if err != nil {
		log.Fatal(err)
	}
	groupService := services.NewGroupService(groupRepo)
	groupHandler := handlers.NewGroupHandler(groupService)

	cfg := api.DefaultConfig()
	cfg.Address = consulAddr
	consulClient, err := api.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// ---- Router
	r := mux.NewRouter()

	// Middleware
	r.Use(middleware.MetricsMiddleware)
	r.Use(otelmux.Middleware(
		"config-service",
		otelmux.WithFilter(func(r *http.Request) bool {
			return r.URL.Path != "/metrics"
		}),
	))

	// Routes
	r.Handle("/metrics", metrics.MetricsHandler())

	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	r.Handle("/docs", sh)

	r.Handle("/configs", middleware.IdempotencyMiddleware(consulClient)(http.HandlerFunc(configHandler.CreateConfig))).Methods("POST")
	r.HandleFunc("/configs/{name}/versions/{version}", configHandler.GetConfigByVersion).Methods("GET")
	r.HandleFunc("/configs/{name}/versions/{version}", configHandler.DeleteConfigByVersion).Methods("DELETE")

	r.HandleFunc("/groups", groupHandler.CreateGroup).Methods("POST")
	r.HandleFunc("/groups/{name}/versions/{version}", groupHandler.GetGroup).Methods("GET")
	r.HandleFunc("/groups/{name}/versions/{version}", groupHandler.DeleteGroup).Methods("DELETE")
	r.HandleFunc("/groups/{name}/versions/{version}/add-config", groupHandler.AddConfig).Methods("POST")
	r.HandleFunc("/groups/{name}/versions/{version}/remove-config", groupHandler.RemoveConfig).Methods("POST")
	r.HandleFunc("/groups/{name}/versions/{version}/configs", groupHandler.GetConfigsByLabels).Methods("GET")

	// ---- Server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// OS signal channel
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Start server
	go func() {
		log.Println("Config service running on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Wait for shutdown
	<-quit
	log.Println("Shutting down Config service...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server stopped gracefully")
}
