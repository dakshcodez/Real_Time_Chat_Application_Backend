package middleware

import (
	"net/http"

	"github.com/dakshcodez/gdg_chat_app_backend_task/internal/ratelimit"
	"github.com/google/uuid"
)

func RateLimit(limiter *ratelimit.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			if !limiter.Allow(userID.String()) {
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
