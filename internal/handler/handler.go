package handler

import (
	"net/http"

	"github.com/danya1733/practiceGO/internal/repository"
	"github.com/danya1733/practiceGO/pkg/logger"
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
