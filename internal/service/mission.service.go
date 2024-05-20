package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/LaQuannT/astronaut-api/internal/model"
	"github.com/lib/pq"
	"net/http"
	"time"
)

func AddMission(ctx context.Context, r model.MissionRepository, m *model.Mission) (*model.Mission, error) {
	if err := validate(m, "Mission"); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := r.CreateMission(ctx, m); err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, &model.APIError{
				Code:      http.StatusBadRequest,
				Message:   "Mission already exists",
				Exception: pgErr.Message,
			}
		}
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "fail to add mission",
			Exception: err.Error(),
		}
	}

	return m, nil
}

func GetMission(ctx context.Context, r model.MissionRepository, id int) (*model.Mission, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m, err := r.FindMissionByID(ctx, id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, &model.APIError{
			Code:      http.StatusNotFound,
			Message:   "Mission not found",
			Exception: err.Error(),
		}
	case err != nil:
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "fail to get mission",
			Exception: err.Error(),
		}
	default:
		return m, nil
	}
}

func GetMissions(ctx context.Context, r model.MissionRepository) ([]*model.Mission, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	missions, err := r.FindAllMissions(ctx)
	if err != nil {
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "fail to get missions",
			Exception: err.Error(),
		}
	}
	return missions, nil
}

func SearchMissionName(ctx context.Context, r model.MissionRepository, target string) ([]*model.Mission, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	missions, err := r.FindMissionByNameOrAlias(ctx, target)
	if err != nil {
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   fmt.Sprintf("failed to get any mission with or containing name: %s", target),
			Exception: err.Error(),
		}
	}
	return missions, nil
}

func UpdateMission(ctx context.Context, r model.MissionRepository, m *model.Mission) error {
	if err := validate(m, "Mission"); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := r.UpdateMission(ctx, m); err != nil {
		return &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "fail to update mission",
			Exception: err.Error(),
		}
	}
	return nil
}

func RegisterAstronautToMission(ctx context.Context, r model.MissionRepository, astronautID, missionID int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := r.CreateAstronautMission(ctx, astronautID, missionID); err != nil {
		return &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to register astronaut to mission",
			Exception: err.Error(),
		}
	}
	return nil
}

func GetMissionsByAstronaut(ctx context.Context, r model.MissionRepository, astronautID int) ([]*model.Mission, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	missions, err := r.FindMissionsByAstronaut(ctx, astronautID)
	if err != nil {
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to get missions by astronaut",
			Exception: err.Error(),
		}
	}
	return missions, nil
}

func RemoveAstronautFromMission(ctx context.Context, r model.MissionRepository, astronautID, missionID int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := r.DeleteAstronautMission(ctx, astronautID, missionID)
	switch {
	case errors.Is(err, model.ErrNoChange):
		return &model.APIError{
			Code:      http.StatusNotFound,
			Message:   "Mission and/or Astronaut not found",
			Exception: err.Error(),
		}
	case err != nil:
		return &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to remove astronaut from mission",
			Exception: err.Error(),
		}
	default:
		return nil
	}
}

func DeleteMission(ctx context.Context, r model.MissionRepository, id int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := r.DeleteMission(ctx, id)
	switch {
	case errors.Is(err, model.ErrNoChange):
		return &model.APIError{
			Code:      http.StatusNotFound,
			Message:   "Mission not found",
			Exception: err.Error(),
		}
	case err != nil:
		return &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "fail to delete mission",
			Exception: err.Error(),
		}
	default:
		return nil

	}
}
