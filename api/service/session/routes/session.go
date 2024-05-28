package routes

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/AdityaP1502/Instant-Messanging/api/database"
	"github.com/AdityaP1502/Instant-Messanging/api/date"
	httpx "github.com/AdityaP1502/Instant-Messanging/api/http"
	"github.com/AdityaP1502/Instant-Messanging/api/http/middleware"
	"github.com/AdityaP1502/Instant-Messanging/api/http/responseerror"
	"github.com/AdityaP1502/Instant-Messanging/api/jsonutil"
	"github.com/AdityaP1502/Instant-Messanging/api/service/auth/jwtutil"
	"github.com/AdityaP1502/Instant-Messanging/api/service/session/config"
	"github.com/AdityaP1502/Instant-Messanging/api/service/session/payload"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	MAX_RETRIES = 3
)

const (
	SESSION_MANAGER_CREATE_ENDPOINT = "v1/vms"
	SESSION_MANAGER_DELETE_ENDPOINT = "v1/vms/%s"
	SESSION_MANAGER_CHECK_TEMPLATES = "v1/vms/templates"
	VM_PIN_ENDPOINT                 = "sendpin"
)

const (
	PROVISIONING           string = "Provisioning"
	WAITING_FOR_CONNECTION string = "WaitingForConnection"
	PAIRING                string = "Pairing"
	RUNNING                string = "Running"
	FAILED                 string = "Failed"
	TERMINATED             string = "Terminated"
)

var querynator = database.Querynator{
	DriverName: "postgres",
}

// // Use Locking
// func deattachGPUFromUsers(body *payload.UserSession, db *sql.DB) responseerror.HTTPCustomError {
// 	return nil
// }

// TODO: Improve Querynator for non string data types
func attachGPUToUsers(body *payload.UserSession, db *sql.DB, retry int) (*sql.Tx, responseerror.HTTPCustomError) {
	var dest payload.GPU

	if retry >= MAX_RETRIES {
		return nil, responseerror.CreateConflictError(
			responseerror.UpdateConflictErr,
			responseerror.UpdateConflictErrorMessage,
			nil,
		)
	}

	// Check if the gpu is available
	err := querynator.FindOne(&payload.GPU{GPUName: body.SessionMetadata.GPUName}, &dest, db, "gpu_list", "n_available", "version", "gpu_alt_name", "gpu_id")

	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		return nil, responseerror.CreateNotFoundError(map[string]string{"resourceName": "gpu_name"})
	default:
		return nil, responseerror.CreateInternalServiceError(err)
	}

	// Optimistic cast
	n, _ := strconv.Atoi(dest.Count)
	v, _ := strconv.Atoi(dest.Version)

	if n < 1 {
		return nil, responseerror.CreateBadRequestError(responseerror.GPUNotAvailable,
			responseerror.GPUNotAvailableMessage, map[string]string{
				"gpuName": body.SessionMetadata.GPUName,
			})
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, responseerror.CreateInternalServiceError(err)
	}

	// Update the gpu entry
	result, err := querynator.UpdateWithResults(
		&payload.GPU{
			Version: fmt.Sprintf("%d", v+1),
			Count:   fmt.Sprintf("%d", n-1),
		},
		[]string{"gpu_id"},
		[]any{dest.GPUID},
		tx,
		"gpu_list",
	)

	if err != nil {
		return nil, responseerror.CreateInternalServiceError(err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return nil, responseerror.CreateInternalServiceError(err)
	}

	// Recursive retry
	if rowsAffected == 0 {
		return attachGPUToUsers(body, db, retry+1)
	}

	tx.Commit()

	return tx, nil
}

