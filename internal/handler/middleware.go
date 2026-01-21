package handler

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// contextKey тип для предотвращения коллизий ключей контекста
type contextKey string

// Ключи контекста
const (
	requestIDContextKey contextKey = "request_id"
)

// requestIDMiddleware добавляет заголовок request_id если его нет
func (h *Handler) requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
			r.Header.Set("X-Request-ID", requestID)
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, requestIDContextKey, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// loggingMiddleware логирует информацию о запросе
func (h *Handler) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID, _ := ctx.Value(requestIDContextKey).(string)

		h.logger.Info("HTTP запрос получен",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("request_id", requestID),
		)

		next.ServeHTTP(w, r)
	})
}
