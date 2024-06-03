package middlewares

import (
	"context"
	"net/http"

	"github.com/LaQuannT/astronaut-api/internal/model"
	"github.com/LaQuannT/astronaut-api/internal/service"
	"github.com/LaQuannT/astronaut-api/internal/transport/handlers"
	"github.com/google/uuid"
)

type requestUser string

const requestUsr requestUser = "request-user"

func VerifyAPIKey(repository model.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("X-API-KEY")
			if key == "" {
				err := &model.APIError{
					Code:      http.StatusUnauthorized,
					Message:   "User unathorised to make request",
					Exception: "user did not supply api key in request header",
				}
				handlers.WriteError(w, err)
				return
			}

			_, err := uuid.Parse(key)
			if err != nil {
				err = &model.APIError{
					Code:      http.StatusUnauthorized,
					Message:   "User unathorised to make request",
					Exception: "user supplied invalid api key format",
				}
				handlers.WriteError(w, err)
				return
			}

			usr, err := service.SearchAPIKey(r.Context(), repository, key)
			if err != nil {
				handlers.WriteError(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), requestUsr, usr)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
