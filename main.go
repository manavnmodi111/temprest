package main

import (
	"errors"
	"net/http"
	"temprest/config"
	"temprest/geolocationapi"
	"temprest/healthcheck"
	"temprest/logging"

	"github.com/go-chi/chi"
)

// Retrieve the server port from configuration
var serverPort = config.GetString("HttpServer.Port")

func main() {
	// Create a new router using Chi
	startServer()
}

func startServer() {

	r := chi.NewRouter()

	r.Mount("/", healthcheck.GetRoutes())
	r.Mount("/geolocationapi", geolocationapi.GetRoutes())

	logging.DoLoggingLevelBasedLogs(logging.Info, "starting http server at port: "+serverPort, nil)
	if err := http.ListenAndServe(":"+serverPort, r); err != nil {
		logging.DoLoggingLevelBasedLogs(logging.Error, "", logging.EnrichErrorWithStackTrace(errors.New("http connection error: "+err.Error())))
	}
}
