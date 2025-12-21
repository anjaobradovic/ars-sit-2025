package middleware

import (
	"net/http"
)

type SwaggerUIOpts struct {
	SpecURL string
}

func SwaggerUI(opts SwaggerUIOpts, _ interface{}) http.Handler {
	return http.FileServer(http.Dir("./swagger"))
}
