package responseerror

import (
	"net/http"
)

const (
	MissingParameter     errorType = "missing_parameter"
	HeaderValueMistmatch errorType = "header_value_mismatch"
	NonUniqueValue       errorType = "non_unique_value"
	UsernameInvalid      errorType = "username_invalid"
	PasswordWeak         errorType = "password_weak"
	EmailInvalid         errorType = "invalid_email"
	PayloadInvalid       errorType = "invalid_payload"
	OTPInvalid           errorType = "invalid_otp"
	OTPExpired           errorType = "otp_expired"
	CallsStatusInvalid   errorType = "invalid_status_value"
	SteamNotLinked       errorType = "steam_not_linked"
	SteamAlreadyLinked   errorType = "steam_has_been_linked"
	MalformedSessionID   errorType = "malformed_uuidv7_id"
	InvalidPIN           errorType = "invalid_pin"
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
	OTPExpiredMessage           errorMessageTemplate = "otp is expired. please resend a new otp to your email"
	StatusInvalidMessage        errorMessageTemplate = "trying to update status to {{.reqStatus}} when user status is {{.status}}"
	SteamAlreadyLinkedMessage   errorMessageTemplate = "this account has been linked to steam"
	SteamNotLinkedMessage       errorMessageTemplate = "account has not been linked to steam yet"
	MalformedSessionIDMessage   errorMessageTemplate = "{{.id}} cannot be parsed as a valid uuidv7 id "
	InvalidPINMessage           errorMessageTemplate = "{{.pin}} is invalid"
	NonUniqueValueMessage       errorMessageTemplate = "{{.message}}"
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
