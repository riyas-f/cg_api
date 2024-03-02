package routes

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/AdityaP1502/Instant-Messanging/api/database"
	"github.com/AdityaP1502/Instant-Messanging/api/date"
	httpx "github.com/AdityaP1502/Instant-Messanging/api/http"
	"github.com/AdityaP1502/Instant-Messanging/api/http/middleware"
	"github.com/AdityaP1502/Instant-Messanging/api/http/responseerror"
	"github.com/AdityaP1502/Instant-Messanging/api/jsonutil"
	"github.com/AdityaP1502/Instant-Messanging/api/service/account/config"
	accountMiddleware "github.com/AdityaP1502/Instant-Messanging/api/service/account/middleware"
	"github.com/AdityaP1502/Instant-Messanging/api/service/account/otp"
	"github.com/AdityaP1502/Instant-Messanging/api/service/account/payload"
	"github.com/AdityaP1502/Instant-Messanging/api/service/account/pwdutil"
	"github.com/AdityaP1502/Instant-Messanging/api/service/auth/jwtutil"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

var querynator = &database.Querynator{}

type RegisterResponse struct {
	Status       string `json:"status"`
	OTPConfirmID string `json:"otp_confirmation_id"`
	AccessToken  string `json:"access_token"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginResponse struct {
	Status string `json:"status"`
	Token  Token  `json:"token"`
}

type GenericResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var AUTH_ISSUE_TOKEN_ENDPOINT string = "v1/auth/token/issue"
var SEND_MAIL_ENDPOINT string = "mail/send"
var AUTH_REVOKE_TOKEN_ENDPOINT string = "v1/auth/token/revoke"

// func sendMailHTTP(message string, subject string, to string, url string) error {
// 	//TODO: Send http request to node js server

// 	var client = &http.Client{}

// 	var mail struct {
// 		To      string `json:"to"`
// 		Subject string `json:"subject"`
// 		Message string `json:"message"`
// 	}

// 	mail.To = to
// 	mail.Subject = subject
// 	mail.Message = message

// 	json, err := jsonutil.CreateJSONResponse(&mail)

// 	if err != nil {
// 		return err
// 	}

// 	r, err := http.NewRequest("POST", url, bytes.NewReader(json))

// 	if err != nil {
// 		return err
// 	}

// 	r.Header.Set("Content-Type", "application/json")

// 	resp, err := client.Do(r)

// 	if err != nil {
// 		return err
// 	}

// 	defer resp.Body.Close()

// 	if resp.StatusCode != 200 {
// 		return fmt.Errorf("failed to send mail. mail server return with status code %d", resp.StatusCode)
// 	}

// 	return nil
// }

func registerHandler(db *sql.DB, conf interface{}, w http.ResponseWriter, r *http.Request) responseerror.HTTPCustomError {
	var req *httpx.HTTPRequest

	cf := conf.(*config.Config)

	body := r.Context().Value(middleware.PayloadKey).(*payload.Account)

	// Check username and email exist or not
	exist, err := querynator.IsExists(&payload.Account{Email: body.Email}, db, "account")

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	if exist {
		return responseerror.CreateBadRequestError(
			responseerror.EmailExists,
			responseerror.EmailsExistMessage,
			nil,
		)
	}

	exist, err = querynator.IsExists(&payload.Account{Username: body.Username}, db, "account")
	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	if exist {
		return responseerror.CreateBadRequestError(
			responseerror.UsernameExists,
			responseerror.UsernameExistMessage,
			nil,
		)
	}

	newUser, err := payload.NewRegisteredAccountPayload(body.Username, body.Name, body.Email, body.Password, cf.Hash.SecretKeyRaw)

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	tx, err := sqlx.NewDb(db, cf.Database.Driver).Beginx()

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	_, err = querynator.Insert(newUser, tx, "account", "account_id")
	if err != nil {
		tx.Rollback()
		return responseerror.CreateInternalServiceError(err)
	}

	otpData, err := payload.NewOTPPayload(body.Username)

	if err != nil {
		tx.Rollback()
		return responseerror.CreateInternalServiceError(err)
	}

	_, err = querynator.Insert(otpData, tx, "user_otp", "otp_id")

	if err != nil {
		tx.Rollback()
		return responseerror.CreateInternalServiceError(err)
	}

	//TODO: Get access token from auth endpoint
	req = &httpx.HTTPRequest{}
	req, err = req.CreateRequest(
		cf.Services.Auth.Scheme,
		cf.Services.Auth.Host,
		cf.Services.Auth.Port,
		AUTH_ISSUE_TOKEN_ENDPOINT,
		http.MethodPost,
		http.StatusOK,
		struct {
			Username  string `json:"username"`
			Email     string `json:"email"`
			Roles     string `json:"roles"`
			TokenType string `json:"token_type"`
		}{
			Username:  body.Username,
			Email:     body.Email,
			Roles:     "user",
			TokenType: "access",
		},
		cf.Config,
	)

	if err != nil {
		tx.Rollback()
		return err.(responseerror.HTTPCustomError)
	}

	token := &Token{}
	err = req.Send(token)

	if err != nil {
		tx.Rollback()
		return err.(responseerror.HTTPCustomError)
	}

	// Create an API Call to mail service
	req = &httpx.HTTPRequest{}
	req, err = req.CreateRequest(
		cf.Services.Mail.Scheme,
		cf.Services.Mail.Host,
		cf.Services.Mail.Port,
		SEND_MAIL_ENDPOINT,
		http.MethodPost,
		http.StatusOK,
		struct {
			To      string `json:"to"`
			Subject string `json:"subject"`
			Message string `json:"message"`
		}{
			To:      body.Email,
			Subject: "Email Verification",
			Message: fmt.Sprintf("Dont share this with anyone. This is your OTP %s. Your token will expired in %d minutes",
				otpData.OTP, cf.OTP.OTPDurationMinutes),
		},
		cf.Config,
	)

	if err != nil {
		tx.Rollback()
		return err.(responseerror.HTTPCustomError)
	}

	err = req.Send(nil)
	if err != nil {
		tx.Rollback()
		return err.(responseerror.HTTPCustomError)
	}

	resp := &RegisterResponse{
		Status:       "success",
		OTPConfirmID: otpData.OTPConfirmID,
		AccessToken:  token.AccessToken,
	}

	jsonResponse, err := jsonutil.EncodeToJson(resp)

	if err != nil {
		tx.Rollback()
		return responseerror.CreateInternalServiceError(err)
	}

	tx.Commit()

	w.WriteHeader(200)
	w.Write(jsonResponse)

	return nil
}

func loginHandler(db *sql.DB, conf interface{}, w http.ResponseWriter, r *http.Request) responseerror.HTTPCustomError {
	cf := conf.(*config.Config)

	body := r.Context().Value(middleware.PayloadKey).(*payload.Account)

	// Grab password and salt from db associated with username
	user := &payload.Account{}
	err := querynator.FindOne(&payload.Account{Email: body.Email}, user, db, "account", "username", "password", "password_salt", "is_active")

	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		return responseerror.CreateNotFoundError(map[string]string{"resourceName": "email"})
	default:
		return responseerror.CreateInternalServiceError(err)
	}

	if isMatch, err := pwdutil.CheckPassword(body.Password, user.Salt, user.Password, cf.Hash.SecretKeyRaw); err != nil {
		return responseerror.CreateInternalServiceError(err)
	} else if !isMatch {
		return responseerror.CreateUnauthenticatedError(
			responseerror.InvalidCredentials,
			responseerror.InvalidCredentialsMessage,
			nil,
		)
	} else if user.IsActive == strconv.FormatBool(false) {
		return responseerror.CreateUnauthenticatedError(
			responseerror.UserMarkedInActive,
			responseerror.UserMarkedInActiveMessage,
			nil,
		)
	}

	// user is good and dandy
	req := &httpx.HTTPRequest{}
	req, err = req.CreateRequest(
		cf.Services.Auth.Scheme,
		cf.Services.Auth.Host,
		cf.Services.Auth.Port,
		AUTH_ISSUE_TOKEN_ENDPOINT,
		http.MethodPost,
		http.StatusOK,
		struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Roles    string `json:"roles"`
		}{
			Username: user.Username,
			Email:    body.Email,
			Roles:    "user",
		},
		cf.Config,
	)

	if err != nil {
		return err.(responseerror.HTTPCustomError)
	}

	token := Token{}
	err = req.Send(&token)

	if err != nil {
		return err.(responseerror.HTTPCustomError)
	}

	json, err := jsonutil.EncodeToJson(&LoginResponse{Status: "success", Token: token})

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	w.WriteHeader(200)
	w.Write(json)

	return nil
}

func resendOTPHandler(db *sql.DB, conf interface{}, w http.ResponseWriter, r *http.Request) responseerror.HTTPCustomError {
	// TODO: Use Transaction when inserting data or update data into the database
	cf := conf.(*config.Config)

	vars := mux.Vars(r)
	confirmID := vars["otp_confirmation_id"]
	u := &payload.UserOTP{OTPConfirmID: confirmID}

	claims := r.Context().Value(accountMiddleware.ClaimsKey).(*jwtutil.Claims)
	token := r.Context().Value(accountMiddleware.TokenKey).(string)

	// check if confirmation id exists
	err := querynator.FindOne(&payload.UserOTP{OTPConfirmID: confirmID, Username: claims.Username, MarkedForDeletion: strconv.FormatBool(false)}, u, db, "user_otp",
		"otp_id",
		"last_resend",
	)

	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		return responseerror.CreateNotFoundError(map[string]string{
			"resourceName": "email",
		})
	default:
		return responseerror.CreateInternalServiceError(err)
	}

	// Check last resend duration
	t, err := date.ParseTimestamp(u.LastResend)

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	duration := date.MinutesDifferenceFronNow(t)

	if duration < cf.OTP.ResendDurationMinutes {
		return responseerror.CreateTooManyRequestError(
			responseerror.ResendIntervalNotReachedErr,
			responseerror.ResendIntervalNotReachedMessage,
			nil,
		)
	}

	// revoked user token
	req := &httpx.HTTPRequest{}
	req, err = req.CreateRequest(
		cf.Services.Auth.Scheme,
		cf.Services.Auth.Host,
		cf.Services.Auth.Port,
		AUTH_REVOKE_TOKEN_ENDPOINT,
		http.MethodPost,
		http.StatusOK,
		&Token{
			AccessToken: token,
		},
		cf.Config,
	)

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	err = req.Send(nil)

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	newToken := &Token{}

	req = &httpx.HTTPRequest{}
	req, err = req.CreateRequest(
		cf.Services.Auth.Scheme,
		cf.Services.Auth.Host,
		cf.Services.Auth.Port,
		AUTH_ISSUE_TOKEN_ENDPOINT,
		http.MethodPost,
		http.StatusOK,
		struct {
			Username  string `json:"username"`
			Email     string `json:"email"`
			Roles     string `json:"roles"`
			TokenType string `json:"token_type"`
		}{
			Username:  claims.Username,
			Email:     claims.Email,
			Roles:     "user",
			TokenType: "access",
		},
		cf.Config,
	)

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	err = req.Send(newToken)
	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	otp, err := otp.GenerateOTP()

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	req = &httpx.HTTPRequest{}
	req, err = req.CreateRequest(
		cf.Services.Mail.Scheme,
		cf.Services.Mail.Host,
		cf.Services.Mail.Port,
		SEND_MAIL_ENDPOINT,
		http.MethodPost,
		http.StatusOK,
		struct {
			To      string `json:"to"`
			Subject string `json:"subject"`
			Message string `json:"message"`
		}{
			To:      claims.Email,
			Subject: "Email Verification",
			Message: fmt.Sprintf("Dont share this with anyone. This is your OTP %s. Your token will expired in %d minutes",
				otp, cf.OTP.OTPDurationMinutes),
		},
		cf.Config,
	)

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	err = req.Send(nil)

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	json, err := jsonutil.EncodeToJson(struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Token   string `json:"token"`
	}{Status: "success", Message: "OTP has been re-send to your email.", Token: newToken.AccessToken})

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	w.WriteHeader(200)
	w.Write(json)

	// Update the otp
	err = querynator.Update(&payload.UserOTP{OTP: otp, LastResend: date.GenerateTimestamp()}, []string{"otp_id"}, []any{u.OTPID}, db, "user_otp")

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	return nil
}

func verifyOTPHandler(db *sql.DB, conf interface{}, w http.ResponseWriter, r *http.Request) responseerror.HTTPCustomError {
	var validOTP = &payload.UserOTP{}

	cf := conf.(*config.Config)

	body := r.Context().Value(middleware.PayloadKey).(*payload.UserOTP)
	claims := r.Context().Value(accountMiddleware.ClaimsKey).(*jwtutil.Claims)
	token := r.Context().Value(accountMiddleware.TokenKey).(string)

	// Fill the username in the payload
	body.Username = claims.Username

	err := querynator.FindOne(&payload.UserOTP{
		Username: body.Username, OTPConfirmID: body.OTPConfirmID,
		MarkedForDeletion: strconv.FormatBool(false)},
		validOTP, db, "user_otp", "otp", "otp_id",
	)

	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		return responseerror.CreateNotFoundError(map[string]string{"resourceName": "email"})
	default:
		return responseerror.CreateInternalServiceError(err)
	}

	if validOTP.OTP != body.OTP {
		// otp is wrong
		return responseerror.CreateBadRequestError(
			responseerror.OTPInvalid,
			responseerror.OTPInvalidMessage,
			nil,
		)
	}

	json, err := jsonutil.EncodeToJson(&GenericResponse{Status: "success", Message: "your account has been activated successfully"})

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	// otp is correct, update user to be an active user, marked otp entry, and add token to revoked list

	// create sqlx connection
	sqlxDb := sqlx.NewDb(db, "postgres")
	tx, err := sqlxDb.Beginx()

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	err = querynator.Update(&payload.Account{IsActive: strconv.FormatBool(true)}, []string{"username"}, []any{claims.Username}, tx, "account")
	if err != nil {
		rollError := tx.Rollback()
		fmt.Println(rollError.Error())
		return responseerror.CreateInternalServiceError(err)
	}

	err = querynator.Update(&payload.UserOTP{MarkedForDeletion: strconv.FormatBool(true)}, []string{"otp_id"}, []any{validOTP.OTPID}, tx, "user_otp")

	if err != nil {
		rollError := tx.Rollback()
		fmt.Println(rollError.Error())
		return responseerror.CreateInternalServiceError(err)
	}

	req := &httpx.HTTPRequest{}
	req, err = req.CreateRequest(
		cf.Services.Auth.Scheme,
		cf.Services.Auth.Host,
		cf.Services.Auth.Port,
		AUTH_REVOKE_TOKEN_ENDPOINT,
		http.MethodPost,
		http.StatusOK,
		&Token{
			AccessToken: token,
		},
		cf.Config,
	)

	if err != nil {
		tx.Rollback()
		return responseerror.CreateInternalServiceError(err)
	}

	err = req.Send(nil)

	if err != nil {
		tx.Rollback()
		return responseerror.CreateInternalServiceError(err)
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return responseerror.CreateInternalServiceError(err)
	}

	w.WriteHeader(200)
	w.Write(json)

	return nil
}

// Register account subrouter
func SetAccountRoute(r *mux.Router, db *sql.DB, config *config.Config) {
	subrouter := r.PathPrefix("/account").Subrouter()

	subrouter.Use(middleware.RouteGetterMiddleware)

	// Create middleware here
	userPayloadMiddleware, err := middleware.PayloadCheckMiddleware(&payload.Account{}, "Username", "Name", "Email", "Password")

	if err != nil {
		log.Fatal(err)
	}

	otpPayloadMiddleware, err := middleware.PayloadCheckMiddleware(&payload.UserOTP{}, "OTPConfirmID", "OTP")

	if err != nil {
		log.Fatal(err)
	}

	loginPayloadMIddleware, err := middleware.PayloadCheckMiddleware(&payload.Account{}, "Email", "Password")

	if err != nil {
		log.Fatal(err)
	}

	basicAccessAuthMiddleware := accountMiddleware.AuthMiddleware

	// REGISTER ROUTE //
	register := &httpx.Handler{
		DB:      db,
		Config:  config,
		Handler: registerHandler,
	}

	subrouter.Handle("/register", middleware.UseMiddleware(db, config, register, userPayloadMiddleware)).Methods("POST")

	// VERIFY OTP ROUTE //
	verifyOTP := &httpx.Handler{
		DB:      db,
		Config:  config,
		Handler: verifyOTPHandler,
	}

	subrouter.Handle("/otp/verify", middleware.UseMiddleware(db, config, verifyOTP, basicAccessAuthMiddleware, otpPayloadMiddleware)).Methods("POST")

	// RESEND OTP ROUTE //
	resendOTP := &httpx.Handler{
		DB:      db,
		Config:  config,
		Handler: resendOTPHandler,
	}

	subrouter.Handle("/otp/{otp_confirmation_id}/resend", middleware.UseMiddleware(db, config, resendOTP, basicAccessAuthMiddleware)).Methods("POST")

	// LOGIN ROUTE //
	login := &httpx.Handler{
		DB:      db,
		Config:  config,
		Handler: loginHandler,
	}

	subrouter.Handle("/login", middleware.UseMiddleware(db, config, login, loginPayloadMIddleware)).Methods("POST")

	// subrouter.HandleFunc("/logout", logOutHandler).Methods("POST")
	// subrouter.HandleFunc("/{username}", patchUserInfoHandler).Methods("PATCH")
}
