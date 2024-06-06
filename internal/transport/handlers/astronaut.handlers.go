package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/LaQuannT/astronaut-api/internal/model"
	"github.com/LaQuannT/astronaut-api/internal/service"
)

func HandleCreateAstronaut(repository model.AstronautRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a := new(model.Astronaut)

		if err := json.NewDecoder(r.Body).Decode(a); err != nil {
			WriteError(w, err)
			return
		}

		a, err := service.AddAstronaut(r.Context(), a, repository)
		if err != nil {
			WriteError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, a)
	}
}

func HandleGetAstronaut(repository model.AstronautRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		aid := r.PathValue("astronautID")

		id, err := strconv.Atoi(aid)
		if err != nil {
			WriteError(w, err)
			return
		}

		a, err := service.GetAstronaut(r.Context(), repository, id)
		if err != nil {
			WriteError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, a)
	}
}

func HandleGetAstronauts(repository model.AstronautRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// add limit and offset

		as, err := service.GetAstronauts(r.Context(), repository)
		if err != nil {
			WriteError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, as)
	}
}

func HandleUpdateAstronaut(repository model.AstronautRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		aid := r.PathValue("astroanutID")

		id, err := strconv.Atoi(aid)
		if err != nil {
			WriteError(w, err)
			return
		}

		a, err := service.GetAstronaut(r.Context(), repository, id)
		if err != nil {
			WriteError(w, err)
			return
		}

		err = json.NewDecoder(r.Body).Decode(a)
		switch {
		case errors.Is(err, io.EOF):
			WriteError(w, &model.APIError{
				Code:      http.StatusBadRequest,
				Message:   "astronaut data not provided in request body",
				Exception: err.Error(),
			})
			return

		case err != nil:
			WriteError(w, err)
			return
		}

		err = service.UpdateAstronaut(r.Context(), a, repository)
		if err != nil {
			WriteError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{"Message": "Astronaut has been updated"})
	}
}

func HandleDeleteAstronaut(repository model.AstronautRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		aid := r.PathValue("astroanutID")

		id, err := strconv.Atoi(aid)
		if err != nil {
			WriteError(w, err)
			return
		}

		err = service.DeleteAstronaut(r.Context(), repository, id)
		if err != nil {
			WriteError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{"Message": "Astronaut has been deleted"})
	}
}

func HandleSearchAstronautName(repository model.AstronautRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			WriteError(w, err)
			return
		}

		name := params.Get("name")

		a, err := service.SearchAstronautByName(r.Context(), repository, name)
		if err != nil {
			WriteError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, a)
	}
}
