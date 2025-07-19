package rest

import (
	"log"
	"net/http"

	"github.com/pabloeclair/rest-subscription/internal/sbscrb"
)

func LoggingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := sbscrb.NewLoggingResponseWriter(w)
		handler.ServeHTTP(lrw, r)
		log.Printf(
			"%s %s: %d - %s",
			r.Method,
			r.URL.Path,
			lrw.StatusCode,
			lrw.StatusMessage,
		)
	})
}
