package middleware

import (
	"net/http"
	"time"
	"log/slog"
)

type StatusRecorder struct {
	http.ResponseWriter
	StatusCode int
}

func (rec *StatusRecorder) WriteHeader(code int) {
	rec.StatusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		start := time.Now()

		recorder := &StatusRecorder{
			ResponseWriter: w,
			StatusCode: http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		duration := time.Since(start)

		slog.Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", recorder.StatusCode,
			"duration_ms", duration.Milliseconds(),)
	})
}