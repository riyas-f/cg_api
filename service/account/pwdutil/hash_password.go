package pwdutil

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

func fillMaximumRandomNumber(totalByte int) (*big.Int, *big.Int) {
	a := new(big.Int)
	b := new(big.Int)

	min := "1"
	max := "e"

	for i := 0; i < 2*totalByte-1; i++ {
		min += "0"
		max += "f"
	}

	a, ok := a.SetString(min, 16)

	if !ok {
		panic(errors.New("cannot create big number for rng"))
	}

	b, ok = b.SetString(max, 16)
	if !ok {
		panic(errors.New("cannot create big number for rng"))
	}

	return a, b

}

var RANDOM_NUMBER_MAX, RANDOM_NUMBER_MIN = fillMaximumRandomNumber(16)

func generateSalt() (string, error) {

	r, err := rand.Int(rand.Reader, RANDOM_NUMBER_MAX)

	if err != nil {
		return "", err
	}

	r = r.Add(RANDOM_NUMBER_MIN, r)
	salt := hex.EncodeToString(r.Bytes())

	return salt, nil
}

func HashPassword(password string, secretKey []byte) (string, string, error) {
	salt, err := generateSalt()

	if err != nil {
		return "", "", err
	}

	// prehashing
	h := hmac.New(sha512.New384, secretKey)
	h.Write([]byte(password + salt))

	prehashPW := h.Sum(nil)

	hash, err := bcrypt.GenerateFromPassword(prehashPW, bcrypt.DefaultCost)

	if err != nil {
		return "", "", err
	}

	return base64.StdEncoding.EncodeToString(hash), salt, nil
}

func CheckPassword(password string, salt string, base64Hash string, secretKey []byte) (bool, error) {
	h := hmac.New(sha512.New384, secretKey)
	h.Write([]byte(password + salt))
	prehashPW := h.Sum(nil)

	passwordHashByte, err := base64.StdEncoding.DecodeString(base64Hash)

	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword(passwordHashByte, prehashPW)
	return err == nil, nil
}
