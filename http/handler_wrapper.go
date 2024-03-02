package httpx

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/AdityaP1502/Instant-Messanging/api/http/responseerror"
	"github.com/AdityaP1502/Instant-Messanging/api/jsonutil"
)

type HandlerLogic func(db *sql.DB, conf interface{}, w http.ResponseWriter, r *http.Request) responseerror.HTTPCustomError

type Handler struct {
	DB      *sql.DB
	Config  interface{}
	Handler HandlerLogic
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if err := h.Handler(h.DB, h.Config, w, r); err != nil {
		if internalErr, ok := err.(*responseerror.InternalServiceError); ok {
			fmt.Println(internalErr.Description)
		}

		requestErr := err.Get()

		errorResponse := responseerror.FailedRequestResponse{
			Status:    "fail",
			ErrorType: requestErr.Name,
			Message:   requestErr.Message,
		}

		w.WriteHeader(requestErr.Code)

		json, err := jsonutil.EncodeToJson(&errorResponse)

		if err != nil {
			http.Error(w, "Something wrong with server!", 500)
		}

		w.Write(json)
	}
}

func CreateHTTPHandler(db *sql.DB, conf interface{}, logic HandlerLogic) http.Handler {
	return &Handler{
		DB:      db,
		Config:  conf,
		Handler: logic,
	}
}
