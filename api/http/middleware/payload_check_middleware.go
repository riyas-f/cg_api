package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"reflect"

	httpx "github.com/AdityaP1502/Instant-Messanging/api/http"
	"github.com/AdityaP1502/Instant-Messanging/api/http/responseerror"
)

var PayloadKey ContextKey = "payload"

func PayloadCheckMiddleware(template httpx.Payload, requiredFields ...string) (Middleware, error) {
	var payload httpx.Payload

	p := reflect.ValueOf(template)

	if p.Kind() != reflect.Ptr {
		err := fmt.Errorf("cannot create middleware. template isn't a pointer")
		return nil, err
	}

	// Struct reflect
	s := p.Elem()
	sType := s.Type()

	// Check if the requiredFields is valid
	for _, field := range requiredFields {
		v := s.FieldByName(field)
		if !v.IsValid() {
			return nil, fmt.Errorf("struct of %s don't have field named %s", s.Type(), field)
		}

		f, _ := sType.FieldByName(field)

		tag := f.Tag.Get("json")

		if tag == "" || tag == "-" {
			return nil, fmt.Errorf("required fields of %s have json tag of %s, which isn't valid", field, tag)
		}
	}

	return func(next http.Handler, db *sql.DB, config interface{}) http.Handler {
		fn := func(db *sql.DB, config interface{}, w http.ResponseWriter, r *http.Request) responseerror.HTTPCustomError {

			if r.Header.Get("Content-Type") != "application/json" {
				return responseerror.CreateBadRequestError(
					responseerror.HeaderValueMistmatch,
					responseerror.HeaderValueMistmatchMessage,
					map[string]string{
						"name": "Content-Type",
					},
				)
			}

			// create a new struct with the same type as template
			payload = reflect.New(sType).Interface().(httpx.Payload)

			if r.Body == nil {
				return responseerror.CreateBadRequestError(
					responseerror.PayloadInvalid,
					responseerror.PayloadInvalidMessage,
					nil,
				)
			}

			err := payload.FromJSON(r.Body, true, requiredFields)

			if err != nil {
				return err
			}

			ctx := context.WithValue(r.Context(), PayloadKey, payload)
			next.ServeHTTP(w, r.WithContext(ctx))

			return nil
		}

		return &httpx.Handler{
			DB:      db,
			Config:  config,
			Handler: httpx.HandlerLogic(fn),
		}
	}, nil
}
