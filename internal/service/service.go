package service

import (
	"fmt"
	"github.com/LaQuannT/astronaut-api/internal/model"
	"net/http"
)

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
