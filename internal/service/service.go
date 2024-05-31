package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/LaQuannT/astronaut-api/internal/model"
	"golang.org/x/crypto/bcrypt"
)

const hashingCost = 12

func validate(validator model.Validator, name string) error {
	problems, isValid := validator.Valid()
	if !isValid {
		msg := fmt.Sprintf("Invalid %s input", name)
		for _, problem := range problems {
			msg += fmt.Sprintf("; %s", problem)
		}
		return &model.APIError{
			Code:      http.StatusBadRequest,
			Message:   msg,
			Exception: fmt.Sprintf("Invalid %s input", name),
		}
	}
	return nil
}

func generatePasswordHash(pwd string) (string, error) {
	if pwd == "" {
		return "", errors.New("password not provided")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), hashingCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func validatePasswordHash(hash, plain string) bool {
	if hash == "" || plain == "" {
		return false
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)); err != nil {
		return false
	}
	return true
}
