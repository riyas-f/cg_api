package responseerror

const (
	UserMarkedInActive errorType = "user_marked_inactive"
	InvalidCredentials errorType = "invalid_credentials"
)

const (
	UserMarkedInActiveMessage errorMessageTemplate = "your account hasn't been verified and currently marked inactive"
	InvalidCredentialsMessage errorMessageTemplate = "email or password is wrong"
)

func CreateUnauthenticatedError(t errorType, tmp errorMessageTemplate, namedArgs map[string]string) HTTPCustomError {
	return &ResponseError{
		Code:    401,
		Message: ParseMessage(tmp, namedArgs),
		Name:    string(t),
	}
}
