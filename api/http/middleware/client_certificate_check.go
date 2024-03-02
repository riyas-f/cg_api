package middleware

import (
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"net/http"
	"time"

	httpx "github.com/AdityaP1502/Instant-Messanging/api/http"
	"github.com/AdityaP1502/Instant-Messanging/api/http/responseerror"
)

func CertMiddleware(rootCACerts *x509.CertPool) Middleware {
	return func(next http.Handler, db *sql.DB, conf interface{}) http.Handler {
		fn := func(db *sql.DB, conf interface{}, w http.ResponseWriter, r *http.Request) responseerror.HTTPCustomError {
			if r.TLS != nil && len(r.TLS.PeerCertificates) > 0 {
				next.ServeHTTP(w, r)
			}

			// client request is being redirected by the proxy
			if x := r.Header.Get("x-client-cert"); r.TLS != nil && x != "" {
				certBytes, err := base64.StdEncoding.DecodeString(x)
				if err != nil {
					return responseerror.CreateInternalServiceError(err)
				}

				cert, err := x509.ParseCertificate(certBytes)

				if err != nil {
					return responseerror.CreateInternalServiceError(err)
				}

				opts := x509.VerifyOptions{
					Roots:       rootCACerts,
					CurrentTime: time.Now(),
				}

				chains, err := cert.Verify(opts)

				if err != nil {
					return responseerror.CreateInternalServiceError(err)
				}

				if len(chains) > 0 {
					next.ServeHTTP(w, r)
				}

			}

			return responseerror.CreateUnauthorizedError(
				responseerror.AccessDenied,
				responseerror.AccessDeniedMessage,
				nil,
			)
		}

		return &httpx.Handler{
			DB:      db,
			Config:  conf,
			Handler: httpx.HandlerLogic(fn),
		}
	}
}
