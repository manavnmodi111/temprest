package geolocationapi

import (
	"net/http"

	"github.com/go-chi/chi"
)

func GetRoutes() http.Handler {

	r := chi.NewRouter()

	r.Get("/items/{name}", GetItemByID)
	r.Post("/items", CreateItem)
	r.Get("/location/{id}", GetLocationByID)
	r.Get("/location", GetLocation)
	r.Post("/location", CreateLocation)

	return r
}
