package middleware

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"

	"inference-gateway/types"
)

func generateRequestID() string {
	bytes := make([]byte, 3)
	rand.Read(bytes)
	return fmt.Sprintf("req-%x", bytes)
}

func RequestID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		ctx := context.WithValue(r.Context(), types.RequestIDKey, requestID)

		w.Header().Set("X-Request-ID", requestID)

		next(w, r.WithContext(ctx))
	}
}

func GetRequestID(r *http.Request) string {
	if id, ok := r.Context().Value(types.RequestIDKey).(string); ok {
		return id
	}
	return "unknown"
}
