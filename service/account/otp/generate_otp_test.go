package otp

import (
	"strconv"
	"testing"
)

func TestGenerateOTP(t *testing.T) {
	otp, err := GenerateOTP()

	if err != nil {
		t.Error(err)
	}

	otpI, err := strconv.Atoi(otp)

	if err != nil {
		t.Error(err)
	}

	if otpI > 999999 || otpI < 100000 {
		t.Errorf("OTP %s isn't a 6 digit number", otp)
		return
	}

	t.Logf("Success. %s is a valid 6 digit number", otp)
}
