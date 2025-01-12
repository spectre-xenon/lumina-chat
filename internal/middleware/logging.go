package middleware

import (
	"log"
	"net/http"
	"time"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	if w.wroteHeader {
		return
	}

	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
	w.wroteHeader = true
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		log.Println(wrapped.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}
