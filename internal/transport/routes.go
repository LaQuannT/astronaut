package transport

import (
	"net/http"

	"github.com/LaQuannT/astronaut-api/internal/model"
	"github.com/LaQuannT/astronaut-api/internal/transport/handlers"
)

func addRoutes(
	mux *http.ServeMux,
	userRepository model.UserRepository,
	astronautRepository model.AstronautRepository,
) {
	// user routes
	mux.Handle("POST /api/v1/register", handlers.HandleRegisterUser(userRepository))
	mux.Handle("GET /api/v1/user", handlers.HandleGetUser(userRepository))
	mux.Handle("GET /api/v1/users", handlers.HandleGetUsers(userRepository))
	mux.Handle("PUT /api/v1/users/{userID}", handlers.HandleUpdateUser(userRepository))
	mux.Handle("DELETE /api/v1/users/{userID}", handlers.HandleDeleteUser(userRepository))
	mux.Handle("PUT /api/v1/users/password/{userID}", handlers.HandlePasswordReset(userRepository))
	mux.Handle("PUT /api/v1/users/apikey/{userID}", handlers.HandleAPIKeyReset(userRepository))

	// astronaut routes
	mux.Handle("POST /api/v1/astonauts", handlers.HandleCreateAstronaut(astronautRepository))
	mux.Handle("GET /api/v1/astonauts", handlers.HandleGetAstronauts(astronautRepository))
	mux.Handle("GET /api/v1/astronauts/search", handlers.HandleSearchAstronautName(astronautRepository))
	mux.Handle("GET /api/v1/astonauts/{astronautID}", handlers.HandleGetAstronaut(astronautRepository))
	mux.Handle("PUT /api/v1/astronauts/{astronautID}", handlers.HandleUpdateAstronaut(astronautRepository))
	mux.Handle("DELETE /api/v1/astronauts/{astronautID}", handlers.HandleDeleteAstronaut(astronautRepository))
}