func getGPUStatusHandler(db *sql.DB, conf interface{}, w http.ResponseWriter, r *http.Request) responseerror.HTTPCustomError {
	available := strings.ToLower(r.URL.Query().Get("only_available"))
	limit := 0

	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil {
		limit = l
	}

	min_count := 0

	if available == "true" {
		min_count = 1
	}

	var dest []payload.GPU

	err := querynator.FindWithCondition(
		[]database.QueryCondition{{TableName: "gpu_list", ColumnName: "n_available", MatchValue: min_count, Operand: database.GEQ}},
		&dest,
		limit,
		db,
		"gpu_list",
		"gpu_name", "n_available", "version",
	)

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	response := struct {
		Status string        `json:"status"`
		GPUs   []payload.GPU `json:"gpu"`
	}{
		Status: "success",
		GPUs:   dest,
	}

	json, err := jsonutil.EncodeToJson(&response)

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	w.Write(json)

	return nil
}
func createNewSessionHandler(db *sql.DB, conf interface{}, w http.ResponseWriter, r *http.Request) responseerror.HTTPCustomError {
	cf := conf.(*config.Config)
	body := r.Context().Value(middleware.PayloadKey).(*payload.UserSession)

	uuidv7, err := uuid.NewV7()

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	sessionId, err := uuidv7.MarshalBinary()
	sessionIdString := uuidv7.String()

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	// Set additional information
	body.SessionID = sessionId
	body.LastUpdatedAt = date.GenerateTimestamp()
	body.MarkedForDeletion = strconv.FormatBool(false)
	body.RequestStatus = PROVISIONING

	body.SessionMetadata.SessionID = sessionId
	body.SessionMetadata.GameLocation.Server.Host = fmt.Sprintf("%s:%d",
		body.SessionMetadata.GameLocation.Server.IP,
		body.SessionMetadata.GameLocation.Server.Port,
	)
	body.SessionMetadata.CreatedAt = date.GenerateTimestamp()

	// This need improvement
	gpuTX, err := attachGPUToUsers(body, db, 0)

	if err != nil {
		return err.(responseerror.HTTPCustomError)
	}

	tx, err := db.Begin()

	if err != nil {
		gpuTX.Rollback()
		return responseerror.CreateInternalServiceError(err)
	}

	// only insert the root file on the body
	_, err = querynator.Insert(body, tx, "user_session", "")

	if err != nil {
		gpuTX.Rollback()
		tx.Rollback()
		return responseerror.CreateInternalServiceError(err)
	}

	// Insert the metadata into the database
	_, err = querynator.Insert(body.SessionMetadata, tx, "session_metadata", "")
	if err != nil {
		gpuTX.Rollback()
		tx.Rollback()
		return responseerror.CreateInternalServiceError(err)
	}

	req := &httpx.HTTPRequest{}
	req, err_ := req.CreateRequest(
		"http",
		cf.Service.SessionManager.Host,
		cf.Service.SessionManager.Port,
		SESSION_MANAGER_CHECK_TEMPLATES,
		http.MethodGet,
		200,
		nil,
		cf.Config,
	)

	if err_ != nil {
		gpuTX.Rollback()
		tx.Rollback()
		return responseerror.CreateInternalServiceError(err_)
	}

	err_ = req.Send(nil)

	// Propagate the error to the user
	if err_ != nil {
		gpuTX.Rollback()
		tx.Rollback()
		if _, ok := err_.(*responseerror.InternalServiceError); ok {
			return err_
		}

		w.WriteHeader(err_.Get().Code)
		w.Write(req.Payload)
		return nil
	}

	var sessionRequest struct {
		Name        string `json:"name"`
		SessionID   string `json:"SID"`
		Description string `json:"desc"`
		//PCIDevice   string `json:"pci_device"`
	}

	sessionRequest.Name = body.Username
	sessionRequest.SessionID = sessionIdString
	sessionRequest.Description = "VM Request"
	//sessionRequest.PCIDevice = ""

	req = &httpx.HTTPRequest{}
	req, err_ = req.CreateRequest(
		"http",
		cf.Service.SessionManager.Host,
		cf.Service.SessionManager.Port,
		SESSION_MANAGER_CREATE_ENDPOINT,
		http.MethodPost,
		200,
		sessionRequest,
		cf.Config,
	)

	if err_ != nil {
		gpuTX.Rollback()
		tx.Rollback()
		return responseerror.CreateInternalServiceError(err_)
	}

	err_ = req.Send(nil)

	// Propagate the error to the user
	if err_ != nil {
		gpuTX.Rollback()
		tx.Rollback()
		if _, ok := err_.(*responseerror.InternalServiceError); ok {
			return err_
		}

		w.WriteHeader(err_.Get().Code)
		w.Write(req.Payload)
		return nil
	}

	tmp := struct {
		Status    string `json:"status"`
		SessionID string `json:"session_id"`
	}{
		Status:    "success",
		SessionID: sessionIdString,
	}

	json, err := jsonutil.EncodeToJson(tmp)

	if err != nil {
		gpuTX.Rollback()
		tx.Rollback()
		return responseerror.CreateInternalServiceError(err)
	}

	w.Write(json)
	gpuTX.Commit()
	tx.Commit()

	return nil

}

