package responseerror

import (
	"net/http"
)

const (
	MissingParameter     errorType = "missing_parameter"
	HeaderValueMistmatch errorType = "header_value_mismatch"
	UsernameExists       errorType = "username_exists"
	EmailExists          errorType = "email_exists"
	UsernameInvalid      errorType = "username_invalid"
	PasswordWeak         errorType = "password_weak"
	EmailInvalid         errorType = "invalid_email"
	PayloadInvalid       errorType = "invalid_payload"
	OTPInvalid           errorType = "invalid_otp"
	CallsStatusInvalid   errorType = "invalid_status_value"
)

const (
	MissingParameterMessage     errorMessageTemplate = "required field {{.field}} is missing"
	HeaderValueMistmatchMessage errorMessageTemplate = "mismatch value in header {{.name}}"
	UsernameExistMessage        errorMessageTemplate = "email already taken"
	EmailsExistMessage          errorMessageTemplate = "username already taken"
	UsernameInvalidMessage      errorMessageTemplate = "username is invalid"
	PasswordWeakMessage         errorMessageTemplate = "password is weak"
	EmailInvalidMessage         errorMessageTemplate = "email is invalid"
	PayloadInvalidMessage       errorMessageTemplate = "payload is invalid"
	OTPInvalidMessage           errorMessageTemplate = "otp is invalid"
	StatusInvalidMessage        errorMessageTemplate = "trying to update status to {{.reqStatus}} when user status is {{.status}}"
)

// create response error with 400 Code and name string(t)
// the error message will be formatted from tmp and namedArgs using
// ParseMessage(tmp, namedArgs)
//
// refer to ParseMessage docs for more details
func CreateBadRequestError(t errorType, tmp errorMessageTemplate, namedArgs map[string]string) HTTPCustomError {
	return &ResponseError{
		Code:    http.StatusBadRequest,
		Name:    string(t),
		Message: ParseMessage(tmp, namedArgs),
	}
}
