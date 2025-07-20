package rest

import (
	"log"
	"net/http"

	"github.com/pabloeclair/rest-subscription/internal/sbscrb/models"
)

func LoggingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := models.NewLoggingResponseWriter(w)
		handler.ServeHTTP(lrw, r)

		var statusMsg string
		if lrw.StatusCode == http.StatusInternalServerError {
			statusMsg = models.RedString(lrw.StatusMessage)
		} else {
			statusMsg = lrw.StatusMessage
		}
		log.Printf(
			"%s %s: %d - %s",
			r.Method,
			r.URL.Path,
			lrw.StatusCode,
			statusMsg,
		)
	})
}
