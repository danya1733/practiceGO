package repository

import (
	"context"

	"github.com/danya1733/practiceGO/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AnalyticsRepository представляет репозиторий для работы с аналитикой
type AnalyticsRepository struct {
	pool *pgxpool.Pool
}

// NewAnalyticsRepository создает новый репозиторий для работы с аналитикой
func NewAnalyticsRepository(pool *pgxpool.Pool) *AnalyticsRepository {
	return &AnalyticsRepository{pool: pool}
}

// GetWarehouseAnalytics возвращает аналитику по складу
func (r *AnalyticsRepository) GetWarehouseAnalytics(ctx context.Context, warehouseID uuid.UUID) ([]domain.Analytics, float64, error) {
	query := `
		SELECT a.id, a.warehouse_id, a.product_id, a.sold_quantity, a.total_sum,
			   SUM(a.total_sum) OVER() as total_sum
		FROM analytics a
		WHERE a.warehouse_id = $1
		ORDER BY a.total_sum DESC
	`

	rows, err := r.pool.Query(ctx, query, warehouseID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var analytics []domain.Analytics
	var totalSum float64

	for rows.Next() {
		var a domain.Analytics
		if err := rows.Scan(
			&a.ID,
			&a.WarehouseID,
			&a.ProductID,
			&a.SoldQuantity,
			&a.TotalSum,
			&totalSum,
		); err != nil {
			return nil, 0, err
		}
		analytics = append(analytics, a)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return analytics, totalSum, nil
}

// GetTopWarehouses возвращает топ-N складов по выручке
func (r *AnalyticsRepository) GetTopWarehouses(ctx context.Context, limit int) ([]domain.WarehouseAnalytics, error) {
	query := `
		SELECT w.id, w.address, COALESCE(SUM(a.total_sum), 0) as total_sum
		FROM warehouses w
		LEFT JOIN analytics a ON w.id = a.warehouse_id
		GROUP BY w.id, w.address
		ORDER BY total_sum DESC
		LIMIT $1
	`

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var warehouses []domain.WarehouseAnalytics
	for rows.Next() {
		var w domain.WarehouseAnalytics
		if err := rows.Scan(
			&w.WarehouseID,
			&w.Address,
			&w.TotalSum,
		); err != nil {
			return nil, err
		}
		warehouses = append(warehouses, w)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return warehouses, nil
}
