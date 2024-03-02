package payload

import (
	"io"

	"github.com/AdityaP1502/Instant-Messanging/api/http/httputil"
	"github.com/AdityaP1502/Instant-Messanging/api/http/responseerror"
	"github.com/AdityaP1502/Instant-Messanging/api/jsonutil"
	"github.com/AdityaP1502/Instant-Messanging/api/service/auth/config"
	"github.com/AdityaP1502/Instant-Messanging/api/service/auth/jwtutil"
)

type Token struct {
	RefreshToken string `json:"refresh_token,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
}

func (t *Token) FromJSON(r io.Reader, checkRequired bool, requiredFields []string) responseerror.HTTPCustomError {
	err := jsonutil.DecodeJSON(r, t)

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	if checkRequired {
		return httputil.CheckParametersUnity(t, requiredFields)
	}

	return nil
}

func (t *Token) ToJSON(checkRequired bool, requiredFields []string) ([]byte, error) {
	var err error

	if checkRequired {
		if err = httputil.CheckParametersUnity(t, requiredFields); err != nil {
			return nil, err
		}
	}

	return jsonutil.EncodeToJson(t)
}

func (t *Token) GenerateTokenPair(config *config.Config, username string, email string, role jwtutil.Roles) error {
	claims := jwtutil.GenerateClaims(config, username, email, role)
	refreshClaism := jwtutil.GenerateRefreshClaims(config, username, email, role)

	accessToken, err := jwtutil.GenerateToken(claims, config.Session.SecretKeyRaw)

	if err != nil {
		return err
	}

	refreshToken, err := jwtutil.GenerateToken(refreshClaism, config.Session.SecretKeyRaw)
	if err != nil {
		return err
	}

	t.AccessToken = accessToken
	t.RefreshToken = refreshToken

	return err
}

func (t *Token) CheckRefreshEligibility(config *config.Config) (*jwtutil.Claims, error) {
	claims, err := jwtutil.VerifyToken(t.AccessToken, config.Session.SecretKeyRaw)

	if err == nil {
		return nil, responseerror.CreateUnauthorizedError(
			responseerror.RefreshDenied,
			responseerror.RefreshDeniedMessage,
			nil,
		)
	}

	if err.Get().Name == string(responseerror.InvalidToken) {
		return nil, err
	}

	refreshClaims, err := jwtutil.VerifyToken(t.RefreshToken, config.Session.SecretKeyRaw)
	if err != nil {
		return nil, err
	}

	if claims.Username != refreshClaims.Username && claims.Email != refreshClaims.Email {
		return nil, responseerror.CreateUnauthorizedError(
			responseerror.ClaimsMismatch,
			responseerror.ClaimsMismatchMessage,
			nil,
		)
	}

	return claims, nil
}
