package main

import (
	"errors"
	"net/http"
	"temprest/config"
	"temprest/geolocationapi"
	"temprest/healthcheck"
	"temprest/logging"

	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Retrieve the server port from configuration
var serverPort = config.GetString("HttpServer.Port")

// @title Swagger Example API
// @version 1.0
// @description This is a sample API for demonstrating Swagger in Golang with the chi package.
// @host localhost:8080
// @BasePath /
func main() {
	// Create a new router using Chi
	startServer()
}

func startServer() {

	r := chi.NewRouter()
	// Serve Swagger JSON

	r.Mount("/", healthcheck.GetRoutes())
	r.Mount("/geolocationapi", geolocationapi.GetRoutes())

	// Serve Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"), // The path to your swagger.json file
	))

	// Serve Swagger JSON
	r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})
	// Define your API routes here
	r.Get("/hello", hello)
	logging.DoLoggingLevelBasedLogs(logging.Info, "starting http server at port: "+serverPort, nil)
	if err := http.ListenAndServe(":"+serverPort, r); err != nil {
		logging.DoLoggingLevelBasedLogs(logging.Error, "", logging.EnrichErrorWithStackTrace(errors.New("http connection error: "+err.Error())))
	}

}

// hello godoc
// @Summary Get a hello message
// @Description Get a simple hello message
// @Tags hello
// @Accept  json
// @Produce  json
// @Success 200 {string} string	"Hello, world!"
// @Router /hello [get]
func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}
