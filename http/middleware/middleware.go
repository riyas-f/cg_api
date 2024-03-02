package middleware

import (
	"database/sql"
	"net/http"
)

type ContextKey string

type Middleware func(next http.Handler, db *sql.DB, conf interface{}) http.Handler

func UseMiddleware(db *sql.DB, conf interface{}, handler http.Handler, middlewares ...Middleware) http.Handler {
	chained := handler

	for i := len(middlewares) - 1; i > -1; i-- {
		chained = middlewares[i](chained, db, conf)
	}

	return chained
}
