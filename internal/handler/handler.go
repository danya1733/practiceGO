package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/danya1733/practiceGO/internal/domain"
	"github.com/danya1733/practiceGO/internal/repository"
	"github.com/danya1733/practiceGO/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Тип контекстного ключа для предотвращения коллизий
type contextKey string

// Ключи контекста
const (
	requestIDContextKey contextKey = "request_id"
)

// Handler представляет обработчик HTTP запросов
type Handler struct {
	warehouseRepo *repository.WarehouseRepository
	productRepo   *repository.ProductRepository
	inventoryRepo *repository.InventoryRepository
	analyticsRepo *repository.AnalyticsRepository
	logger        *logger.Logger
}

// NewHandler создает новый обработчик HTTP запросов
func NewHandler(
	warehouseRepo *repository.WarehouseRepository,
	productRepo *repository.ProductRepository,
	inventoryRepo *repository.InventoryRepository,
	analyticsRepo *repository.AnalyticsRepository,
	logger *logger.Logger,
) *Handler {
	return &Handler{
		warehouseRepo: warehouseRepo,
		productRepo:   productRepo,
		inventoryRepo: inventoryRepo,
		analyticsRepo: analyticsRepo,
		logger:        logger,
	}
}

// RegisterRoutes регистрирует маршруты HTTP
func (h *Handler) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Swagger UI - используем FileServer для статических файлов
	fileServer := http.FileServer(http.Dir("./docs/swagger"))

	// Update these handlers for better compatibility
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fileServer))

	// Swagger JSON - explicitly serve the swagger.json file
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, "./docs/swagger/swagger.json")
	})

	// Redirect from /swagger to /swagger/
	mux.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})

	// Обработчик для проверки здоровья
	mux.HandleFunc("GET /api/health", h.HealthCheck)

	// Маршруты для работы со складами
	mux.HandleFunc("GET /api/warehouses", h.GetWarehouses)
	mux.HandleFunc("POST /api/warehouses", h.CreateWarehouse)

	// Маршруты для работы с товарами
	mux.HandleFunc("GET /api/products", h.GetProducts)
	mux.HandleFunc("POST /api/products", h.CreateProduct)
	mux.HandleFunc("PUT /api/products/{id}", h.UpdateProduct)

	// Маршруты для работы с инвентаризацией
	mux.HandleFunc("POST /api/inventory", h.CreateInventory)
	mux.HandleFunc("PUT /api/inventory/quantity", h.UpdateInventoryQuantity)
	mux.HandleFunc("PUT /api/inventory/discount", h.UpdateInventoryDiscount)
	mux.HandleFunc("GET /api/warehouses/{id}/products", h.GetWarehouseProducts)
	mux.HandleFunc("GET /api/warehouses/{warehouse_id}/products/{product_id}", h.GetWarehouseProduct)
	mux.HandleFunc("POST /api/warehouses/calculate", h.CalculateProductsPrice)
	mux.HandleFunc("POST /api/warehouses/purchase", h.PurchaseProducts)

	// Маршруты для работы с аналитикой
	mux.HandleFunc("GET /api/analytics/warehouses/{id}", h.GetWarehouseAnalytics)
	mux.HandleFunc("GET /api/analytics/warehouses/top", h.GetTopWarehouses)

	// Применение middleware для логирования и обработки request_id
	return h.requestIDMiddleware(h.loggingMiddleware(mux))
}

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
		http.Error(w, "Ошибка при получении списка складов", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(warehouses)
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
		http.Error(w, "Некорректный формат запроса", http.StatusBadRequest)
		return
	}

	createdWarehouse, err := h.warehouseRepo.Create(ctx, warehouse)
	if err != nil {
		logger.Error("Ошибка при создании склада", zap.Error(err))
		http.Error(w, "Ошибка при создании склада", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdWarehouse)
}

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
		http.Error(w, "Ошибка при получении списка товаров", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
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
		http.Error(w, "Некорректный формат запроса", http.StatusBadRequest)
		return
	}

	createdProduct, err := h.productRepo.Create(ctx, product)
	if err != nil {
		logger.Error("Ошибка при создании товара", zap.Error(err))
		http.Error(w, "Ошибка при создании товара", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdProduct)
}

