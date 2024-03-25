package payload

import (
	"io"
	"strconv"

	"github.com/AdityaP1502/Instant-Messanging/api/http/httputil"
	"github.com/AdityaP1502/Instant-Messanging/api/http/responseerror"
	"github.com/AdityaP1502/Instant-Messanging/api/jsonutil"
	"github.com/AdityaP1502/Instant-Messanging/api/service/account/pwdutil"
)

type Account struct {
	AccountID string `json:"-" db:"account_id"`
	Username  string `json:"username" db:"username"`
	Name      string `json:"name" db:"name"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"password" db:"password"`
	Salt      string `json:"-" db:"password_salt"`
	IsActive  string `json:"-" db:"is_active"`
}

func (acc *Account) FromJSON(r io.Reader, checkRequired bool, requiredFields []string) responseerror.HTTPCustomError {
	err := jsonutil.DecodeJSON(r, acc)
	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	if checkRequired {
		return httputil.CheckParametersUnity(acc, requiredFields)
	}

	return nil
}

func (acc *Account) ToJSON(checkRequired bool, requiredFields []string) ([]byte, error) {
	var err error

	var tmp struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
	}

	tmp.Username = acc.Username
	tmp.Name = acc.Name
	tmp.Email = acc.Email

	if checkRequired {
		if err = httputil.CheckParametersUnity(&tmp, requiredFields); err != nil {
			return nil, err
		}
	}

	return jsonutil.EncodeToJson(&tmp)
}

func NewRegisteredAccountPayload(username string, name string, email string, password string, secretKey []byte) (*Account, error) {
	hash, salt, err := pwdutil.HashPassword(password, secretKey)

	if err != nil {
		return nil, err
	}

	return &Account{
		Username: username,
		Email:    email,
		Name:     name,
		Password: hash,
		Salt:     salt,
		IsActive: strconv.FormatBool(false),
	}, nil
}
