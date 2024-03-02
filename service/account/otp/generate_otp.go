package otp

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateOTP() (string, error) {
	// Generate user OTP
	n, err := rand.Int(rand.Reader, big.NewInt(900000))

	if err != nil {
		return "0", nil
	}

	otp := n.Add(n, big.NewInt(100000)) // produce a 6 digit otp

	return fmt.Sprintf("%d", otp.Int64()), nil
}
