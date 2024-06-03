package transport

import (
	"log/slog"
	"net/http"

	"github.com/LaQuannT/astronaut-api/internal/model"
	"github.com/LaQuannT/astronaut-api/internal/transport/middlewares"
)

func NewServer(logger *slog.Logger, usrRepository model.UserRepository) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, usrRepository)

	var handler http.Handler = mux
	mw := middlewares.RequestLogger(logger)
	handler = mw(handler)
	return handler
}
