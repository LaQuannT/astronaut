package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/LaQuannT/astronaut-api/internal/model"
	"github.com/LaQuannT/astronaut-api/internal/service"
)

func HandleCreateMission(repository model.MissionRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := new(model.Mission)

		if err := json.NewDecoder(r.Body).Decode(m); err != nil {
			WriteError(w, err)
			return
		}

		m, err := service.AddMission(r.Context(), repository, m)
		if err != nil {
			WriteError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, m)
	}
}

func HandleGetMission(repository model.MissionRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mid := r.PathValue("missionID")

		id, err := strconv.Atoi(mid)
		if err != nil {
			WriteError(w, err)
			return
		}

		m, err := service.GetMission(r.Context(), repository, id)
		if err != nil {
			WriteError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, m)
	}
}

func HandleGetMissions(repository model.MissionRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// add limit and offset

		ms, err := service.GetMissions(r.Context(), repository)
		if err != nil {
			WriteError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, ms)
	}
}
