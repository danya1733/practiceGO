package repository

import (
	"context"

	"github.com/danya1733/practiceGO/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProductRepository представляет репозиторий для работы с товарами
type ProductRepository struct {
	pool *pgxpool.Pool
}

// NewProductRepository создает новый репозиторий для работы с товарами
func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

// Create создает новый товар
func (r *ProductRepository) Create(ctx context.Context, product domain.Product) (domain.Product, error) {
	query := `
		INSERT INTO products (id, name, description, characteristics, weight, barcode)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, description, characteristics, weight, barcode
	`

	if product.ID == uuid.Nil {
		product.ID = uuid.New()
	}

	err := r.pool.QueryRow(ctx, query,
		product.ID,
		product.Name,
		product.Description,
		product.Characteristics,
		product.Weight,
		product.Barcode,
	).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Characteristics,
		&product.Weight,
		&product.Barcode,
	)

	if err != nil {
		return domain.Product{}, err
	}

	return product, nil
}

// GetAll возвращает список всех товаров
func (r *ProductRepository) GetAll(ctx context.Context) ([]domain.Product, error) {
	query := `
		SELECT id, name, description, characteristics, weight, barcode
		FROM products
		ORDER BY name
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Characteristics,
			&p.Weight,
			&p.Barcode,
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

// GetByID возвращает товар по его ID
func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	query := `
		SELECT id, name, description, characteristics, weight, barcode
		FROM products
		WHERE id = $1
	`

	var product domain.Product
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Characteristics,
		&product.Weight,
		&product.Barcode,
	)

	if err != nil {
		return domain.Product{}, err
	}

	return product, nil
}

// Update обновляет информацию о товаре
func (r *ProductRepository) Update(ctx context.Context, product domain.Product) (domain.Product, error) {
	query := `
		UPDATE products
		SET name = $2, description = $3, characteristics = $4, weight = $5, barcode = $6
		WHERE id = $1
		RETURNING id, name, description, characteristics, weight, barcode
	`

	err := r.pool.QueryRow(ctx, query,
		product.ID,
		product.Name,
		product.Description,
		product.Characteristics,
		product.Weight,
		product.Barcode,
	).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Characteristics,
		&product.Weight,
		&product.Barcode,
	)

	if err != nil {
		return domain.Product{}, err
	}

	return product, nil
}
