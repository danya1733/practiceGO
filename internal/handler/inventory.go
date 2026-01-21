package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/danya1733/practiceGO/internal/domain"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateInventory создает новую запись инвентаризации
func (h *Handler) CreateInventory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithRequestID(ctx)

	var inventory domain.Inventory
	if err := json.NewDecoder(r.Body).Decode(&inventory); err != nil {
		logger.Error("Ошибка при декодировании запроса", zap.Error(err))
		writeError(w, "Некорректный формат запроса", http.StatusBadRequest)
		return
	}

	createdInventory, err := h.inventoryRepo.Create(ctx, inventory)
	if err != nil {
		logger.Error("Ошибка при создании записи инвентаризации", zap.Error(err))
		writeError(w, "Ошибка при создании записи инвентаризации", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, createdInventory)
}

// UpdateInventoryQuantity обновляет количество товара на складе
func (h *Handler) UpdateInventoryQuantity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithRequestID(ctx)

	var data struct {
		WarehouseID string `json:"warehouse_id"`
		ProductID   string `json:"product_id"`
		Quantity    int    `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		logger.Error("Ошибка при декодировании запроса", zap.Error(err))
		writeError(w, "Некорректный формат запроса", http.StatusBadRequest)
		return
	}

	warehouseID, err := uuid.Parse(data.WarehouseID)
	if err != nil {
		logger.Error("Некорректный формат ID склада", zap.Error(err))
		writeError(w, "Некорректный формат ID склада", http.StatusBadRequest)
		return
	}

	productID, err := uuid.Parse(data.ProductID)
	if err != nil {
		logger.Error("Некорректный формат ID товара", zap.Error(err))
		writeError(w, "Некорректный формат ID товара", http.StatusBadRequest)
		return
	}

	updatedInventory, err := h.inventoryRepo.UpdateQuantity(ctx, warehouseID, productID, data.Quantity)
	if err != nil {
		logger.Error("Ошибка при обновлении количества товара", zap.Error(err))
		writeError(w, "Ошибка при обновлении количества товара", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, updatedInventory)
}

// UpdateInventoryDiscount обновляет скидку на товар
func (h *Handler) UpdateInventoryDiscount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithRequestID(ctx)

	var data struct {
		WarehouseID string  `json:"warehouse_id"`
		ProductID   string  `json:"product_id"`
		Discount    float64 `json:"discount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		logger.Error("Ошибка при декодировании запроса", zap.Error(err))
		writeError(w, "Некорректный формат запроса", http.StatusBadRequest)
		return
	}

	warehouseID, err := uuid.Parse(data.WarehouseID)
	if err != nil {
		logger.Error("Некорректный формат ID склада", zap.Error(err))
		writeError(w, "Некорректный формат ID склада", http.StatusBadRequest)
		return
	}

	productID, err := uuid.Parse(data.ProductID)
	if err != nil {
		logger.Error("Некорректный формат ID товара", zap.Error(err))
		writeError(w, "Некорректный формат ID товара", http.StatusBadRequest)
		return
	}

	updatedInventory, err := h.inventoryRepo.UpdateDiscount(ctx, warehouseID, productID, data.Discount)
	if err != nil {
		logger.Error("Ошибка при обновлении скидки", zap.Error(err))
		writeError(w, "Ошибка при обновлении скидки", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, updatedInventory)
}

// GetWarehouseProducts возвращает список товаров на складе
func (h *Handler) GetWarehouseProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithRequestID(ctx)

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Error("Некорректный формат ID склада", zap.Error(err))
		writeError(w, "Некорректный формат ID склада", http.StatusBadRequest)
		return
	}

	page := 1
	limit := 10

	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		pageVal, err := strconv.Atoi(pageStr)
		if err == nil && pageVal > 0 {
			page = pageVal
		}
	}

	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		limitVal, err := strconv.Atoi(limitStr)
		if err == nil && limitVal > 0 {
			limit = limitVal
		}
	}

	products, err := h.inventoryRepo.GetProductsByWarehouse(ctx, id, page, limit)
	if err != nil {
		logger.Error("Ошибка при получении списка товаров на складе", zap.Error(err))
		writeError(w, "Ошибка при получении списка товаров на складе", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, products)
}