func getRequestStatus(db *sql.DB, _ interface{}, w http.ResponseWriter, r *http.Request) responseerror.HTTPCustomError {
	vars := mux.Vars(r)

	claims := r.Context().Value(middleware.ClaimsKey).(*jwtutil.Claims)

	sessionIDstring := vars["session_id"]
	uuidv7, err := uuid.Parse(sessionIDstring)

	if err != nil {
		return responseerror.CreateBadRequestError(responseerror.MalformedSessionID, responseerror.MalformedSessionIDMessage,
			map[string]string{
				"id": "session_id",
			},
		)
	}

	sessionID, err := uuidv7.MarshalBinary()
	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	// err = querynator.FindOne(&payload.UserSession{
	// 	SessionID:         sessionID,
	// 	MarkedForDeletion: strconv.FormatBool(false),
	// }, userSession, db, "user_session", "username", "request_status", "last_updated")

	joinTables := []struct {
		payload.UserSession
		payload.SessionHost
	}{}

	joinExecutor := querynator.PrepareJoinOperation()
	joinExecutor.AddJoinTable("session_host", "session_id", "user_session", "session_id")

	err = joinExecutor.Find(db, []database.QueryCondition{
		{TableName: "user_session", ColumnName: "session_id", MatchValue: sessionID, Operand: database.EQ},
	}, &joinTables, "user_session", database.LEFT_JOIN, map[string][]string{
		"user_session": {"username", "request_status", "last_update"},
		"session_host": {"network_id,string"},
	})

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	if len(joinTables) < 1 {
		return responseerror.CreateNotFoundError(map[string]string{
			"resourceName": "session_id",
		})
	}

	if claims.Username != joinTables[0].Username {
		return responseerror.CreateUnauthorizedError(responseerror.AccessDenied, responseerror.AccessDeniedMesage, nil)
	}

	tmp := struct {
		Status        string `json:"status"`
		Username      string `json:"username"`
		LastUpdatedAt string `json:"last_updated"`
		RequestStatus string `json:"request_status"`
		NetworkID     string `json:"network_id"`
	}{
		Status:        "success",
		Username:      joinTables[0].Username,
		LastUpdatedAt: joinTables[0].LastUpdatedAt,
		RequestStatus: joinTables[0].RequestStatus,
		NetworkID:     joinTables[0].SessionHost.NetworkID,
	}

	json, err := jsonutil.EncodeToJson(tmp)

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	w.Write(json)

	return nil
}

func startConnectionEstablishmentHandler(db *sql.DB, _ interface{}, w http.ResponseWriter, r *http.Request) responseerror.HTTPCustomError {
	vars := mux.Vars(r)
	sessionIDString := vars["session_id"]

	uuidv7, err := uuid.Parse(sessionIDString)

	if err != nil {
		return responseerror.CreateBadRequestError(responseerror.MalformedSessionID, responseerror.MalformedSessionIDMessage, map[string]string{
			"id": "session_id",
		})
	}

	sessionID, err := uuidv7.MarshalBinary()
	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	body := r.Context().Value(middleware.PayloadKey).(*payload.SessionHost)
	body.SessionID_ = sessionID

	tx, err := db.Begin()

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	_, err = querynator.Insert(body, tx, "session_host", "")
	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	err = querynator.Update(&payload.UserSession{RequestStatus: WAITING_FOR_CONNECTION}, []string{"session_id"}, []any{sessionID}, tx, "user_session")
	if err != nil {
		tx.Rollback()
		return responseerror.CreateInternalServiceError(err)
	}

	tmp := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  "success",
		Message: "instance has been attached successfully",
	}

	json, err := jsonutil.EncodeToJson(tmp)

	if err != nil {
		tx.Rollback()
		return responseerror.CreateInternalServiceError(err)
	}

	w.Write(json)
	tx.Commit()
	return nil
}

func pairHandler(db *sql.DB, conf interface{}, w http.ResponseWriter, r *http.Request) responseerror.HTTPCustomError {
	cf := conf.(*config.Config)
	vars := mux.Vars(r)
	sessionIDString := vars["session_id"]

	uuidv7, err := uuid.Parse(sessionIDString)
	if err != nil {
		return responseerror.CreateBadRequestError(
			responseerror.MalformedSessionID,
			responseerror.MalformedSessionIDMessage,
			map[string]string{
				"id": "session_id",
			},
		)
	}

	sessionID, err := uuidv7.MarshalBinary()
	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	body := r.Context().Value(middleware.PayloadKey).(*payload.SessionPIN)

	// get the host id
	host := &payload.Webhook{}
	err = querynator.FindOne(&payload.SessionHost{SessionID_: sessionID}, host, db, "session_host", "webhook_host", "webhook_port")
	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		return responseerror.CreateNotFoundError(map[string]string{
			"resourceName": "session_id",
		})
	default:
		return responseerror.CreateInternalServiceError(err)
	}

	// update the status to Running
	tx, err := db.Begin()
	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	err = querynator.Update(&payload.UserSession{RequestStatus: RUNNING}, []string{"session_id"}, []any{sessionID}, tx, "user_session")
	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	port, err := strconv.Atoi(host.Port)

	if err != nil {
		return responseerror.CreateInternalServiceError(err)
	}

	req := &httpx.HTTPRequest{}
	req, err = req.CreateRequest(
		"http",
		host.Host,
		port,
		VM_PIN_ENDPOINT,
		http.MethodPost,
		200,
		body,
		cf.Config,
	)

	err = req.Send(nil)

	if err != nil {
		tx.Rollback()

		if internalErr, ok := err.(*responseerror.InternalServiceError); ok {
			return internalErr
		}

		return responseerror.CreateBadRequestError(responseerror.InvalidPIN, responseerror.InvalidPINMessage, map[string]string{
			"pin": body.PIN,
		})

	}

	tmp := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  "success",
		Message: "instance has been attached successfully",
	}

	json, err := jsonutil.EncodeToJson(tmp)

	if err != nil {
		tx.Rollback()
		return responseerror.CreateInternalServiceError(err)
	}

	w.Write(json)
	tx.Commit()
	return nil

}

