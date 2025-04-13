package app

import (
	"net/http"

	"github.com/danya1733/practiceGO/internal/config"
	"github.com/danya1733/practiceGO/internal/handler"
	"github.com/danya1733/practiceGO/internal/repository"
	"github.com/danya1733/practiceGO/pkg/logger"
)

// App представляет приложение
type App struct {
	cfg           *config.Config
	logger        *logger.Logger
	db            *repository.PostgresDB
	handler       *handler.Handler
	warehouseRepo *repository.WarehouseRepository
	productRepo   *repository.ProductRepository
	inventoryRepo *repository.InventoryRepository
	analyticsRepo *repository.AnalyticsRepository
}

// NewApp создает новое приложение
func NewApp(cfg *config.Config, logger *logger.Logger) (*App, error) {
	// Инициализация базы данных
	db, err := repository.NewPostgresDB(cfg.Database)
	if err != nil {
		return nil, err
	}

	// Инициализация репозиториев
	warehouseRepo := repository.NewWarehouseRepository(db.GetPool())
	productRepo := repository.NewProductRepository(db.GetPool())
	inventoryRepo := repository.NewInventoryRepository(db.GetPool())
	analyticsRepo := repository.NewAnalyticsRepository(db.GetPool())

	// Инициализация обработчика HTTP запросов
	h := handler.NewHandler(warehouseRepo, productRepo, inventoryRepo, analyticsRepo, logger)

	return &App{
		cfg:           cfg,
		logger:        logger,
		db:            db,
		handler:       h,
		warehouseRepo: warehouseRepo,
		productRepo:   productRepo,
		inventoryRepo: inventoryRepo,
		analyticsRepo: analyticsRepo,
	}, nil
}

// Router возвращает маршрутизатор HTTP
func (a *App) Router() http.Handler {
	return a.handler.RegisterRoutes()
}

// Close закрывает соединения
func (a *App) Close() error {
	a.db.Close()
	return nil
}
