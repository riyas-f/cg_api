package payload

import (
	"io"
	"strconv"

	"github.com/AdityaP1502/Instant-Messanging/api/date"
	"github.com/AdityaP1502/Instant-Messanging/api/http/httputil"
	"github.com/AdityaP1502/Instant-Messanging/api/http/responseerror"
	"github.com/AdityaP1502/Instant-Messanging/api/jsonutil"
	"github.com/AdityaP1502/Instant-Messanging/api/service/account/otp"
	"github.com/google/uuid"
)

type UserOTP struct {
	OTPID             string `json:"-" db:"otp_id"`
	Username          string `json:"-" db:"username"`
	OTPConfirmID      string `json:"otp_confirmation_id" db:"otp_confirmation_id"`
	OTP               string `json:"otp" db:"otp"`
	LastResend        string `json:"-" db:"last_resend"`
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

func NewOTPPayload(username string) (*UserOTP, error) {
	otp, err := otp.GenerateOTP()

	if err != nil {
		return nil, err
	}

	return &UserOTP{
		Username:          username,
		OTPConfirmID:      uuid.NewString(),
		OTP:               otp,
		LastResend:        date.GenerateTimestamp(),
		MarkedForDeletion: strconv.FormatBool(false),
	}, nil

}
