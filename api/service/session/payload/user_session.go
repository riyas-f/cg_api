package payload

import (
	"io"

	"github.com/AdityaP1502/Instant-Messanging/api/http/httputil"
	"github.com/AdityaP1502/Instant-Messanging/api/http/responseerror"
	"github.com/AdityaP1502/Instant-Messanging/api/jsonutil"
)

type UserSession struct {
	SessionID []byte `json:"-" db:"session_id"`
	// SessionIDString   string          `json:"session_id" db:"-"`
	Username          string          `json:"username" db:"username"`
	RequestStatus     string          `json:"status" db:"request_status"`
	LastUpdatedAt     string          `json:"last_updated_at" db:"last_update"`
	SessionMetadata   SessionMetadata `json:"session_metadata" db:"-"` // disable recurse
	MarkedForDeletion string          `json:"-" db:"marked_for_deletion"`
}

func (c *UserSession) FromJSON(r io.Reader, checkRequired bool, requiredFields []string) responseerror.HTTPCustomError {
	err := jsonutil.DecodeJSON(r, c)

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	if checkRequired {
		return httputil.CheckParametersUnity(c, requiredFields)
	}

	return nil
}

func (c *UserSession) ToJSON(checkRequired bool, requiredFields []string) ([]byte, error) {
	var err error

	if checkRequired {
		if err = httputil.CheckParametersUnity(c, requiredFields); err != nil {
			return nil, err
		}
	}

	return jsonutil.EncodeToJson(c)
}
