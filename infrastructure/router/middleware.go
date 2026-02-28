package router

import (
	"errors"
	"log/slog"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, err := GetUser(r)
		if errors.Is(err, ErrUserNotFoundInRequest) {
			slog.InfoContext(ctx, "request received",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.String("host", r.Host),
			)
		} else {
			slog.InfoContext(ctx, "request received",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.String("host", r.Host),
				slog.Any("userID", user.ID),
			)
		}

		next.ServeHTTP(w, r)
	})
}
