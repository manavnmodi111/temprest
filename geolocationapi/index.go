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
	r.Get("/location/{id}", GetLocationByID)
	r.Get("/location", GetLocation)
	r.Post("/location", CreateLocation)
	r.Get("/community", GetCommunity)
	r.Get("/community/{id}", GetCommunityByID)
	r.Post("/community", CreateCommunity)
	return r
}
