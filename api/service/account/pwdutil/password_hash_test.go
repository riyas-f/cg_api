package pwdutil

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "SUPER_SECURE_PASSWORD"

	hashedPassword, salt, err := HashPassword(password, []byte("Super_Secret_Key"))

	if err != nil {
		t.Error(err)
		return
	}

	isMatch, err := CheckPassword(password, salt, hashedPassword, []byte("Super_Secret_Key"))

	if err != nil {
		t.Error(err)
		return
	}

	if !isMatch {
		t.Errorf("Password should match")
		return
	}

	t.Logf("Success: %s", hashedPassword)
	t.Logf("Password has length is %d", len(hashedPassword))
}
