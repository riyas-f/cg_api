package payload

import (
	"io"
	"strconv"

	"github.com/AdityaP1502/Instant-Messanging/api/date"
	"github.com/AdityaP1502/Instant-Messanging/api/http/httputil"
	"github.com/AdityaP1502/Instant-Messanging/api/http/responseerror"
	"github.com/AdityaP1502/Instant-Messanging/api/jsonutil"
	"github.com/AdityaP1502/Instant-Messanging/api/service/account/otp"
)

type UserOTP struct {
	OTPID             string `json:"-" db:"otp_id"`
	Email             string `json:"email" db:"email"`
	OTP               string `json:"otp" db:"otp"`
	LastResend        string `json:"-" db:"last_resend"`
	ExpiredAt         string `json:"-" db:"expired_at"`
	MarkedForDeletion string `json:"-" db:"marked_for_deletion"`
}

func (o *UserOTP) FromJSON(r io.Reader, checkRequired bool, requiredFields []string) responseerror.HTTPCustomError {
	err := jsonutil.DecodeJSON(r, o)

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	if checkRequired {
		return httputil.CheckParametersUnity(o, requiredFields)
	}

	return nil
}

func (o *UserOTP) ToJSON(checkRequired bool, requiredFields []string) ([]byte, error) {
	var err error

	if checkRequired {
		if err = httputil.CheckParametersUnity(o, requiredFields); err != nil {
			return nil, err
		}
	}

	return jsonutil.EncodeToJson(o)
}

func NewOTPPayload(email string, otpTTL int) (*UserOTP, error) {
	otp, err := otp.GenerateOTP()

	if err != nil {
		return nil, err
	}

	return &UserOTP{
		Email:             email,
		OTP:               otp,
		ExpiredAt:         date.GenerateTimestampWithOffset(otpTTL),
		LastResend:        date.GenerateTimestamp(),
		MarkedForDeletion: strconv.FormatBool(false),
	}, nil

}
