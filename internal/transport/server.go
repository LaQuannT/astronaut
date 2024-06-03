package transport

import (
	"net/http"

	"github.com/LaQuannT/astronaut-api/internal/model"
)

func NewServer(usrRepository model.UserRepository) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, usrRepository)

	var handler http.Handler = mux
	// add top level middleware CORS, auth, and logging
	return handler
}
