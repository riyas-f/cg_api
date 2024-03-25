package payload

import (
	"strconv"
	"testing"
	"time"

	"github.com/AdityaP1502/Instant-Messanging/api/database"
	"github.com/AdityaP1502/Instant-Messanging/api/service/account/otp"
	"github.com/google/uuid"
)

func TestInserUserOTP(t *testing.T) {
	db := connectToDB(t)

	querynator := database.Querynator{}
	_, err := querynator.Insert(newUser, db.DB, "account", "")

	if err != nil {
		t.Error(err)
		return
	}

	otp, err := otp.GenerateOTP()
	if err != nil {
		t.Error(err)
		return
	}

	otpData := &UserOTP{
		Username:          newUser.Username,
		OTPConfirmID:      uuid.NewString(),
		OTP:               otp,
		LastResend:        time.Now().Format(time.RFC3339),
		MarkedForDeletion: strconv.FormatBool(false),
	}

	otpID, err := querynator.Insert(otpData, db.DB, "user_otp", "otp_id")

	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("Success. The ID is %d", otpID)
}

func TestUpdateOTP(t *testing.T) {
	db := connectToDB(t)

	data := &Account{
		Username: "Dimwit",
		Email:    "email@domain.com",
		Name:     "MyGuyisGay",
		Salt:     "random_ah_salt",
		Password: "Test123",
		IsActive: strconv.FormatBool(false),
	}

	querynator := database.Querynator{}
	_, err := querynator.Insert(data, db.DB, "account", "")

	if err != nil {
		t.Error(err)
		return
	}

	otpData, err := NewOTPPayload(data.Username)

	if err != nil {
		t.Error(err)
		return
	}

	_, err = querynator.Insert(otpData, db.DB, "user_otp", "")

	if err != nil {
		t.Error(err)
		return
	}

	searchData := &Account{}
	err = querynator.FindOne(&Account{Username: data.Username}, searchData, db.DB, "account", "username")

	if err != nil {
		t.Error(err)
		return
	}

	o, err := otp.GenerateOTP()

	if err != nil {
		t.Error(err)
		return
	}

	now := time.Now().Format(time.RFC3339)
	loc := time.Now().Location()

	err = querynator.Update(&UserOTP{
		OTP:        o,
		LastResend: now,
	}, []string{"username", "otp_confirmation_id"}, []any{searchData.Username, otpData.OTPConfirmID}, db.DB, "user_otp")

	if err != nil {
		t.Error(err)
		return
	}

	searchOTP := &UserOTP{}

	err = querynator.FindOne(&UserOTP{OTPConfirmID: otpData.OTPConfirmID}, searchOTP, db.DB, "user_otp", "otp", "last_resend")

	if err != nil {
		t.Error(err)
		return
	}

	if searchOTP.OTP != o {
		t.Errorf("Data hasn't been updated properly")
		return
	}

	resendTime, err := time.Parse(time.RFC3339, searchOTP.LastResend)

	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("Last resend in: %s", now)
	t.Logf("Last resend in: %s", resendTime.In(loc).Format(time.RFC3339))
	t.Log("Success")
}
