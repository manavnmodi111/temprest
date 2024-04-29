package geolocationapi

import (
	"net/http"

	"github.com/go-chi/chi"
)

func GetRoutes() http.Handler {

	r := chi.NewRouter()

	//Endpoints for membership
	r.Get("/membership", GetMembership)
	r.Get("/membership/{id}", GetMembershipByID)
	r.Post("/membership", CreateMembership)
	r.Put("/membership/{id}", UpdateMembershipByID)
	r.Delete("/membership/{id}", DeleteMembershipByID)

	//Endpoints for location
	r.Get("/location/{id}", GetLocationByID)
	r.Get("/location", GetLocation)
	r.Post("/location", CreateLocation)
	r.Put("/location/{id}", UpdateLocationByID)
	r.Delete("/location/{id}", DeleteLocationByID)

	//Endpoints for membership
	r.Get("/community", GetCommunity)
	r.Get("/community/{id}", GetCommunityByID)
	r.Post("/community", CreateCommunity)
	r.Put("/community/{id}", UpdateCommunityByID)
	r.Delete("/community/{id}", DeleteCommunityByID)

	return r
}