// UpdateProduct обновляет существующий товар
func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithRequestID(ctx)

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Error("Некорректный формат ID", zap.Error(err))
		http.Error(w, "Некорректный формат ID", http.StatusBadRequest)
		return
	}

	var product domain.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		logger.Error("Ошибка при декодировании запроса", zap.Error(err))
		http.Error(w, "Некорректный формат запроса", http.StatusBadRequest)
		return
	}

	product.ID = id
	updatedProduct, err := h.productRepo.Update(ctx, product)
	if err != nil {
		logger.Error("Ошибка при обновлении товара", zap.Error(err))
		http.Error(w, "Ошибка при обновлении товара", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedProduct)
}

// CreateInventory создает новую запись инвентаризации
func (h *Handler) CreateInventory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithRequestID(ctx)

	var inventory domain.Inventory
	if err := json.NewDecoder(r.Body).Decode(&inventory); err != nil {
		logger.Error("Ошибка при декодировании запроса", zap.Error(err))
		http.Error(w, "Некорректный формат запроса", http.StatusBadRequest)
		return
	}

	createdInventory, err := h.inventoryRepo.Create(ctx, inventory)
	if err != nil {
		logger.Error("Ошибка при создании записи инвентаризации", zap.Error(err))
		http.Error(w, "Ошибка при создании записи инвентаризации", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdInventory)
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
		http.Error(w, "Некорректный формат запроса", http.StatusBadRequest)
		return
	}

	warehouseID, err := uuid.Parse(data.WarehouseID)
	if err != nil {
		logger.Error("Некорректный формат ID склада", zap.Error(err))
		http.Error(w, "Некорректный формат ID склада", http.StatusBadRequest)
		return
	}

	productID, err := uuid.Parse(data.ProductID)
	if err != nil {
		logger.Error("Некорректный формат ID товара", zap.Error(err))
		http.Error(w, "Некорректный формат ID товара", http.StatusBadRequest)
		return
	}

	updatedInventory, err := h.inventoryRepo.UpdateQuantity(ctx, warehouseID, productID, data.Quantity)
	if err != nil {
		logger.Error("Ошибка при обновлении количества товара", zap.Error(err))
		http.Error(w, "Ошибка при обновлении количества товара", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedInventory)
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
		http.Error(w, "Некорректный формат запроса", http.StatusBadRequest)
		return
	}

	warehouseID, err := uuid.Parse(data.WarehouseID)
	if err != nil {
		logger.Error("Некорректный формат ID склада", zap.Error(err))
		http.Error(w, "Некорректный формат ID склада", http.StatusBadRequest)
		return
	}

	productID, err := uuid.Parse(data.ProductID)
	if err != nil {
		logger.Error("Некорректный формат ID товара", zap.Error(err))
		http.Error(w, "Некорректный формат ID товара", http.StatusBadRequest)
		return
	}

	updatedInventory, err := h.inventoryRepo.UpdateDiscount(ctx, warehouseID, productID, data.Discount)
	if err != nil {
		logger.Error("Ошибка при обновлении скидки", zap.Error(err))
		http.Error(w, "Ошибка при обновлении скидки", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedInventory)
}

