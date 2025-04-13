package repository

import (
	"context"
	"fmt"

	"github.com/danya1733/practiceGO/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresDB представляет подключение к базе данных PostgreSQL
type PostgresDB struct {
	pool *pgxpool.Pool
}

// NewPostgresDB создает новое подключение к базе данных
func NewPostgresDB(cfg config.DatabaseConfig) (*PostgresDB, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга URL базы данных: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.MaxOpenConns)
	poolConfig.MinConns = int32(cfg.MaxIdleConns)
	poolConfig.MaxConnLifetime = cfg.ConnMaxLifetime

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	// Проверка подключения
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("ошибка проверки подключения к базе данных: %w", err)
	}

	return &PostgresDB{pool: pool}, nil
}

// Close закрывает соединение с базой данных
func (db *PostgresDB) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}

// GetPool возвращает пул соединений
func (db *PostgresDB) GetPool() *pgxpool.Pool {
	return db.pool
}
