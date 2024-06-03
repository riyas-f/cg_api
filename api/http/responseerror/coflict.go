package responseerror

const (
	UpdateConflictErr errorType = "update_conflict_encountered"
)

const (
	UpdateConflictErrorMessage errorMessageTemplate = "encountered conflict on update"
)

func CreateConflictError(t errorType, tmp errorMessageTemplate, namedArgs map[string]string) HTTPCustomError {
	return &ResponseError{
		Code:    409,
		Message: ParseMessage(tmp, namedArgs),
		Name:    string(t),
	}
}
