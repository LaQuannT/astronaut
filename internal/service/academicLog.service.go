package service

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/LaQuannT/astronaut-api/internal/model"
	"github.com/lib/pq"
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

func UpdateMajor(ctx context.Context, repository model.AcademicLogRepository, major *model.Major) error {
	if err := validate(major, "Major"); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := repository.UpdateMajor(ctx, major); err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return &model.APIError{
				Code:      http.StatusBadRequest,
				Message:   "Major already exist",
				Exception: pgErr.Message,
			}
		}

		if errors.Is(err, model.ErrNoChange) {
			return &model.APIError{
				Code:      http.StatusBadRequest,
				Message:   "Major not found",
				Exception: err.Error(),
			}
		}
		return &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to update Major",
			Exception: err.Error(),
		}
	}
	return nil
}

func UpdateAlaMater(ctx context.Context, repository model.AcademicLogRepository, almaMater *model.AlmaMater) error {
	if err := validate(almaMater, "Alma Mater"); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := repository.UpdateAlmaMater(ctx, almaMater); err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			return &model.APIError{
				Code:      http.StatusBadRequest,
				Message:   "Alma Mater already exists",
				Exception: pgErr.Message,
			}
		}

		if errors.Is(err, model.ErrNoChange) {
			return &model.APIError{
				Code:      http.StatusBadRequest,
				Message:   "Alma Mater not found",
				Exception: err.Error(),
			}
		}
		return &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to update Alma Mater",
			Exception: err.Error(),
		}
	}
	return nil
}

func GetMajorByID(ctx context.Context, repository model.AcademicLogRepository, id int) (*model.Major, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m, err := repository.FindMajorByID(ctx, id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, &model.APIError{
			Code:      http.StatusNotFound,
			Message:   "Major Not Found",
			Exception: err.Error(),
		}
	case err != nil:
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to find major by id",
			Exception: err.Error(),
		}
	default:
		return m, nil
	}
}

func GetAlmaMaterByID(ctx context.Context, repository model.AcademicLogRepository, id int) (*model.AlmaMater, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	a, err := repository.FindAlmaMaterByID(ctx, id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, &model.APIError{
			Code:      http.StatusNotFound,
			Message:   "Alma Mater Not Found",
			Exception: err.Error(),
		}
	case err != nil:
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to find Alma Mater by id",
			Exception: err.Error(),
		}
	default:
		return a, nil
	}
}

func GetAstronautUndergradMajors(ctx context.Context, repository model.AcademicLogRepository, astronautID int) ([]*model.Major, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ms, err := repository.FindAstronautUnderGradMajors(ctx, astronautID)
	if err != nil {
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to find Astronaut Undergrad Majors",
			Exception: err.Error(),
		}
	}
	return ms, nil
}

func GetAstronautGradMajors(ctx context.Context, repository model.AcademicLogRepository, astronautID int) ([]*model.Major, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ms, err := repository.FindAstronautGradMajors(ctx, astronautID)
	if err != nil {
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to find Astronaut Grad Majors",
			Exception: err.Error(),
		}
	}
	return ms, nil
}

func GetAstronautAlmaMaters(ctx context.Context, repository model.AcademicLogRepository, astronautID int) ([]*model.AlmaMater, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	as, err := repository.FindAstronautAlmaMaters(ctx, astronautID)
	if err != nil {
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to find Astronaut Alma Maters",
			Exception: err.Error(),
		}
	}
	return as, nil
}

func DeleteMajor(ctx context.Context, repository model.AcademicLogRepository, id int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := repository.DeleteMajor(ctx, id)
	switch {
	case errors.Is(err, model.ErrNoChange):
		return &model.APIError{
			Code:      http.StatusNotFound,
			Message:   "Major Not Found",
			Exception: err.Error(),
		}
	case err != nil:
		return &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to delete major",
			Exception: err.Error(),
		}
	default:
		return nil
	}
}

func DeleteUnderGradMajor(ctx context.Context, repository model.AcademicLogRepository, astronautID, majorID int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := repository.DeleteAstronautUnderGradMajor(ctx, astronautID, majorID)
	if err != nil {
		return &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to delete astronaut undergrad major",
			Exception: err.Error(),
		}
	}
	return nil
}

func DeleteGradeMajor(ctx context.Context, repository model.AcademicLogRepository, astronautID, majorID int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := repository.DeleteAstronautGradMajor(ctx, astronautID, majorID)
	if err != nil {
		return &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to delete astronaut grad major",
			Exception: err.Error(),
		}
	}
	return nil
}

func DeleteAlmaMater(ctx context.Context, repository model.AcademicLogRepository, id int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := repository.DeleteAlmaMater(ctx, id)
	switch {
	case errors.Is(err, model.ErrNoChange):
		return &model.APIError{
			Code:      http.StatusBadRequest,
			Message:   "Alma Mater not found",
			Exception: err.Error(),
		}
	case err != nil:
		return &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to delete Alma Mater",
			Exception: err.Error(),
		}
	default:
		return nil
	}
}

func DeleteAstronautAlmaMater(ctx context.Context, repository model.AcademicLogRepository, astronautID, almaMaterID int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := repository.DeleteAstronautAlmaMater(ctx, astronautID, almaMaterID); err != nil {
		return &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to delete astronaut alma mater",
			Exception: err.Error(),
		}
	}
	return nil
}

func GetAstronautAcademicLog(ctx context.Context, repository model.AcademicLogRepository, astronautID int) (*model.AcademicLog, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	al, err := repository.GetAcademicLog(ctx, astronautID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, &model.APIError{
			Code:      http.StatusNotFound,
			Message:   "Astronaut academic log not found",
			Exception: err.Error(),
		}
	case err != nil:
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to get astronat academic log",
			Exception: err.Error(),
		}
	default:
		return al, nil
	}
}