func terminateSessionHandler(db *sql.DB, conf interface{}, w http.ResponseWriter, r *http.Request) responseerror.HTTPCustomError {
	cf := conf.(*config.Config)
	vars := mux.Vars(r)
	sessionIDString := vars["session_id"]

	_, err := uuid.Parse(sessionIDString)
	if err != nil {
		return responseerror.CreateBadRequestError(
			responseerror.MalformedSessionID,
			responseerror.MalformedSessionIDMessage,
			map[string]string{
				"id": "session_id",
			},
		)
	}

	req := &httpx.HTTPRequest{}
	req, err = req.CreateRequest(
		"http",
		cf.Service.SessionManager.Host,
		cf.Service.SessionManager.Port,
		fmt.Sprintf(SESSION_MANAGER_DELETE_ENDPOINT, sessionIDString),
		http.MethodDelete,
		200,
		nil,
		cf.Config,
	)

	err = req.Send(nil)

	if err != nil {
		if internalErr, ok := err.(*responseerror.InternalServiceError); ok {
			return internalErr
		}

		return err.(responseerror.HTTPCustomError)
	}

	return nil

}

func SetSessionRoute(r *mux.Router, db *sql.DB, conf *config.Config) {
	subrouter := r.PathPrefix("/session").Subrouter()

	subrouter.Use(middleware.RouteGetterMiddleware)

	authMiddleware := middleware.AuthMiddleware(conf.Service.Auth, conf.Config)
	certMiddleware := middleware.CertMiddleware(conf.RootCAs)

	createNewSessionPayloadMiddleware, err := middleware.PayloadCheckMiddleware(&payload.UserSession{},
		"Username",
		"SessionMetadata:GameID",
		"SessionMetadata:GPUName",
		"SessionMetadata:GameLocation:Protocol",
		"SessionMetadata:GameLocation:Path",
		"SessionMetadata:GameLocation:Server:IP",
	)

	if err != nil {
		log.Fatal(err)
	}

	startConnectionPayloadMiddleware, err := middleware.PayloadCheckMiddleware(&payload.SessionHost{},
		"Webhook:Port",
		"Webhook:Host",
		"NetworkID",
	)

	if err != nil {
		log.Fatal(err)
	}

	pinPairPayloadMiddleware, err := middleware.PayloadCheckMiddleware(&payload.SessionPIN{}, "PIN")

	if err != nil {
		log.Fatal(err)
	}

	createSession := httpx.CreateHTTPHandler(db, conf, createNewSessionHandler)
	subrouter.Handle("/create", middleware.UseMiddleware(db, conf, createSession, certMiddleware, createNewSessionPayloadMiddleware))

	getGPUStatus := httpx.CreateHTTPHandler(db, conf, getGPUStatusHandler)
	subrouter.Handle("/gpu", getGPUStatus).Methods("GET")

	getStatus := httpx.CreateHTTPHandler(db, conf, getRequestStatus)
	subrouter.Handle("/{session_id}/status", middleware.UseMiddleware(db, conf, getStatus, authMiddleware)).Methods("GET")

	startConnection := httpx.CreateHTTPHandler(db, conf, startConnectionEstablishmentHandler)
	subrouter.Handle("/{session_id}/connection/start", middleware.UseMiddleware(db, conf, startConnection, startConnectionPayloadMiddleware)).Methods("POST")

	pair := httpx.CreateHTTPHandler(db, conf, pairHandler)
	subrouter.Handle("/{session_id}/pair", middleware.UseMiddleware(db, conf, pair, authMiddleware, pinPairPayloadMiddleware)).Methods("POST")

	terminate := httpx.CreateHTTPHandler(db, conf, terminateSessionHandler)
	subrouter.Handle("/{session_id}/terminate", middleware.UseMiddleware(db, conf, terminate, authMiddleware)).Methods("DELETE")
}
