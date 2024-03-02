package routes

import (
	"fmt"
	"testing"
	"time"

	"github.com/AdityaP1502/Instant-Messanging/api/database"
	"github.com/AdityaP1502/Instant-Messanging/api/service/auth/config"
	"github.com/AdityaP1502/Instant-Messanging/api/service/auth/jwtutil"
	"github.com/jmoiron/sqlx"
)

const (
	host     = "localhost"
	port     = 3000
	user     = "instant"
	password = "4jBWgQ7qpmYq19+0y07Gc/VAts4QyBKrv1/UeORklQc="
	dbname   = "account_db"
)

func connectToDB(t *testing.T) *sqlx.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func TestInsertRevokedToken(t *testing.T) {
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
			SecretKeyRaw:      []byte("super super secret key"),
		},
	}

	db := connectToDB(t)

	claims := jwtutil.GenerateClaims(config, "lucas", "lucas12@gmail.com", jwtutil.User)
	token, err := jwtutil.GenerateToken(claims, config.Session.SecretKeyRaw)

	if err != nil {
		t.Error(err.Error())
		return
	}

	querynator := database.Querynator{}

	_, err = querynator.Insert(&RevokedToken{Token: token, Username: "lucas", ExpiredAt: claims.ExpiresAt.Local().Format(time.RFC3339), TokenType: string(claims.AccessType)}, db.DB, "revoked_token", "token_id")

	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Log("Success")
}

func TestSearchToken(t *testing.T) {
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
			SecretKeyRaw:      []byte("super super secret key"),
		},
	}

	db := connectToDB(t)

	claims := jwtutil.GenerateClaims(config, "lucas14", "lucas14@gmail.com", jwtutil.User)
	token, err := jwtutil.GenerateToken(claims, config.Session.SecretKeyRaw)

	if err != nil {
		t.Error(err.Error())
		return
	}

	querynator := database.Querynator{}

	id, err := querynator.Insert(&RevokedToken{Token: token, Username: "lucas14", ExpiredAt: claims.ExpiresAt.Local().Format(time.RFC3339), TokenType: string(claims.AccessType)}, db.DB, "revoked_token", "token_id")

	if err != nil {
		t.Error(err.Error())
		return
	}

	tokenData := &RevokedToken{}

	err = querynator.FindOne(&RevokedToken{TokenID: fmt.Sprintf("%d", id)}, tokenData, db.DB, "revoked_token", "expired_at", "token")

	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Log("Success")
}
