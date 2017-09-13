package httpx

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

// BasicAuthentication provides a middleware that prompts a simple username/password dialog
// Follows https://tools.ietf.org/html/rfc2617
func BasicAuthentication(user, password string) Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			credentials := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, password)))
			if r.Header.Get("Authorization") == fmt.Sprintf("Basic %s", credentials) {
				next.ServeHTTP(w, r)
				return
			}

			w.Header().Add("WWW-Authenticate", "Basic")
			w.WriteHeader(http.StatusUnauthorized)
		}

		return http.HandlerFunc(fn)
	}
}
