package middlewares

import "net/http"

const (
	allowedOrigin  = "*"
	allowedMethods = "GET, POST, PUT, DELETE, OPTIONS"
	allowedHeaders = "Origin, Content-Type, Accept"
)

func EnableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", allowedMethods)
		w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
		next.ServeHTTP(w, r)
	})
}
