// Package docs предоставляет документацию Swagger для API.
//
// Документация API для системы управления складами.
//
// Схемы: http
// BasePath: /api
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package docs

import (
	"github.com/danya1733/practiceGO/internal/domain"
)

// Успешный ответ со статусом OK
// swagger:response okResponse
type OkResponse struct {
	// Статус операции
	// in: body
	Body struct {
		Status string `json:"status"`
	}
}

// Ответ с ошибкой
// swagger:response errorResponse
type ErrorResponse struct {
	// Сообщение об ошибке
	// in: body
	Body struct {
		Error string `json:"error"`
	}
}

// Запрос на обновление количества товара
// swagger:parameters updateInventoryQuantity
type UpdateQuantityRequest struct {
	// Данные для обновления количества товара
	// in: body
	Body struct {
		WarehouseID string `json:"warehouse_id"`
		ProductID   string `json:"product_id"`
		Quantity    int    `json:"quantity"`
	}
}

// Запрос на обновление скидки на товар
// swagger:parameters updateInventoryDiscount
type UpdateDiscountRequest struct {
	// Данные для обновления скидки на товар
	// in: body
	Body struct {
		WarehouseID string  `json:"warehouse_id"`
		ProductID   string  `json:"product_id"`
		Discount    float64 `json:"discount"`
	}
}

// Запрос на расчет стоимости покупки
// swagger:parameters calculateProductsPrice purchaseProducts
type PurchaseRequest struct {
	// Данные для расчета стоимости покупки
	// in: body
	Body domain.PurchaseRequest
}

// Результат расчета стоимости покупки
// swagger:response calculationResponse
type CalculationResponse struct {
	// Результат расчета стоимости
	// in: body
	Body domain.CalculationResult
}

// Результат получения аналитики по складу
// swagger:response warehouseAnalyticsResponse
type WarehouseAnalyticsResponse struct {
	// Аналитика по складу
	// in: body
	Body struct {
		TotalSum  float64            `json:"total_sum"`
		Analytics []domain.Analytics `json:"analytics"`
	}
}

// Ответ с топ-N складами по выручке
// swagger:response topWarehousesResponse
type TopWarehousesResponse struct {
	// Топ-N складов по выручке
	// in: body
	Body []domain.WarehouseAnalytics
}

// Запрос получения списка товаров на складе с пагинацией
// swagger:parameters getWarehouseProducts
type GetWarehouseProductsParams struct {
	// ID склада
	// in: path
	// required: true
	ID string `json:"id"`

	// Номер страницы, начиная с 1
	// in: query
	// default: 1
	Page int `json:"page"`

	// Количество товаров на странице
	// in: query
	// default: 10
	Limit int `json:"limit"`
}

// Запрос получения топ-N складов по выручке
// swagger:parameters getTopWarehouses
type GetTopWarehousesParams struct {
	// Количество складов в выборке
	// in: query
	// default: 5
	Limit int `json:"limit"`
}
