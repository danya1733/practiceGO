package config

import (
	"os"
	"time"
)

// Config представляет конфигурацию приложения
type Config struct {
	HTTP     HTTPConfig
	Database DatabaseConfig
	Log      LogConfig
}

// HTTPConfig содержит настройки HTTP сервера
type HTTPConfig struct {
	Port            string
	ShutdownTimeout time.Duration
}

// DatabaseConfig содержит настройки базы данных
type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// LogConfig содержит настройки логирования
type LogConfig struct {
	Level string
}

// NewConfig создает новую конфигурацию на основе переменных окружения
func NewConfig() (*Config, error) {
	return &Config{
		HTTP: HTTPConfig{
			Port:            getEnv("HTTP_PORT", ":8080"),
			ShutdownTimeout: 30 * time.Second,
		},
		Database: DatabaseConfig{
			URL:             getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/warehouse"),
			MaxOpenConns:    20,
			MaxIdleConns:    5,
			ConnMaxLifetime: 5 * time.Minute,
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}, nil
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
