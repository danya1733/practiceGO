package domain

import (
	"encoding/json"

	"github.com/google/uuid"
)

// Warehouse представляет склад
type Warehouse struct {
	ID      uuid.UUID `json:"id"`
	Address string    `json:"address"`
}

// Product представляет товар
type Product struct {
	ID              uuid.UUID       `json:"id"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	Characteristics json.RawMessage `json:"characteristics"`
	Weight          float64         `json:"weight"`
	Barcode         string          `json:"barcode"`
}

// Inventory представляет связь между товаром и складом
type Inventory struct {
	ID          uuid.UUID `json:"id"`
	WarehouseID uuid.UUID `json:"warehouse_id"`
	ProductID   uuid.UUID `json:"product_id"`
	Quantity    int       `json:"quantity"`
	Price       float64   `json:"price"`
	Discount    float64   `json:"discount"` // в процентах
}

// InventoryWithProduct представляет инвентарь с информацией о товаре
type InventoryWithProduct struct {
	Inventory
	Product Product `json:"product"`
}

// Analytics представляет аналитику продаж
type Analytics struct {
	ID           uuid.UUID `json:"id"`
	WarehouseID  uuid.UUID `json:"warehouse_id"`
	ProductID    uuid.UUID `json:"product_id"`
	SoldQuantity int       `json:"sold_quantity"`
	TotalSum     float64   `json:"total_sum"`
}

// WarehouseAnalytics представляет аналитику по складу
type WarehouseAnalytics struct {
	WarehouseID uuid.UUID `json:"warehouse_id"`
	Address     string    `json:"address"`
	TotalSum    float64   `json:"total_sum"`
}

// ProductPurchase представляет информацию о покупке товара
type ProductPurchase struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}

// PurchaseRequest представляет запрос на покупку товаров
type PurchaseRequest struct {
	WarehouseID uuid.UUID         `json:"warehouse_id"`
	Products    []ProductPurchase `json:"products"`
}

// CalculationResult представляет результат расчета стоимости товаров
type CalculationResult struct {
	TotalSum float64 `json:"total_sum"`
	Items    []struct {
		ProductID         uuid.UUID `json:"product_id"`
		Name              string    `json:"name"`
		Quantity          int       `json:"quantity"`
		Price             float64   `json:"price"`
		PriceWithDiscount float64   `json:"price_with_discount"`
		TotalPrice        float64   `json:"total_price"`
	} `json:"items"`
}
