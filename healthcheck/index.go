package healthcheck

import (
	"net/http"

	"github.com/go-chi/chi"
)

func GetRoutes() http.Handler {
	r := chi.NewRouter()
	r.Get("/healthcheck", healthCheck)
	return r
}
