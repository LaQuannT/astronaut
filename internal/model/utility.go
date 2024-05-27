package model

import (
	"errors"
	"regexp"
	"unicode"
)

var ErrNoChange = errors.New("sql: no rows effected")
var errInvalidPassword = errors.New("password must be at least 8 characters containing at least one uppercase, lowercase, number and special character")
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type APIError struct {
	Message   string
	Code      int
	Exception string
}

func (e APIError) Error() string {
	return e.Exception
}

type Validator interface {
	Valid() (map[string]string, bool)
}

func ValidatePassword(password string) error {
	minPasswordLength := 8

	var hasUpper, hasLower, hasSpecial, hasNumber bool

	if len(password) < minPasswordLength {
		return errInvalidPassword
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasLower || !hasUpper || !hasNumber || !hasSpecial {
		return errInvalidPassword
	}
	return nil

}
