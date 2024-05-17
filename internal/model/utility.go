package model

type APIError struct {
	Message   string
	Code      int
	Exception string
}

func (e *APIError) Error() string {
	return e.Exception
}

type Validator interface {
	Valid() (map[string]string, bool)
}
