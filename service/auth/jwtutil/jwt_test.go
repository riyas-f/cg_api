package jwtutil

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/AdityaP1502/Instant-Messanging/api/http/responseerror"
	"github.com/AdityaP1502/Instant-Messanging/api/service/auth/config"
)

var jwtKey = []byte("This is a super secret key")

func TestJWTGeneration(t *testing.T) {
	config := &config.Config{
		ServiceName: "test-app",
		Session: struct {
			ExpireTime        int    "json:\"expireTimeMinutes,string\""
			RefreshExpireTime int    "json:\"refreshExpireTimeMinutes,string\""
			SecretKeyBase64   string "json:\"secretKey\""
			SecretKeyRaw      []byte "json:\"-\""
		}{
			ExpireTime:        5,
			RefreshExpireTime: 60,
			SecretKeyRaw:      jwtKey,
		},
	}

	claims := GenerateClaims(config, "aditya", "adityanotgeh@email.com", User)
	token, err := GenerateToken(claims, config.Session.SecretKeyRaw)

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(token)
	t.Log("Success")
}

func TestValidJWTToken(t *testing.T) {
	config := &config.Config{
		ServiceName: "test-app",
		Session: struct {
			ExpireTime        int    "json:\"expireTimeMinutes,string\""
			RefreshExpireTime int    "json:\"refreshExpireTimeMinutes,string\""
			SecretKeyBase64   string "json:\"secretKey\""
			SecretKeyRaw      []byte "json:\"-\""
		}{
			ExpireTime:        5,
			RefreshExpireTime: 60,
			SecretKeyRaw:      jwtKey,
		},
	}

	claims := GenerateClaims(config, "aditya", "adityanotgeh@email.com", User)
	token, err := GenerateToken(claims, config.Session.SecretKeyRaw)

	if err != nil {
		t.Error(err)
		return
	}

	decodedClaims, err := VerifyToken(token, config.Session.SecretKeyRaw)

	if err != nil {
		t.Error(err)
		return
	}

	if decodedClaims.Username != claims.Username {
		t.Errorf("Claim username not match. expected %s received %s", claims.Username, decodedClaims.Username)
		return
	}

	t.Log(decodedClaims)
	t.Log(decodedClaims.ExpiresAt.Format(time.RFC3339))
	t.Log("Success")
}

func TestInvalidToken(t *testing.T) {
	token := "xxxxxxxddddddd"
	config := &config.Config{
		ServiceName: "test-app",
		Session: struct {
			ExpireTime        int    "json:\"expireTimeMinutes,string\""
			RefreshExpireTime int    "json:\"refreshExpireTimeMinutes,string\""
			SecretKeyBase64   string "json:\"secretKey\""
			SecretKeyRaw      []byte "json:\"-\""
		}{
			ExpireTime:        5,
			RefreshExpireTime: 60,
			SecretKeyRaw:      jwtKey,
		},
	}

	_, err := VerifyToken(token, config.Session.SecretKeyRaw)

	if err.Get().Name != string(responseerror.InvalidToken) {
		t.Errorf("Wrong error type found")
		t.Error(err)
		return
	}

	t.Log("Success")
}

// func TestExpiredToken(t *testing.T) {
// 	config := &config.Config{
// 		ServiceName: "test-app",
// 		Session: struct {
// 			ExpireTime        int    "json:\"expireTimeMinutes,string\""
// 			RefreshExpireTime int    "json:\"refreshExpireTimeMinutes,string\""
// 			SecretKeyBase64   string "json:\"secretKey\""
// 			SecretKeyRaw      []byte "json:\"-\""
// 		}{
// 			ExpireTime:        1,
// 			RefreshExpireTime: 60,
// 			SecretKeyRaw:      jwtKey,
// 		},
// 	}

// 	claims := GenerateClaims(config, "aditya", "adityanotgeh@email.com", User)

// 	json, err := jsonutil.EncodeToJson(claims)

// 	fmt.Println(string(json))

// 	t.Logf("Token expired at: %s", claims.RegisteredClaims.ExpiresAt.Format(time.RFC3339))
// 	token, err := GenerateToken(claims, config.Session.SecretKeyRaw)

// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	time.Sleep(time.Duration(2) * time.Minute)

// 	_, err = VerifyToken(token, config.Session.SecretKeyRaw)

// 	if err == nil {
// 		t.Errorf("token is supposedly expired, not valid!")
// 		return
// 	}

// 	if err.(responseerror.HTTPCustomError).Get().Name != string(responseerror.TokenExpired) {
// 		t.Errorf("Wrong error type")
// 		t.Error(err)
// 		return
// 	}

// 	t.Log("Success")
// }

func TestToken(t *testing.T) {
	base64Key := "yTTeQ6qTkxWgr3oeDa53FuPpoen9hINxF+zyAZIxJb6WyrdPuZvRMPA6S0wVgl4mPuMojOpll5EcXXYWDax7q6oqe12WV/iN1G7OVouKuOy6R0OuMVylC5y+fxdI4kYUS/L0AIS0aVeQKDyWUmVHbRXjO58v3GSjnuw/hJeuqUQ="
	raw, err := base64.StdEncoding.DecodeString(base64Key)

	if err != nil {
		t.Error(err)
		return
	}

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkaXQxeTlhMTIxMDE1MSIsImVtYWlsIjoiaS5tLjFhQDExZzJtYWlsLmNvbSIsInJvbGVzIjoidXNlciIsImlzcyI6Ikluc3RhbnQtTWVzc2FnaW5nIEFQSSIsImV4cCI6MTcwODc5NzAxN30.fFKgOJOjm5ZCTPW5kV-xAyxaAHpBtZqmwdX-xvSxKYM"
	claims, err := VerifyToken(token, raw)

	t.Log(claims)
}