// GetWarehouseProducts возвращает список товаров на складе
func (h *Handler) GetWarehouseProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithRequestID(ctx)

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Error("Некорректный формат ID склада", zap.Error(err))
		http.Error(w, "Некорректный формат ID склада", http.StatusBadRequest)
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
		http.Error(w, "Ошибка при получении списка товаров на складе", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// GetWarehouseProduct возвращает конкретный товар на складе
func (h *Handler) GetWarehouseProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithRequestID(ctx)

	warehouseIDStr := r.PathValue("warehouse_id")
	warehouseID, err := uuid.Parse(warehouseIDStr)
	if err != nil {
		logger.Error("Некорректный формат ID склада", zap.Error(err))
		http.Error(w, "Некорректный формат ID склада", http.StatusBadRequest)
		return
	}

	productIDStr := r.PathValue("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		logger.Error("Некорректный формат ID товара", zap.Error(err))
		http.Error(w, "Некорректный формат ID товара", http.StatusBadRequest)
		return
	}

	inventory, err := h.inventoryRepo.GetByWarehouseAndProduct(ctx, warehouseID, productID)
	if err != nil {
		logger.Error("Ошибка при получении товара на складе", zap.Error(err))
		http.Error(w, "Ошибка при получении товара на складе", http.StatusInternalServerError)
		return
	}

	// Получаем информацию о товаре
	product, err := h.productRepo.GetByID(ctx, productID)
	if err != nil {
		logger.Error("Ошибка при получении информации о товаре", zap.Error(err))
		http.Error(w, "Ошибка при получении информации о товаре", http.StatusInternalServerError)
		return
	}

	result := domain.InventoryWithProduct{
		Inventory: inventory,
		Product:   product,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// CalculateProductsPrice рассчитывает стоимость товаров с учетом скидок
func (h *Handler) CalculateProductsPrice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithRequestID(ctx)

	var request domain.PurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Error("Ошибка при декодировании запроса")
		http.Error(w, "Некорректный формат запроса", http.StatusBadRequest)
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
			http.Error(w, "Ошибка при расчете стоимости: товар не найден на складе", http.StatusBadRequest)
			return
		}

		// Получаем информацию о товаре
		product, err := h.productRepo.GetByID(ctx, p.ProductID)
		if err != nil {
			logger.Error("Ошибка при получении информации о товаре",
				zap.Error(err),
				zap.String("product_id", p.ProductID.String()))
			http.Error(w, "Ошибка при расчете стоимости: товар не найден", http.StatusBadRequest)
			return
		}

		// Проверяем наличие достаточного количества товара
		if inventory.Quantity < p.Quantity {
			logger.Error("Недостаточное количество товара на складе",
				zap.String("product_id", p.ProductID.String()),
				zap.Int("available", inventory.Quantity),
				zap.Int("requested", p.Quantity))
			http.Error(w, "Недостаточное количество товара на складе", http.StatusBadRequest)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// PurchaseProducts обрабатывает покупку товаров
func (h *Handler) PurchaseProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithRequestID(ctx)

	var request domain.PurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Error("Ошибка при декодировании запроса", zap.Error(err))
		http.Error(w, "Некорректный формат запроса", http.StatusBadRequest)
		return
	}

	err := h.inventoryRepo.PurchaseProducts(ctx, request.WarehouseID, request.Products)
	if err != nil {
		logger.Error("Ошибка при обработке покупки", zap.Error(err))
		http.Error(w, "Ошибка при обработке покупки: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}

// GetWarehouseAnalytics возвращает аналитику по складу
func (h *Handler) GetWarehouseAnalytics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithRequestID(ctx)

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Error("Некорректный формат ID склада", zap.Error(err))
		http.Error(w, "Некорректный формат ID склада", http.StatusBadRequest)
		return
	}

	analytics, totalSum, err := h.analyticsRepo.GetWarehouseAnalytics(ctx, id)
	if err != nil {
		logger.Error("Ошибка при получении аналитики по складу", zap.Error(err))
		http.Error(w, "Ошибка при получении аналитики по складу", http.StatusInternalServerError)
		return
	}

	result := struct {
		TotalSum  float64            `json:"total_sum"`
		Analytics []domain.Analytics `json:"analytics"`
	}{
		TotalSum:  totalSum,
		Analytics: analytics,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
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
		http.Error(w, "Ошибка при получении топ складов", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(warehouses)
}
