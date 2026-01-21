package handler

import (
	"net/http"
	"strconv"

	"github.com/danya1733/practiceGO/internal/domain"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// GetWarehouseAnalytics возвращает аналитику по складу
func (h *Handler) GetWarehouseAnalytics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithRequestID(ctx)

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Error("Некорректный формат ID склада", zap.Error(err))
		writeError(w, "Некорректный формат ID склада", http.StatusBadRequest)
		return
	}

	analytics, totalSum, err := h.analyticsRepo.GetWarehouseAnalytics(ctx, id)
	if err != nil {
		logger.Error("Ошибка при получении аналитики по складу", zap.Error(err))
		writeError(w, "Ошибка при получении аналитики по складу", http.StatusInternalServerError)
		return
	}

	result := struct {
		TotalSum  float64            `json:"total_sum"`
		Analytics []domain.Analytics `json:"analytics"`
	}{
		TotalSum:  totalSum,
		Analytics: analytics,
	}

	writeJSON(w, http.StatusOK, result)
}

// GetTopWarehouses возвращает топ складов по выручке
func (h *Handler) GetTopWarehouses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithRequestID(ctx)

	limit := 5
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		limitVal, err := strconv.Atoi(limitStr)
		if err == nil && limitVal > 0 {
			limit = limitVal
		}
	}

	warehouses, err := h.analyticsRepo.GetTopWarehouses(ctx, limit)
	if err != nil {
		logger.Error("Ошибка при получении топ складов", zap.Error(err))
		writeError(w, "Ошибка при получении топ складов", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, warehouses)
}
