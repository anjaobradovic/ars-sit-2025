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
	rootCtx := context.Background()

	shutdownTracer, err := tracing.InitTracer(rootCtx)
	if err != nil {
		log.Fatal(err)
	}

	consulAddr := "consul:8500"

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

	r := mux.NewRouter()

	// Metrics
	r.Use(middleware.MetricsMiddleware)

	// Rate limiter (skip metrics/docs/spec)
	rl := middleware.NewRateLimiter(10, 20, 2*time.Minute)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			switch req.URL.Path {
			case "/metrics", "/docs", "/swagger.yaml":
				next.ServeHTTP(w, req)
				return
			default:
				rl.Middleware(next).ServeHTTP(w, req)
			}
		})
	})

	// Tracing (skip metrics)
	r.Use(otelmux.Middleware(
		"config-service",
		otelmux.WithFilter(func(req *http.Request) bool {
			return req.URL.Path != "/metrics"
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

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Config service running on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-quit
	log.Println("Shutting down Config service...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	_ = shutdownTracer(ctx)

	log.Println("Server stopped gracefully")
}
