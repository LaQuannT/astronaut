package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/LaQuannT/astronaut-api/internal/service"

	"github.com/LaQuannT/astronaut-api/internal/model"
)

func HandleRegisterUser(repository model.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usr := new(model.User)

		if err := json.NewDecoder(r.Body).Decode(usr); err != nil {
			writeError(w, err)
			return
		}

		usr, err := service.RegisterUser(r.Context(), repository, usr)
		if err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, usr)
	}
}

func HandleGetUser(repository model.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			writeError(w, err)
		}

		email := params.Get("email")
		userID := params.Get("uid")

		usr := new(model.User)

		switch {
		case email != "":
			usr, err = service.SearchUserEmail(r.Context(), repository, email)
			if err != nil {
				writeError(w, err)
				return
			}
			writeJSON(w, http.StatusOK, usr)
			return

		case userID != "":
			uid, err := strconv.Atoi(userID)
			if err != nil {
				writeError(w, err)
				return
			}
			usr, err = service.SearchUserID(r.Context(), repository, uid)
			if err != nil {
				writeError(w, err)
				return
			}
			writeJSON(w, http.StatusOK, usr)
			return
		default:
			writeError(w, &model.APIError{
				Code:      http.StatusBadRequest,
				Message:   "email or uid param not supplied",
				Exception: "email or user ID params not supplied",
			})
		}
	}
}

func HandleGetUsers(repository model.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// add limit and offset

		urs, err := service.GetUsers(r.Context(), repository)
		if err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, urs)
	}
}

func HandleUpdateUser(repository model.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.PathValue("userID")

		id, err := strconv.Atoi(uid)
		if err != nil {
			writeError(w, err)
			return
		}

		usr, err := service.SearchUserID(r.Context(), repository, id)
		if err != nil {
			writeError(w, err)
			return
		}

		err = json.NewDecoder(r.Body).Decode(usr)
		switch {
		case errors.Is(err, io.EOF):
			writeError(w, &model.APIError{
				Code:      http.StatusBadRequest,
				Message:   "user data not provided in request body",
				Exception: err.Error(),
			})
			return

		case err != nil:
			writeError(w, err)
			return
		}

		err = service.UpdateUser(r.Context(), repository, usr)
		if err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{"Message": "User has been updated"})
	}
}

func HandleDeleteUser(repository model.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.PathValue("userID")

		id, err := strconv.Atoi(uid)
		if err != nil {
			writeError(w, err)
			return
		}

		if err := service.DeleteUser(r.Context(), repository, id); err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{"Message": "User has been deleted"})
	}
}

func HandlePasswordReset(repository model.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.PathValue("userID")

		id, err := strconv.Atoi(uid)
		if err != nil {
			writeError(w, err)
			return
		}

		usr := new(model.User)

		if err := json.NewDecoder(r.Body).Decode(usr); err != nil {
			writeError(w, err)
			return
		}

		err = service.ResetPassword(r.Context(), repository, usr.Password, id)
		if err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{"Message": "User password has been reset"})
	}
}

func HandleAPIKeyReset(repository model.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.PathValue("userID")

		id, err := strconv.Atoi(uid)
		if err != nil {
			writeError(w, err)
			return
		}

		key, err := service.GenerateNewAPIKey(r.Context(), repository, id)
		if err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{"apiKey": key})
	}
}
