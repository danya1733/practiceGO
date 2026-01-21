package handler

import (
	"encoding/json"
	"net/http"

	"github.com/danya1733/practiceGO/internal/domain"
	"go.uber.org/zap"
)

// GetWarehouses возвращает список всех складов
// @Summary Получить список всех складов
// @Description Возвращает список всех складов в системе
// @Tags warehouses
// @Produce json
// @Success 200 {array} domain.Warehouse
// @Failure 500 {object} map[string]string
// @Router /api/warehouses [get]
func (h *Handler) GetWarehouses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithRequestID(ctx)

	warehouses, err := h.warehouseRepo.GetAll(ctx)
	if err != nil {
		logger.Error("Ошибка при получении списка складов", zap.Error(err))
		writeError(w, "Ошибка при получении списка складов", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, warehouses)
}

// CreateWarehouse создает новый склад
// @Summary Создать новый склад
// @Description Создает новый склад в системе
// @Tags warehouses
// @Accept json
// @Produce json
// @Param warehouse body domain.Warehouse true "Информация о складе"
// @Success 201 {object} domain.Warehouse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/warehouses [post]
func (h *Handler) CreateWarehouse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithRequestID(ctx)

	var warehouse domain.Warehouse
	if err := json.NewDecoder(r.Body).Decode(&warehouse); err != nil {
		logger.Error("Ошибка при декодировании запроса", zap.Error(err))
		writeError(w, "Некорректный формат запроса", http.StatusBadRequest)
		return
	}

	createdWarehouse, err := h.warehouseRepo.Create(ctx, warehouse)
	if err != nil {
		logger.Error("Ошибка при создании склада", zap.Error(err))
		writeError(w, "Ошибка при создании склада", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, createdWarehouse)
}
