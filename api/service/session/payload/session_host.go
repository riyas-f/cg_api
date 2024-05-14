package payload

import (
	"io"

	"github.com/AdityaP1502/Instant-Messanging/api/http/httputil"
	"github.com/AdityaP1502/Instant-Messanging/api/http/responseerror"
	"github.com/AdityaP1502/Instant-Messanging/api/jsonutil"
)

type Webhook struct {
	Host    string `json:"host" db:"webhook_host"`
	HostAlt string `json:"host_alt" db:"webhook_host_alt"`
	Port    string `json:"port" db:"webhook_port"`
}

type SessionHost struct {
	Webhook         Webhook `json:"webhook"`
	NetworkID       string  `json:"network_id" db:"network_id"`
	SessionIDString string  `json:"session_id" db:"-"`
	SessionID_      []byte  `json:"-" db:"session_id"` // avoid confilict with UserSession struct
}

func (c *SessionHost) FromJSON(r io.Reader, checkRequired bool, requiredFields []string) responseerror.HTTPCustomError {
	err := jsonutil.DecodeJSON(r, c)

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	if checkRequired {
		return httputil.CheckParametersUnity(c, requiredFields)
	}

	return nil
}

func (c *SessionHost) ToJSON(checkRequired bool, requiredFields []string) ([]byte, error) {
	var err error

	if checkRequired {
		if err = httputil.CheckParametersUnity(c, requiredFields); err != nil {
			return nil, err
		}
	}

	return jsonutil.EncodeToJson(c)
}
