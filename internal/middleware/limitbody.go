package middleware

import "net/http"

func LimitBodySize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Limit request body size
		r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10 MB

		next.ServeHTTP(w, r)
	})
}
