package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LaQuannT/astronaut-api/internal/model"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error) {
	var apiErr *model.APIError

	e := struct{ Error string }{Error: ""}

	if errors.As(err, &apiErr) {
		e.Error = apiErr.Message
		writeJSON(w, apiErr.Code, e)
		return
	}
	e.Error = "unable to process request"
	writeJSON(w, http.StatusInternalServerError, e)
}
