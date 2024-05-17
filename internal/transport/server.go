package transport

import "net/http"

func NewServer() http.Handler {
	mux := http.NewServeMux()
	// add routes

	var handler http.Handler = mux
	// add top level middleware CORS, auth, and logging
	return handler
}
