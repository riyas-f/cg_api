package payload

import (
	"io"

	"github.com/AdityaP1502/Instant-Messanging/api/http/httputil"
	"github.com/AdityaP1502/Instant-Messanging/api/http/responseerror"
	"github.com/AdityaP1502/Instant-Messanging/api/jsonutil"
)

type GPU struct {
	GPUID          string `json:"-" db:"gpu_id"`
	Version        string `db:"version" json:"-"`
	GPUName        string `json:"gpu_name" db:"gpu_name"`
	GPUDisplayName string `json:"gpu_display_name" db:"gpu_alt_name"`
	Count          string `json:"n_available" db:"n_available"`
}

func (c *GPU) FromJSON(r io.Reader, checkRequired bool, requiredFields []string) responseerror.HTTPCustomError {
	err := jsonutil.DecodeJSON(r, c)

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	if checkRequired {
		return httputil.CheckParametersUnity(c, requiredFields)
	}

	return nil
}

func (c *GPU) ToJSON(checkRequired bool, requiredFields []string) ([]byte, error) {
	var err error

	if checkRequired {
		if err = httputil.CheckParametersUnity(c, requiredFields); err != nil {
			return nil, err
		}
	}

	return jsonutil.EncodeToJson(c)
}
