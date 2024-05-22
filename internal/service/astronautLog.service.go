package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/LaQuannT/astronaut-api/internal/model"
	"net/http"
	"time"
)

func AddAstronautLog(ctx context.Context, astroLogRepo model.AstronautLogRepository, al *model.AstronautLog) (*model.AstronautLog, error) {
	if err := validate(al, "AstronautLog"); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := astroLogRepo.CreateAstronautLog(ctx, al)
	switch {
	case err != nil:
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to add AstronautLog",
			Exception: err.Error(),
		}
	default:
		return al, err

	}
}

func GetAstronautLog(ctx context.Context, astroLogRepo model.AstronautLogRepository, id int) (*model.AstronautLog, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	al, err := astroLogRepo.FindAstronautLogById(ctx, id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, &model.APIError{
			Code:      http.StatusNotFound,
			Message:   "AstronautLog not found",
			Exception: err.Error(),
		}
	case err != nil:
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to find AstronautLog",
			Exception: err.Error(),
		}
	default:
		return al, nil
	}
}

func GetAstronautLogs(ctx context.Context, astroLogRepo model.AstronautLogRepository) ([]*model.AstronautLog, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	als, err := astroLogRepo.FindAstronautLogs(ctx)
	if err != nil {
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to find AstronautLogs",
			Exception: err.Error(),
		}
	}
	return als, nil
}

func UpdateAstronautLog(ctx context.Context, astroLogRepo model.AstronautLogRepository, al *model.AstronautLog) error {
	if err := validate(al, "AstronautLog"); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := astroLogRepo.UpdateAstronautLog(ctx, al)
	switch {
	case errors.Is(err, model.ErrNoChange):
		return &model.APIError{
			Code:      http.StatusNotFound,
			Message:   "AstronautLog not found",
			Exception: err.Error(),
		}
	case err != nil:
		return &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to update AstronautLog",
			Exception: err.Error(),
		}
	default:
		return nil
	}
}

func DeleteAstronautLog(ctx context.Context, astroLogRepo model.AstronautLogRepository, id int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := astroLogRepo.DeleteAstronautLog(ctx, id)
	switch {
	case errors.Is(err, model.ErrNoChange):
		return &model.APIError{
			Code:      http.StatusNotFound,
			Message:   "AstronautLog not found",
			Exception: err.Error(),
		}
	case err != nil:
		return &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to delete AstronautLog",
			Exception: err.Error(),
		}
	default:
		return nil
	}
}
