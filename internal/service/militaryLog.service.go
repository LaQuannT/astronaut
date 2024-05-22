package service

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/LaQuannT/astronaut-api/internal/model"
)

func AddMilitaryLog(ctx context.Context, militaryLogRepo model.MilitaryLogRepository, ml *model.MilitaryLog) (*model.MilitaryLog, error) {
	if err := validate(ml, "MilitaryLog"); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := militaryLogRepo.CreateMilitaryLog(ctx, ml)
	if err != nil {
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to add Military Log",
			Exception: err.Error(),
		}
	}
	return ml, nil
}

func GetMilitaryLog(ctx context.Context, militaryLogRepo model.MilitaryLogRepository, astronautID int) (*model.MilitaryLog, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ml, err := militaryLogRepo.FindMilitaryLog(ctx, astronautID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, &model.APIError{
			Code:      http.StatusNotFound,
			Message:   "Military Log not found",
			Exception: err.Error(),
		}
	case err != nil:
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to get Astronaut Military Log",
			Exception: err.Error(),
		}
	default:
		return ml, nil
	}
}

func GetMilitaryLogs(ctx context.Context, militaryLogRepo model.MilitaryLogRepository) ([]*model.MilitaryLog, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	mls, err := militaryLogRepo.FindAllMilitaryLogs(ctx)
	if err != nil {
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to get astronaut Military Logs",
			Exception: err.Error(),
		}
	}
	return mls, err
}

func UpdateMilitaryLog(ctx context.Context, militaryLogRepo model.MilitaryLogRepository, ml *model.MilitaryLog) error {
	if err := validate(ml, "Military Log"); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := militaryLogRepo.UpdateMilitaryLog(ctx, ml)
	switch {
	case errors.Is(err, model.ErrNoChange):
		return &model.APIError{
			Code:      http.StatusNotFound,
			Message:   "Astronaut Military Log not found",
			Exception: err.Error(),
		}
	case err != nil:
		return &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to update Military Log",
			Exception: err.Error(),
		}
	default:
		return nil
	}
}

func DeleteMilitaryLog(ctx context.Context, militaryLogRepo model.MilitaryLogRepository, astronautID int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := militaryLogRepo.DeleteMilitaryLog(ctx, astronautID)
	switch {
	case errors.Is(err, model.ErrNoChange):
		return &model.APIError{
			Code:      http.StatusNotFound,
			Message:   "Astronaut Military Log not found",
			Exception: err.Error(),
		}
	case err != nil:
		return &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to delete astronaut Military Log",
			Exception: err.Error(),
		}
	default:
		return nil
	}
}
