package service

import (
	"context"
	"errors"
	"github.com/LaQuannT/astronaut-api/internal/model"
	"github.com/lib/pq"
	"net/http"
	"time"
)

func AddMajor(ctx context.Context, repository model.AcademicLogRepository, major *model.Major) (*model.Major, error) {
	if err := validate(major, "Major"); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := repository.CreateMajor(ctx, major); err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, &model.APIError{
				Code:      http.StatusBadRequest,
				Message:   "Major already exists",
				Exception: pgErr.Message,
			}
		}
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to add Major",
			Exception: err.Error(),
		}
	}
	return major, nil
}

func AddAlmaMater(ctx context.Context, repository model.AcademicLogRepository, almaMater *model.AlmaMater) (*model.AlmaMater, error) {
	if err := validate(almaMater, "Alma Mater"); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := repository.CreateAlmaMater(ctx, almaMater); err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, &model.APIError{
				Code:      http.StatusBadRequest,
				Message:   "Alma Mater already exists",
				Exception: pgErr.Message,
			}
		}
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to add Alma Mater",
			Exception: err.Error(),
		}
	}
	return almaMater, nil
}

func AddAstronautUndergradMajor(ctx context.Context, repository model.AcademicLogRepository, astronautID, majorID int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := repository.AddUnderGradMajor(ctx, astronautID, majorID); err != nil {
		return model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to add Astronaut Undergrad Major",
			Exception: err.Error(),
		}
	}
	return nil
}

func AddAstronautGradMajor(ctx context.Context, repository model.AcademicLogRepository, astronautID, majorID int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := repository.AddGradMajor(ctx, astronautID, majorID); err != nil {
		return model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to add Astronaut Grad Major",
			Exception: err.Error(),
		}
	}
	return nil
}

func AddAstronautAlmaMater(ctx context.Context, repository model.AcademicLogRepository, astronautID, almaMaterID int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := repository.AddAstronautAlmaMater(ctx, astronautID, almaMaterID); err != nil {
		return model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to add Astronaut Alma Mater",
			Exception: err.Error(),
		}
	}
	return nil
}
