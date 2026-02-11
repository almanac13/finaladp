package middleware

import "net/http"

// Small helper middleware: ensures JSON content-type on responses if handler forgot.
// (Not required, but keeps demo clean.)
func WithJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// default response type
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}
