package repository

import (
	"context"

	"github.com/danya1733/practiceGO/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WarehouseRepository представляет репозиторий для работы со складами
type WarehouseRepository struct {
	pool *pgxpool.Pool
}

// NewWarehouseRepository создает новый репозиторий для работы со складами
func NewWarehouseRepository(pool *pgxpool.Pool) *WarehouseRepository {
	return &WarehouseRepository{pool: pool}
}

// Create создает новый склад
func (r *WarehouseRepository) Create(ctx context.Context, warehouse domain.Warehouse) (domain.Warehouse, error) {
	query := `
		INSERT INTO warehouses (id, address)
		VALUES ($1, $2)
		RETURNING id, address
	`

	if warehouse.ID == uuid.Nil {
		warehouse.ID = uuid.New()
	}

	err := r.pool.QueryRow(ctx, query, warehouse.ID, warehouse.Address).Scan(&warehouse.ID, &warehouse.Address)
	if err != nil {
		return domain.Warehouse{}, err
	}

	return warehouse, nil
}

// GetAll возвращает список всех складов
func (r *WarehouseRepository) GetAll(ctx context.Context) ([]domain.Warehouse, error) {
	query := `
		SELECT id, address
		FROM warehouses
		ORDER BY address
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var warehouses []domain.Warehouse
	for rows.Next() {
		var w domain.Warehouse
		if err := rows.Scan(&w.ID, &w.Address); err != nil {
			return nil, err
		}
		warehouses = append(warehouses, w)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return warehouses, nil
}

// GetByID возвращает склад по его ID
func (r *WarehouseRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Warehouse, error) {
	query := `
		SELECT id, address
		FROM warehouses
		WHERE id = $1
	`

	var warehouse domain.Warehouse
	err := r.pool.QueryRow(ctx, query, id).Scan(&warehouse.ID, &warehouse.Address)
	if err != nil {
		return domain.Warehouse{}, err
	}

	return warehouse, nil
}
