package payload

import (
	"io"

	"github.com/AdityaP1502/Instant-Messanging/api/http/httputil"
	"github.com/AdityaP1502/Instant-Messanging/api/http/responseerror"
	"github.com/AdityaP1502/Instant-Messanging/api/jsonutil"
)

type UserGames struct {
	Username string `db:"username" json:"username"`
	GameID   int    `db:"game_id" json:"game_id"`
	GPU      string `json:"gpu"`
}

type UserCollections struct {
	CollectionsID  string `json:"-" db:"collections_id"`
	GameID         int    `json:"game_id,omitempty" db:"game_id"`
	GameName       string `json:"name" db:"name"`
	GamePictureURL string `json:"game_picture_url" db:"display_picture_url"`
}

func (c *UserGames) FromJSON(r io.Reader, checkRequired bool, requiredFields []string) responseerror.HTTPCustomError {
	err := jsonutil.DecodeJSON(r, c)

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	if checkRequired {
		return httputil.CheckParametersUnity(c, requiredFields)
	}

	return nil
}

func (c *UserGames) ToJSON(checkRequired bool, requiredFields []string) ([]byte, error) {
	var err error

	if checkRequired {
		if err = httputil.CheckParametersUnity(c, requiredFields); err != nil {
			return nil, err
		}
	}

	return jsonutil.EncodeToJson(c)
}
