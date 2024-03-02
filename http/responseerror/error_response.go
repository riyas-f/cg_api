package responseerror

type FailedRequestResponse struct {
	Status    string `json:"status"`
	ErrorType string `json:"error_type"`
	Message   string `json:"message"`
}
