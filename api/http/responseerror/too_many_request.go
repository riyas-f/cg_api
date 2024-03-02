package responseerror

const (
	ResendIntervalNotReachedErr errorType = "otp_resend_interval_not_reached"
)

const (
	ResendIntervalNotReachedMessage errorMessageTemplate = "Mail has already been sent to your registered email"
)

func CreateTooManyRequestError(t errorType, tmp errorMessageTemplate, namedArgs map[string]string) HTTPCustomError {
	return &ResponseError{
		Code:    429,
		Message: ParseMessage(tmp, namedArgs),
		Name:    string(t),
	}
}
