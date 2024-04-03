package middleware

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/rs/cors"
)

type ContextKey string

type Middleware func(next http.Handler, db *sql.DB, conf interface{}) http.Handler

func UseMiddleware(db *sql.DB, conf interface{}, handler http.Handler, corsOptions *cors.Options, middlewares ...Middleware) http.Handler {
	chained := handler

	for i := len(middlewares) - 1; i > -1; i-- {
		chained = middlewares[i](chained, db, conf)
	}

	// append cors middleware into the chain
	if corsOptions != nil {
		fmt.Println(*corsOptions)
		chained = cors.New(*corsOptions).Handler(chained)
	}

	return chained
}
