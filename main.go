package main

import (
	"log"
	"net/http"

	"github.com/anjaobradovic/ars-sit-2025/routers"
)

func main() {
	routers.RegisterRoutes()

	log.Println("Config service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
