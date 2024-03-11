package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/LaQuannT/astronaut-data-api/internal/model"
	"github.com/LaQuannT/astronaut-data-api/internal/transport/util"
	"github.com/google/uuid"
)

type (
	apikeyHeader string
	requestUser  string
)

const (
	APIKeyHeader apikeyHeader = "X-api-key"
	RequestUser  requestUser  = "request-user"
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

func APIKeyValidation(uc model.UserUsecase, log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get(string(APIKeyHeader))

			if key == "" {
				err := errors.New("invalid APIKey")
				log.Warn("failed APIKey validation", slog.Any("error", err))
				util.WriteJSON(w, http.StatusUnauthorized, model.JSONResponse{Error: "Permission denied"})
				return
			}

			_, err := uuid.Parse(key)
			if err != nil {
				log.Warn("faild APIKey validation", slog.Any("error", err))
				util.WriteJSON(w, http.StatusUnauthorized, model.JSONResponse{Error: "Permission denied"})
				return
			}

			ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
			defer cancel()

			u, err := uc.SearchAPIKey(ctx, key)
			if err != nil {
				log.Warn("failed APIKey validation user search", slog.Any("error", err))
				util.WriteJSON(w, http.StatusUnauthorized, model.JSONResponse{Error: "Permission denied"})
				return
			}

			if u == nil {
				util.WriteJSON(w, http.StatusUnauthorized, model.JSONResponse{Error: "Permission denied"})
				return
			} else {
				ctx = context.WithValue(ctx, RequestUser, u)
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
