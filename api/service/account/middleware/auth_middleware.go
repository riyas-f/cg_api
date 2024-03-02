package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	httpx "github.com/AdityaP1502/Instant-Messanging/api/http"
	"github.com/AdityaP1502/Instant-Messanging/api/http/middleware"
	"github.com/AdityaP1502/Instant-Messanging/api/http/responseerror"
	"github.com/AdityaP1502/Instant-Messanging/api/service/account/config"
	"github.com/AdityaP1502/Instant-Messanging/api/service/auth/jwtutil"
)

var ClaimsKey middleware.ContextKey = "claims"
var TokenKey middleware.ContextKey = "token"

var AUTH_VERIFY_TOKEN_ENDPOINT = "v1/auth/token/verify"

func AuthMiddleware(next http.Handler, db *sql.DB, conf interface{}) http.Handler {
	fn := func(db *sql.DB, conf interface{}, w http.ResponseWriter, r *http.Request) responseerror.HTTPCustomError {
		var token string

		endpoint := r.Context().Value(middleware.EndpointKey).(string)

		cf := conf.(*config.Config)

		auth := r.Header.Get("Authorization")

		if auth == "" {
			return responseerror.CreateBadRequestError(
				responseerror.EmptyAuthHeader,
				responseerror.EmptyAuthHeaderMessage,
				nil,
			)
		}

		if authType, authValue, _ := strings.Cut(auth, " "); authType != "Bearer" {
			return responseerror.CreateUnauthorizedError(
				responseerror.InvalidAuthHeader,
				responseerror.InvalidAuthHeaderMessage,
				map[string]string{
					"authType": authType,
				},
			)

		} else {
			token = authValue
		}

		req := &httpx.HTTPRequest{}
		req, err := req.CreateRequest(
			cf.Services.Auth.Scheme,
			cf.Services.Auth.Host,
			cf.Services.Auth.Port,
			AUTH_VERIFY_TOKEN_ENDPOINT,
			http.MethodPost,
			http.StatusOK,
			struct {
				Token    string `json:"token"`
				Endpoint string `json:"endpoint"`
			}{
				Token:    token,
				Endpoint: endpoint,
			},
			cf.Config,
		)

		if err != nil {
			return responseerror.CreateInternalServiceError(err)
		}

		claims := &jwtutil.Claims{}
		err = req.Send(claims)

		if err != nil {
			return err
		}

		// send token and claims into the next middleware chain
		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		ctx = context.WithValue(ctx, TokenKey, token)

		next.ServeHTTP(w, r.WithContext(ctx))
		return nil
	}

	return &httpx.Handler{
		DB:      db,
		Config:  conf,
		Handler: httpx.HandlerLogic(fn),
	}

}
