package geolocationapi

import (
	"net/http"

	"github.com/go-chi/chi"
)

func GetRoutes() http.Handler {

	r := chi.NewRouter()

	r.Get("/membership", GetMembership)
	r.Get("/membership/{id}", GetMembershipByID)
	r.Post("/membership", CreateMembership)
	r.Put("/membership/{id}", UpdateMembershipByID)
	r.Get("/location/{id}", GetLocationByID)
	r.Get("/location", GetLocation)
	r.Post("/location", CreateLocation)
	r.Put("/location/{id}", UpdateLocationByID)
	r.Get("/community", GetCommunity)
	r.Get("/community/{id}", GetCommunityByID)
	r.Post("/community", CreateCommunity)
	r.Put("/community/{id}", UpdateCommunityByID)
	return r
}
