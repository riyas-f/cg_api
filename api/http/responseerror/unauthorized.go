package responseerror

const (
	InvalidAuthHeader errorType = "invalid_auth_header"
	EmptyAuthHeader   errorType = "empty_auth_header"
	InvalidToken      errorType = "invalid_token"
	TokenExpired      errorType = "token_expired"
	RefreshDenied     errorType = "refresh_denied"
	ClaimsMismatch    errorType = "claims_mismatch"
	AccessDenied      errorType = "access_denied"
)

const (
	InvalidAuthHeaderMessage errorMessageTemplate = "Not accepted authorization of type {{.authType}}"
	EmptyAuthHeaderMessage   errorMessageTemplate = "Required authorization header in request header"
	InvalidTokenMessage      errorMessageTemplate = "Invalid Token.{{.description}}"
	TokenExpiredMessage      errorMessageTemplate = "your token has expired"
	RefreshDeniedMessage     errorMessageTemplate = "cannot get new access token when the previous one still active"
	ClaimsMismatchMessage    errorMessageTemplate = "refresh claims and username claims don't share the same credentials"
	AccessDeniedMessage      errorMessageTemplate = "cannot establish trust. certificate is either invalid, revoked, or empty"
)

func CreateUnauthorizedError(t errorType, tmp errorMessageTemplate, namedArgs map[string]string) HTTPCustomError {
	return &ResponseError{
		Code:    403,
		Message: ParseMessage(tmp, namedArgs),
		Name:    string(t),
	}
}
