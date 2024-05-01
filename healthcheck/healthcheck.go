package healthcheck

import (
	"net/http"
	"temprest/logging"
)

// healthcheck godoc
// @Summary Get a healthcheck message
// @Description Get a simple healthcheck message
// @Tags healthcheck
// @Accept  json
// @Produce  json
// @Success 200 {string} string	"success"
// @Router /healthcheck [get]
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	logging.DoLoggingLevelBasedLogs(logging.Info, "healthcheck success", nil)
	w.Write([]byte("success"))
}
