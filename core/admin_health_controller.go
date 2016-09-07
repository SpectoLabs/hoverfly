package hoverfly

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
)

type HealthHandler struct{}

func (this *HealthHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	w.Header().Set("Content-Type", "application/json")

	var response messageResponse
	response.Message = "Hoverfly is healthy"

	response.Encode()

	b, err := response.Encode()
	if err != nil {
		// failed to read response body
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Could not encode response body!")
		http.Error(w, "Failed to encode response", 500)
		return
	}
	w.Write(b)
	return
}
