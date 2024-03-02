package payload

import (
	"io"

	"github.com/AdityaP1502/Instant-Messanging/api/http/httputil"
	"github.com/AdityaP1502/Instant-Messanging/api/http/responseerror"
	"github.com/AdityaP1502/Instant-Messanging/api/jsonutil"
)

type Status string

var (
	statusMap = map[string]Status{
		string(OnCall): OnCall,
		string(OnWait): OnWait,
		string(Free):   Free,
	}
)

func ParseStringToStatus(str string) (Status, bool) {
	st, ok := statusMap[str]

	return st, ok
}

var (
	OnCall Status = "on_call"
	OnWait Status = "on_wait"
	Free   Status = "free"
)

type CallStatus struct {
	Status             string `json:"status" db:"call_status"`
	Username           string `json:"sender_username" db:"username"`
	OtherPartyUsername string `json:"recipient_username" db:"-"`
}

func CreateCallStatus(status Status, username string, otherParty string) *CallStatus {
	return &CallStatus{
		Status:             string(status),
		Username:           username,
		OtherPartyUsername: otherParty,
	}
}

func (c *CallStatus) FromJSON(r io.Reader, checkRequired bool, requiredFields []string) responseerror.HTTPCustomError {
	err := jsonutil.DecodeJSON(r, c)

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	if checkRequired {
		return httputil.CheckParametersUnity(c, requiredFields)
	}

	return nil
}

func (c *CallStatus) ToJSON(checkRequired bool, requiredFields []string) ([]byte, error) {
	var err error

	if checkRequired {
		if err = httputil.CheckParametersUnity(c, requiredFields); err != nil {
			return nil, err
		}
	}

	return jsonutil.EncodeToJson(c)
}
