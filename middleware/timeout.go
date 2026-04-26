package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	gwerrors "inference-gateway/errors"
)

func Timeout(duration time.Duration, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithTimeout(r.Context(), duration)
		defer cancel()

		done := make(chan struct{})

		wrapped := newResponseWriter(w)

		go func() {
			next(wrapped, r.WithContext(ctx))
			close(done)
		}()

		select {
		case <-done:
			return

		case <-ctx.Done():
			if wrapped.statusCode == http.StatusOK {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusRequestTimeout)
				json.NewEncoder(w).Encode(gwerrors.ErrRequestTimeout)
			}
			return
		}
	}
}
