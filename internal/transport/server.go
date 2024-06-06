package transport

import (
	"log/slog"
	"net/http"

	"github.com/LaQuannT/astronaut-api/internal/model"
	"github.com/LaQuannT/astronaut-api/internal/transport/middlewares"
)

func NewServer(
	logger *slog.Logger,
	usrRepository model.UserRepository,
	astronautRepository model.AstronautRepository,
) http.Handler {
	mux := http.NewServeMux()

	addRoutes(
		mux,
		usrRepository,
		astronautRepository,
	)

	var handler http.Handler = mux
	handler = middlewares.EnableCors(handler)
	mw := middlewares.RequestLogger(logger)
	handler = mw(handler)
	return handler
}
