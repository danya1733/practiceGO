package repository

import (
	"context"
	"fmt"

	"github.com/danya1733/practiceGO/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InventoryRepository представляет репозиторий для работы с инвентаризацией
type InventoryRepository struct {
	pool *pgxpool.Pool
}

// NewInventoryRepository создает новый репозиторий для работы с инвентаризацией
func NewInventoryRepository(pool *pgxpool.Pool) *InventoryRepository {
	return &InventoryRepository{pool: pool}
}

// Create создает новую запись инвентаризации
func (r *InventoryRepository) Create(ctx context.Context, inventory domain.Inventory) (domain.Inventory, error) {
	query := `
		INSERT INTO inventory (id, warehouse_id, product_id, quantity, price, discount)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, warehouse_id, product_id, quantity, price, discount
	`

	if inventory.ID == uuid.Nil {
		inventory.ID = uuid.New()
	}

	err := r.pool.QueryRow(ctx, query,
		inventory.ID,
		inventory.WarehouseID,
		inventory.ProductID,
		inventory.Quantity,
		inventory.Price,
		inventory.Discount,
	).Scan(
		&inventory.ID,
		&inventory.WarehouseID,
		&inventory.ProductID,
		&inventory.Quantity,
		&inventory.Price,
		&inventory.Discount,
	)

	if err != nil {
		return domain.Inventory{}, err
	}

	return inventory, nil
}

// GetByWarehouseAndProduct возвращает инвентаризацию по складу и товару
func (r *InventoryRepository) GetByWarehouseAndProduct(ctx context.Context, warehouseID, productID uuid.UUID) (domain.Inventory, error) {
	query := `
		SELECT id, warehouse_id, product_id, quantity, price, discount
		FROM inventory
		WHERE warehouse_id = $1 AND product_id = $2
	`

	var inventory domain.Inventory
	err := r.pool.QueryRow(ctx, query, warehouseID, productID).Scan(
		&inventory.ID,
		&inventory.WarehouseID,
		&inventory.ProductID,
		&inventory.Quantity,
		&inventory.Price,
		&inventory.Discount,
	)

	if err != nil {
		return domain.Inventory{}, err
	}

	return inventory, nil
}

// UpdateQuantity обновляет количество товара на складе
func (r *InventoryRepository) UpdateQuantity(ctx context.Context, warehouseID, productID uuid.UUID, quantity int) (domain.Inventory, error) {
	query := `
		UPDATE inventory
		SET quantity = quantity + $3
		WHERE warehouse_id = $1 AND product_id = $2
		RETURNING id, warehouse_id, product_id, quantity, price, discount
	`

	var inventory domain.Inventory
	err := r.pool.QueryRow(ctx, query, warehouseID, productID, quantity).Scan(
		&inventory.ID,
		&inventory.WarehouseID,
		&inventory.ProductID,
		&inventory.Quantity,
		&inventory.Price,
		&inventory.Discount,
	)

	if err != nil {
		return domain.Inventory{}, err
	}

	return inventory, nil
}

// UpdateDiscount обновляет скидку на товар
func (r *InventoryRepository) UpdateDiscount(ctx context.Context, warehouseID, productID uuid.UUID, discount float64) (domain.Inventory, error) {
	query := `
		UPDATE inventory
		SET discount = $3
		WHERE warehouse_id = $1 AND product_id = $2
		RETURNING id, warehouse_id, product_id, quantity, price, discount
	`

	var inventory domain.Inventory
	err := r.pool.QueryRow(ctx, query, warehouseID, productID, discount).Scan(
		&inventory.ID,
		&inventory.WarehouseID,
		&inventory.ProductID,
		&inventory.Quantity,
		&inventory.Price,
		&inventory.Discount,
	)

	if err != nil {
		return domain.Inventory{}, err
	}

	return inventory, nil
}

// GetProductsByWarehouse возвращает список товаров на складе с пагинацией
func (r *InventoryRepository) GetProductsByWarehouse(ctx context.Context, warehouseID uuid.UUID, page, limit int) ([]domain.InventoryWithProduct, error) {
	query := `
		SELECT i.id, i.warehouse_id, i.product_id, i.quantity, i.price, i.discount,
			   p.id, p.name, p.description, p.characteristics, p.weight, p.barcode
		FROM inventory i
		JOIN products p ON i.product_id = p.id
		WHERE i.warehouse_id = $1
		ORDER BY p.name
		LIMIT $2 OFFSET $3
	`

	offset := (page - 1) * limit
	rows, err := r.pool.Query(ctx, query, warehouseID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.InventoryWithProduct
	for rows.Next() {
		var p domain.InventoryWithProduct
		if err := rows.Scan(
			&p.ID,
			&p.WarehouseID,
			&p.ProductID,
			&p.Quantity,
			&p.Price,
			&p.Discount,
			&p.Product.ID,
			&p.Product.Name,
			&p.Product.Description,
			&p.Product.Characteristics,
			&p.Product.Weight,
			&p.Product.Barcode,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// PurchaseProducts уменьшает количество товаров на складе при покупке
func (r *InventoryRepository) PurchaseProducts(ctx context.Context, warehouseID uuid.UUID, products []domain.ProductPurchase) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Проверяем наличие и достаточное количество каждого товара
	for _, p := range products {
		var currentQuantity int
		err := tx.QueryRow(ctx, `
			SELECT quantity FROM inventory 
			WHERE warehouse_id = $1 AND product_id = $2
		`, warehouseID, p.ProductID).Scan(&currentQuantity)

		if err != nil {
			if err == pgx.ErrNoRows {
				return fmt.Errorf("товар с ID %s не найден на складе %s", p.ProductID, warehouseID)
			}
			return err
		}

		if currentQuantity < p.Quantity {
			return fmt.Errorf("недостаточное количество товара %s на складе: доступно %d, запрошено %d",
				p.ProductID, currentQuantity, p.Quantity)
		}
	}

	// Уменьшаем количество товаров и сохраняем аналитику
	for _, p := range products {
		var price, discount float64

		// Получаем текущую цену и скидку
		err := tx.QueryRow(ctx, `
			SELECT price, discount FROM inventory 
			WHERE warehouse_id = $1 AND product_id = $2
		`, warehouseID, p.ProductID).Scan(&price, &discount)

		if err != nil {
			return err
		}

		// Вычисляем финальную цену с учетом скидки
		finalPrice := price * (1 - discount/100)
		totalSum := finalPrice * float64(p.Quantity)

		// Уменьшаем количество товара
		_, err = tx.Exec(ctx, `
			UPDATE inventory 
			SET quantity = quantity - $3 
			WHERE warehouse_id = $1 AND product_id = $2
		`, warehouseID, p.ProductID, p.Quantity)

		if err != nil {
			return err
		}

		// Записываем аналитику
		_, err = tx.Exec(ctx, `
			INSERT INTO analytics (id, warehouse_id, product_id, sold_quantity, total_sum)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (warehouse_id, product_id) DO UPDATE
			SET sold_quantity = analytics.sold_quantity + $4,
				total_sum = analytics.total_sum + $5
		`, uuid.New(), warehouseID, p.ProductID, p.Quantity, totalSum)

		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
