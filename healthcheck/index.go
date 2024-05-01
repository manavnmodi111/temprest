package healthcheck

import (
	"net/http"

	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
)

func GetRoutes() http.Handler {
	r := chi.NewRouter()
	// Serve Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"), // The path to your swagger.json file
	))

	// Serve Swagger JSON
	r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})
	r.Get("/healthcheck", healthCheck)
	return r
}
