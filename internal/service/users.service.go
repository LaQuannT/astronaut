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

func RegisterUser(ctx context.Context, repository model.UserRepository, user *model.User) (*model.User, error) {
	err := validate(user, "User")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user.Password, err = generatePasswordHash(user.Password)
	if err != nil {
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to register user",
			Exception: err.Error(),
		}
	}

	err = repository.CreateUser(ctx, user)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, &model.APIError{
				Code:      http.StatusBadRequest,
				Message:   "Email already in use",
				Exception: pgErr.Message,
			}
		}
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to register User",
			Exception: err.Error(),
		}
	}
	return user, nil
}

func SearchUserID(ctx context.Context, repository model.UserRepository, id int) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	u, err := repository.FindUserByID(ctx, id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, &model.APIError{
			Code:      http.StatusNotFound,
			Message:   "User not found",
			Exception: err.Error(),
		}
	case err != nil:
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to find User",
			Exception: err.Error(),
		}
	default:
		return u, nil
	}
}

func SearchUserEmail(ctx context.Context, repository model.UserRepository, email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	u, err := repository.FindUserByEmail(ctx, email)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, &model.APIError{
			Code:      http.StatusNotFound,
			Message:   "User not found",
			Exception: err.Error(),
		}
	case err != nil:
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to find User",
			Exception: err.Error(),
		}
	default:
		return u, nil
	}
}

func GetUsers(ctx context.Context, repository model.UserRepository) ([]*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	us, err := repository.FindAllUsers(ctx)
	if err != nil {
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to find Users",
			Exception: err.Error(),
		}
	}
	return us, nil
}

func UpdateUser(ctx context.Context, repository model.UserRepository, user *model.User) error {
	if err := validate(user, "User"); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := repository.UpdateUser(ctx, user)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return &model.APIError{
				Code:      http.StatusBadRequest,
				Message:   "Email already in use",
				Exception: pgErr.Message,
			}
		}

		if errors.Is(err, model.ErrNoChange) {
			return &model.APIError{
				Code:      http.StatusBadRequest,
				Message:   "User not found",
				Exception: err.Error(),
			}
		}

		return &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to register User",
			Exception: err.Error(),
		}
	}
	return nil
}

func DeleteUser(ctx context.Context, repository model.UserRepository, id int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := repository.DeleteUser(ctx, id)
	switch {
	case errors.Is(err, model.ErrNoChange):
		return &model.APIError{
			Code:      http.StatusBadRequest,
			Message:   "User not found",
			Exception: err.Error(),
		}
	case err != nil:
		return &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to delete User",
			Exception: err.Error(),
		}
	default:
		return nil

	}
}

func ResetPassword(ctx context.Context, repository model.UserRepository, password string, userID int) error {
	if err := model.ValidatePassword(password); err != nil {
		return &model.APIError{
			Code:      http.StatusBadRequest,
			Message:   err.Error(),
			Exception: err.Error(),
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	hash, err := generatePasswordHash(password)
	if err != nil {
		return &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to reset user password",
			Exception: err.Error(),
		}
	}

	err = repository.RestUserPassword(ctx, hash, userID)
	switch {
	case errors.Is(err, model.ErrNoChange):
		return &model.APIError{
			Code:      http.StatusNotFound,
			Message:   "User not found",
			Exception: err.Error(),
		}

	case err != nil:
		return &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to reset password",
			Exception: err.Error(),
		}
	default:
		return nil
	}
}

func GenerateNewAPIKey(ctx context.Context, repository model.UserRepository, userID int) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	key, err := repository.GenerateNewUserAPIKey(ctx, userID)
	if err != nil {
		return "", &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to generate new API key",
			Exception: err.Error(),
		}
	}
	return key, nil
}

func CreateAdmin(ctx context.Context, repository model.UserRepository, userID int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := repository.GiveAdminPrivileges(ctx, userID); err != nil {
		return &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to add admin",
			Exception: err.Error(),
		}
	}
	return nil
}

func RemoveAdmin(ctx context.Context, repository model.UserRepository, userID int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := repository.RevokeAdminPrivileges(ctx, userID); err != nil {
		if errors.Is(err, model.ErrNoChange) {
			return &model.APIError{
				Code:      http.StatusBadRequest,
				Message:   "Admin not found",
				Exception: err.Error(),
			}
		}
		return &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to remove admin",
			Exception: err.Error(),
		}
	}
	return nil
}

func SearchAPIKey(ctx context.Context, repository model.UserRepository, key string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	usr, err := repository.FindUserByAPIKey(ctx, key)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, &model.APIError{
			Code:      http.StatusNotFound,
			Message:   "User not found",
			Exception: err.Error(),
		}
	case err != nil:
		return nil, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to get user",
			Exception: err.Error(),
		}
	default:
		return usr, nil
	}
}

func CheckAdminPermission(ctx context.Context, repository model.UserRepository, userID int) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	usrCount, err := repository.IsAdmin(ctx, userID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return false, nil
	case err != nil:
		return false, &model.APIError{
			Code:      http.StatusInternalServerError,
			Message:   "failed to check user permission",
			Exception: err.Error(),
		}
	}

	if usrCount != 1 {
		return false, nil
	}

	return true, nil
}
