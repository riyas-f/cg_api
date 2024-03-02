package responseerror

const (
	InternalServiceErr errorType = "internal_service_error"
)

const (
	InternalServiceErrorTemplate errorMessageTemplate = "sorry, we cannot proceed with your request at the moment. please try again later!. "
)

type InternalServiceError struct {
	*ResponseError
	Description string
}

func CreateInternalServiceError(err error) HTTPCustomError {
	return &InternalServiceError{
		ResponseError: &ResponseError{
			Code:    500,
			Message: string(InternalServiceErrorTemplate),
			Name:    string(InternalServiceErr),
		},
		Description: err.Error(),
	}
}
