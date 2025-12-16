package routers

import (
	"net/http"

	"github.com/anjaobradovic/ars-sit-2025/handlers"
)

func RegisterRoutes() {
	http.HandleFunc("/configs", handlers.CreateConfigHandler)
}
