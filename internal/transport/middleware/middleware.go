package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func HTTPLogger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrappedRW := wrapResponseWriter(w)
			next.ServeHTTP(wrappedRW, r)
			log.Info(
				"incoming request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.EscapedPath()),
				slog.Int("status", wrappedRW.status),
				slog.Duration("duration", time.Since(start)),
			)
		}
		return http.HandlerFunc(fn)
	}
}
