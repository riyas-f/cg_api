package httpx

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"

	"github.com/AdityaP1502/Instant-Messanging/api/http/responseerror"
	"github.com/AdityaP1502/Instant-Messanging/api/jsonutil"
)

type HTTPRequest struct {
	Request            http.Request
	Payload            []byte
	SuccessStatusCode  int
	ReturnedStatusCode int
	Status             int
	IsSuccess          bool
	TLSClientConfig    *tls.Config
}

func (h *HTTPRequest) CreateRequest(scheme string, host string, port int, endpoint string, method string, successStatus int, payload interface{}, tlsConfig *tls.Config) (*HTTPRequest, responseerror.HTTPCustomError) {
	url := fmt.Sprintf("%s://%s:%d/%s", scheme, host, port, endpoint)

	json, err := jsonutil.EncodeToJson(payload)
	if err != nil {
		return nil, responseerror.CreateInternalServiceError(err)
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(json))
	if err != nil {
		return nil, responseerror.CreateInternalServiceError(err)
	}

	req.Header.Set("Content-Type", "application/json")

	return &HTTPRequest{
		Request:           *req,
		Payload:           nil,
		SuccessStatusCode: successStatus,
		TLSClientConfig:   tlsConfig,
	}, nil
}

func (h *HTTPRequest) Send(dest interface{}) responseerror.HTTPCustomError {
	var client = &http.Client{}

	if h.Request.URL.Scheme == "https" {
		client.Transport = &http.Transport{
			TLSClientConfig: h.TLSClientConfig,
		}
	}

	resp, err := client.Do(&h.Request)

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != h.SuccessStatusCode {
		// if not provided a destination or that the status code don't match
		// the expected return code
		// store the payload in the payload field

		respBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			return responseerror.CreateInternalServiceError(err)
		}

		errorResponse := &responseerror.FailedRequestResponse{}
		err = jsonutil.DecodeJSON(bytes.NewReader(respBytes), errorResponse)

		if err != nil {
			return responseerror.CreateInternalServiceError(err)
		}

		fmt.Println(errorResponse)

		return &responseerror.ResponseError{
			Code:    resp.StatusCode,
			Message: errorResponse.Message,
			Name:    errorResponse.ErrorType,
		}
	}

	if dest == nil {
		return nil
	}

	err = jsonutil.DecodeJSON(resp.Body, dest)

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	return nil
}
