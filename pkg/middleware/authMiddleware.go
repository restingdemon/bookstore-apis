package middleware

import (
	"context"
	"fmt"
	"net/http"
	"github.com/restingdemon/go-mysql-mux/pkg/helper"
)

var AuthenticationNotRequired map[string]bool = map[string]bool{
	"/signup": true,
	"/login":  true,
}

// Authenticate is a middleware function that performs authentication
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPath := r.URL.Path
		if AuthenticationNotRequired[requestedPath] {
			// If the requested path is in AuthenticationNotRequired, skip authentication
			next.ServeHTTP(w, r)
			return
		}

		clientToken := r.Header.Get("token")
		if clientToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("No auth header provided")))
			return
		}

		claims, err := helper.ValidateToken(clientToken)
		if err != "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err))
			return
		}

		ctx := context.WithValue(r.Context(), "email", claims.Email)
		r = r.WithContext(ctx)
		
		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
