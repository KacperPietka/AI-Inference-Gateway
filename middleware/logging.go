package middleware

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"inference-gateway/metrics"
	"inference-gateway/types"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// Creates a wrapper with a default status code of 200
func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func GetModel(r *http.Request) string {
	if model, ok := r.Context().Value(types.ModelKey).(string); ok {
		return model
	}
	return ""
}

func NewLogger() *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)
}

func Logger(logger *slog.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := newResponseWriter(w)

		// Extract the user ID
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			userID = "anonymous"
		}

		next(wrapped, r)

		duration := time.Since(start)

		// Record metrics
		metrics.RequestsTotal.WithLabelValues(
			r.Method,
			r.URL.Path,
			strconv.Itoa(wrapped.statusCode),
		).Inc()

		metrics.RequestDuration.WithLabelValues(
			r.Method,
			r.URL.Path,
		).Observe(duration.Seconds())

		model := GetModel(r)

		if model != "" {
			logger.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", wrapped.statusCode,
				"duration", duration.String(),
				"remoted_addr", r.RemoteAddr, // allows to see where requests are coming from
				"user_id", userID,
				"request_id", GetRequestID(r),
				"model", model,
			)
		} else {
			logger.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", wrapped.statusCode,
				"duration", time.Since(start).String(),
				"remote_addr", r.RemoteAddr,
				"user_id", userID,
				"request_id", GetRequestID(r),
			)
		}
	}
}
