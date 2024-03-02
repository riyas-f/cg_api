package middleware

import (
	"context"
	"net/http"
	"strings"
)

var EndpointKey ContextKey = "endpoint"

func RouteGetterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		pathSplit := strings.SplitN(path, "/", 3)
		endpoint := pathSplit[len(pathSplit)-1]

		ctx := context.WithValue(r.Context(), EndpointKey, endpoint)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
