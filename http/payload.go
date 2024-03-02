package httpx

import (
	"io"

	"github.com/AdityaP1502/Instant-Messanging/api/http/responseerror"
)

type Payload interface {
	FromJSON(r io.Reader, checkRequired bool, requiredFields []string) responseerror.HTTPCustomError
	ToJSON(checkRequired bool, requiredFields []string) ([]byte, error)
}
