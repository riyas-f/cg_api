package routes

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/AdityaP1502/Instant-Messanging/api/database"
	httpx "github.com/AdityaP1502/Instant-Messanging/api/http"
	"github.com/AdityaP1502/Instant-Messanging/api/http/middleware"
	"github.com/AdityaP1502/Instant-Messanging/api/http/responseerror"
	"github.com/AdityaP1502/Instant-Messanging/api/jsonutil"
	"github.com/AdityaP1502/Instant-Messanging/api/service/call/config"
	"github.com/AdityaP1502/Instant-Messanging/api/service/call/payload"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

var querynator = &database.Querynator{}

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

var CALL_STATUS_TABLE_NAME = "user_call_status"

func updateCallStatussHandler(db *sql.DB, conf interface{}, w http.ResponseWriter, r *http.Request) responseerror.HTTPCustomError {
	body := r.Context().Value(middleware.PayloadKey).(*payload.CallStatus)

	_, ok := payload.ParseStringToStatus(body.Status)
	if !ok {
		return responseerror.CreateBadRequestError(
			responseerror.PayloadInvalid,
			responseerror.PayloadInvalidMessage,
			nil,
		)
	}

	dbX := sqlx.NewDb(db, "postgres")

	userStatus := &payload.CallStatus{}
	otherUserStatus := &payload.CallStatus{}

	err := querynator.FindOne(&payload.CallStatus{Username: body.Username}, userStatus, db, CALL_STATUS_TABLE_NAME, "call_status")
	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		return responseerror.CreateNotFoundError(map[string]string{
			"resourceName": "username",
		})
	default:
		return responseerror.CreateInternalServiceError(err)
	}

	// Forbidden condition
	if userStatus.Status == string(payload.Free) && body.Status == string(payload.OnCall) || userStatus.Status == string(payload.OnCall) && body.Status == string(payload.OnWait) {
		return responseerror.CreateBadRequestError(
			responseerror.CallsStatusInvalid,
			responseerror.StatusInvalidMessage,
			map[string]string{
				"reqStatus": body.Status,
				"status":    userStatus.Status,
			},
		)
	}

	err = querynator.FindOne(&payload.CallStatus{Username: body.OtherPartyUsername}, otherUserStatus, db, CALL_STATUS_TABLE_NAME, "call_status")
	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	// Forbidden condition
	if userStatus.Status == string(payload.Free) && body.Status == string(payload.OnCall) || userStatus.Status == string(payload.OnCall) && body.Status == string(payload.OnWait) {
		return responseerror.CreateBadRequestError(
			responseerror.CallsStatusInvalid,
			responseerror.StatusInvalidMessage,
			map[string]string{
				"reqStatus": body.Status,
				"status":    userStatus.Status,
			},
		)
	}

	tx, err := dbX.Beginx()

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	err = querynator.Update(&payload.CallStatus{Status: body.Status}, []string{"username"}, []any{body.Username}, tx, CALL_STATUS_TABLE_NAME)

	if err != nil {
		tx.Rollback()
		return responseerror.CreateInternalServiceError(err)
	}

	err = querynator.Update(&payload.CallStatus{Status: body.Status}, []string{"username"}, []any{body.OtherPartyUsername}, tx, CALL_STATUS_TABLE_NAME)

	if err != nil {
		tx.Rollback()
		return responseerror.CreateInternalServiceError(err)
	}

	json, err := jsonutil.EncodeToJson(
		&Response{
			Status:  "success",
			Message: "call status successfully updated",
		},
	)

	if err != nil {
		tx.Rollback()
		return responseerror.CreateInternalServiceError(err)
	}

	w.WriteHeader(200)
	_, err = w.Write(json)
	if err != nil {
		tx.Rollback()
		return responseerror.CreateInternalServiceError(err)
	}
	tx.Commit()

	return nil
}

func createNewUserEntry(db *sql.DB, conf interface{}, w http.ResponseWriter, r *http.Request) responseerror.HTTPCustomError {
	// only username is required here
	body := r.Context().Value(middleware.PayloadKey).(*payload.CallStatus)

	dbX := sqlx.NewDb(db, "postgres")
	tx, err := dbX.Beginx()

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	_, err = querynator.Insert(&payload.CallStatus{
		Status:   string(payload.Free),
		Username: body.Username,
	}, tx, CALL_STATUS_TABLE_NAME, "")

	if err != nil {
		tx.Rollback()
		return responseerror.CreateInternalServiceError(err)
	}

	json, err := jsonutil.EncodeToJson(
		&Response{
			Status:  "success",
			Message: "call status successfully created",
		},
	)

	if err != nil {
		tx.Rollback()
		return responseerror.CreateInternalServiceError(err)
	}

	w.WriteHeader(201)
	_, err = w.Write(json)

	if err != nil {
		tx.Rollback()
		return responseerror.CreateInternalServiceError(err)
	}

	tx.Commit()

	return nil
}

func SetCallRoute(r *mux.Router, db *sql.DB, conf *config.Config) {
	updateStatus := httpx.CreateHTTPHandler(db, conf, updateCallStatussHandler)

	updatePayloadGate, err := middleware.PayloadCheckMiddleware(
		&payload.CallStatus{},
		"Username",
		"Status",
		"OtherPartyUsername",
	)

	if err != nil {
		log.Fatal(err)
	}

	addUserEntry := httpx.CreateHTTPHandler(db, conf, createNewUserEntry)
	addEntryPayloadGate, err := middleware.PayloadCheckMiddleware(
		&payload.CallStatus{},
		"Username",
	)

	if err != nil {
		log.Fatal(err)
	}

	subrouter := r.PathPrefix("/call").Subrouter()

	subrouter.Handle("", middleware.UseMiddleware(db, conf, addUserEntry, addEntryPayloadGate)).Methods("POST")
	subrouter.Handle("", middleware.UseMiddleware(db, conf, updateStatus, updatePayloadGate)).Methods("PATCH")

}