// GetWarehouseProduct возвращает конкретный товар на складе
func (h *Handler) GetWarehouseProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithRequestID(ctx)

	warehouseIDStr := r.PathValue("warehouse_id")
	warehouseID, err := uuid.Parse(warehouseIDStr)
	if err != nil {
		logger.Error("Некорректный формат ID склада", zap.Error(err))
		writeError(w, "Некорректный формат ID склада", http.StatusBadRequest)
		return
	}

	productIDStr := r.PathValue("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		logger.Error("Некорректный формат ID товара", zap.Error(err))
		writeError(w, "Некорректный формат ID товара", http.StatusBadRequest)
		return
	}

	inventory, err := h.inventoryRepo.GetByWarehouseAndProduct(ctx, warehouseID, productID)
	if err != nil {
		logger.Error("Ошибка при получении товара на складе", zap.Error(err))
		writeError(w, "Ошибка при получении товара на складе", http.StatusInternalServerError)
		return
	}

	// Получаем информацию о товаре
	product, err := h.productRepo.GetByID(ctx, productID)
	if err != nil {
		logger.Error("Ошибка при получении информации о товаре", zap.Error(err))
		writeError(w, "Ошибка при получении информации о товаре", http.StatusInternalServerError)
		return
	}

	result := domain.InventoryWithProduct{
		Inventory: inventory,
		Product:   product,
	}

	writeJSON(w, http.StatusOK, result)
}

// CalculateProductsPrice рассчитывает стоимость товаров с учетом скидок
func (h *Handler) CalculateProductsPrice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithRequestID(ctx)

	var request domain.PurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Error("Ошибка при декодировании запроса")
		writeError(w, "Некорректный формат запроса", http.StatusBadRequest)
		return
	}

	result := domain.CalculationResult{
		TotalSum: 0,
		Items: []struct {
			ProductID         uuid.UUID `json:"product_id"`
			Name              string    `json:"name"`
			Quantity          int       `json:"quantity"`
			Price             float64   `json:"price"`
			PriceWithDiscount float64   `json:"price_with_discount"`
			TotalPrice        float64   `json:"total_price"`
		}{},
	}

	for _, p := range request.Products {
		// Получаем информацию о товаре на складе
		inventory, err := h.inventoryRepo.GetByWarehouseAndProduct(ctx, request.WarehouseID, p.ProductID)
		if err != nil {
			logger.Error("Ошибка при получении информации о товаре на складе",
				zap.Error(err),
				zap.String("product_id", p.ProductID.String()))
			writeError(w, "Ошибка при расчете стоимости: товар не найден на складе", http.StatusBadRequest)
			return
		}

		// Получаем информацию о товаре
		product, err := h.productRepo.GetByID(ctx, p.ProductID)
		if err != nil {
			logger.Error("Ошибка при получении информации о товаре",
				zap.Error(err),
				zap.String("product_id", p.ProductID.String()))
			writeError(w, "Ошибка при расчете стоимости: товар не найден", http.StatusBadRequest)
			return
		}

		// Проверяем наличие достаточного количества товара
		if inventory.Quantity < p.Quantity {
			logger.Error("Недостаточное количество товара на складе",
				zap.String("product_id", p.ProductID.String()),
				zap.Int("available", inventory.Quantity),
				zap.Int("requested", p.Quantity))
			writeError(w, "Недостаточное количество товара на складе", http.StatusBadRequest)
			return
		}

		// Рассчитываем цену с учетом скидки
		priceWithDiscount := inventory.Price * (1 - inventory.Discount/100)
		totalPrice := priceWithDiscount * float64(p.Quantity)

		item := struct {
			ProductID         uuid.UUID `json:"product_id"`
			Name              string    `json:"name"`
			Quantity          int       `json:"quantity"`
			Price             float64   `json:"price"`
			PriceWithDiscount float64   `json:"price_with_discount"`
			TotalPrice        float64   `json:"total_price"`
		}{
			ProductID:         p.ProductID,
			Name:              product.Name,
			Quantity:          p.Quantity,
			Price:             inventory.Price,
			PriceWithDiscount: priceWithDiscount,
			TotalPrice:        totalPrice,
		}

		result.Items = append(result.Items, item)
		result.TotalSum += totalPrice
	}

	writeJSON(w, http.StatusOK, result)
}

// PurchaseProducts обрабатывает покупку товаров
func (h *Handler) PurchaseProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithRequestID(ctx)

	var request domain.PurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Error("Ошибка при декодировании запроса", zap.Error(err))
		writeError(w, "Некорректный формат запроса", http.StatusBadRequest)
		return
	}

	err := h.inventoryRepo.PurchaseProducts(ctx, request.WarehouseID, request.Products)
	if err != nil {
		logger.Error("Ошибка при обработке покупки", zap.Error(err))
		writeError(w, "Ошибка при обработке покупки: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}
