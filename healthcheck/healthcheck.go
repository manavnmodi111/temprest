package healthcheck

import (
	"net/http"
	"temprest/logging"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	logging.DoLoggingLevelBasedLogs(logging.Info, "healthcheck success", nil)
	w.Write([]byte("success"))
}
