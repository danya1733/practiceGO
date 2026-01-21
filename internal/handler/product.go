package handler

import (
	"encoding/json"
	"net/http"

	"github.com/danya1733/practiceGO/internal/domain"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// GetProducts возвращает список всех товаров
// @Summary Получить список всех товаров
// @Description Возвращает список всех товаров в системе
// @Tags products
// @Produce json
// @Success 200 {array} domain.Product
// @Failure 500 {object} map[string]string
// @Router /api/products [get]
func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithRequestID(ctx)

	products, err := h.productRepo.GetAll(ctx)
	if err != nil {
		logger.Error("Ошибка при получении списка товаров", zap.Error(err))
		writeError(w, "Ошибка при получении списка товаров", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, products)
}

// CreateProduct создает новый товар
// @Summary Создать новый товар
// @Description Создает новый товар в системе
// @Tags products
// @Accept json
// @Produce json
// @Param product body domain.Product true "Информация о товаре"
// @Success 201 {object} domain.Product
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/products [post]
func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithRequestID(ctx)

	var product domain.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		logger.Error("Ошибка при декодировании запроса", zap.Error(err))
		writeError(w, "Некорректный формат запроса", http.StatusBadRequest)
		return
	}

	createdProduct, err := h.productRepo.Create(ctx, product)
	if err != nil {
		logger.Error("Ошибка при создании товара", zap.Error(err))
		writeError(w, "Ошибка при создании товара", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, createdProduct)
}

// UpdateProduct обновляет существующий товар
func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithRequestID(ctx)

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Error("Некорректный формат ID", zap.Error(err))
		writeError(w, "Некорректный формат ID", http.StatusBadRequest)
		return
	}

	var product domain.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		logger.Error("Ошибка при декодировании запроса", zap.Error(err))
		writeError(w, "Некорректный формат запроса", http.StatusBadRequest)
		return
	}

	product.ID = id
	updatedProduct, err := h.productRepo.Update(ctx, product)
	if err != nil {
		logger.Error("Ошибка при обновлении товара", zap.Error(err))
		writeError(w, "Ошибка при обновлении товара", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, updatedProduct)
}
