package responseerror

import (
	"net/http"
)

type NotFoundError struct {
	Name    string
	Code    int
	Message string
}

const (
	NotFoundErr errorType = "resource_not_found"
)

const (
	NotFoundTemplate errorMessageTemplate = "{{.resourceName}} not found"
)

// Create ResponseError with 404 code and NotFoundErr errorType
//
// and NotFoundTemplate : "{{.resourceName}} not found"
func CreateNotFoundError(namedArgs map[string]string) HTTPCustomError {
	return &ResponseError{
		Code:    http.StatusNotFound,
		Message: ParseMessage(NotFoundTemplate, namedArgs),
		Name:    string(NotFoundErr),
	}
}
