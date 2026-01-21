package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// writeJSON записывает JSON ответ
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError записывает ошибку в ответ
func writeError(w http.ResponseWriter, message string, status int) {
	http.Error(w, message, status)
}

// parseUUID парсит UUID из строки и логирует ошибку если не удалось
func parseUUID(s string, logger *zap.Logger, fieldName string) (uuid.UUID, bool) {
	id, err := uuid.Parse(s)
	if err != nil {
		logger.Error("Некорректный формат ID", zap.String("field", fieldName), zap.Error(err))
		return uuid.Nil, false
	}
	return id, true
}

// HealthCheck проверяет работоспособность сервиса
// @Summary Проверка работоспособности сервиса
// @Description Возвращает статус 200 OK если сервис работает
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/health [get]
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
