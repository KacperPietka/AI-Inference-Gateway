package middleware

import (
	"net/http"
	"strconv"

	gwerrors "inference-gateway/errors"
	"inference-gateway/ratelimit"
)

// Ratelimit middleware checks the user's request count before
// allowing the request through the handler
func RateLimit(limiter *ratelimit.RateLimiter, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			userID = "anonymous"
		}

		allowed, err := limiter.Allow(r.Context(), userID)
		if err != nil {
			next(w, r)
			return
		}

		if !allowed {
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(10))
			w.Header().Set("X-RateLimit-Remaining", "0")
			writeRateLimitError(w)
			return
		}

		remaining, err := limiter.Remaining(r.Context(), userID)
		if err == nil {
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(10))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))

		}

		next(w, r)
	}
}

func writeRateLimitError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(gwerrors.ErrRateLimited.Code)

	w.Write([]byte(`{"error":"rate limit exceeded", "code":429}`))
}
