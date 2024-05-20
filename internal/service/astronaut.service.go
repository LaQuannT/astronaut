package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/LaQuannT/astronaut-api/internal/model"
	"net/http"
	"time"
)

func AddAstronaut(ctx context.Context, a *model.Astronaut, r model.AstronautRepository) (*model.Astronaut, error) {
	problems, isValid := a.Valid()
	if !isValid {
		msg := "Invalid Astronaut input"
		for _, problem := range problems {
			msg = fmt.Sprintf("%s; %s", msg, problem)
		}
		return nil, &model.APIError{
			Code:    http.StatusBadRequest,
			Message: msg,
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := r.CreateAstronaut(ctx, a)
	switch {
	case err != nil:
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to create astronaut",
			Exception: err.Error(),
		}
	default:
		return a, nil
	}
}

func GetAstronaut(
	ctx context.Context,
	ar model.AstronautRepository,
	id int,
) (*model.Astronaut, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	astronaut, err := ar.FindAstronautByID(ctx, id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, &model.APIError{
			Code:      http.StatusNotFound,
			Message:   "not found",
			Exception: err.Error(),
		}
	case err != nil:
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to retrieve astronaut",
			Exception: err.Error(),
		}
	default:
		return astronaut, nil
	}
}

func GetAstronauts(ctx context.Context, r model.AstronautRepository) ([]*model.Astronaut, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	astronauts, err := r.FindAstronauts(ctx)
	if err != nil {
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to retrieve astronauts",
			Exception: err.Error(),
		}
	}

	return astronauts, nil
}

func UpdateAstronaut(ctx context.Context, a *model.Astronaut, r model.AstronautRepository) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	problems, isValid := a.Valid()
	if !isValid {
		msg := "Invalid Astronaut input"
		for _, problem := range problems {
			msg = msg + "; " + problem
		}
		return &model.APIError{
			Code:    http.StatusBadRequest,
			Message: msg,
		}
	}

	err := r.UpdateAstronaut(ctx, a)
	if err != nil {
		return &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to update astronaut",
			Exception: err.Error(),
		}
	}
	return nil
}

func DeleteAstronaut(ctx context.Context, r model.AstronautRepository, id int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := r.DeleteAstronaut(ctx, id)
	switch {
	case errors.Is(err, model.ErrNoChange):
		return &model.APIError{
			Code:      http.StatusNotFound,
			Message:   "not found",
			Exception: err.Error(),
		}

	case err != nil:
		return &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to delete astronaut",
			Exception: err.Error(),
		}
	}
	return nil
}

func SearchAstronautByName(ctx context.Context, r model.AstronautRepository, name string) ([]*model.Astronaut, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	astronauts, err := r.FindAstronautByName(ctx, name)
	if err != nil {
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to search astronauts",
			Exception: err.Error(),
		}
	}
	if len(astronauts) == 0 {
		return nil, nil
	}
	return astronauts, nil
}
