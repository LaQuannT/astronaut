package middlewares

import (
	"context"
	"net/http"

	"github.com/LaQuannT/astronaut-api/internal/model"
	"github.com/LaQuannT/astronaut-api/internal/service"
	"github.com/LaQuannT/astronaut-api/internal/transport/handlers"
)

func AdminOnly(repository model.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			usr, err := getRequestUser(r.Context())
			if err != nil {
				handlers.WriteError(w, err)
				return
			}

			isAdmin, err := service.CheckAdminPermission(r.Context(), repository, usr.ID)
			if err != nil {
				handlers.WriteError(w, err)
				return
			}

			if !isAdmin {
				handlers.WriteError(w, &model.APIError{
					Code:    http.StatusForbidden,
					Message: "User unauthorised",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getRequestUser(ctx context.Context) (*model.User, error) {
	reqUsr, ok := ctx.Value(requestUsr).(*model.User)
	if !ok {
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to process request",
			Exception: "failed to get request user data from request context",
		}
	}

	return reqUsr, nil
}
