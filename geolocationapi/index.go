package geolocationapi

import (
	"net/http"

	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
)

// GetRoutes returns the HTTP handler for the geolocation API.
func GetRoutes() http.Handler {

	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"), // The path to your swagger.json file
	))

	// Serve Swagger JSON
	r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})
	
	// Endpoints for membership
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
