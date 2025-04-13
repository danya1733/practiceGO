// Package main Warehouse Management System API
//
// # API для системы управления складами, товарами и инвентаризацией с аналитикой продаж
//
// Schemes: http
// Host: localhost:8080
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
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danya1733/practiceGO/internal/app"
	"github.com/danya1733/practiceGO/internal/config"
	"github.com/danya1733/practiceGO/pkg/logger"
)

// @title Warehouse Management System API
// @version 1.0
// @description API для системы управления складами, товарами и инвентаризацией с аналитикой продаж

// @contact.name API Support
// @contact.email support@example.com

// @host localhost:8080
// @BasePath /api

func main() {
	// Инициализация конфигурации
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Инициализация логгера
	l, err := logger.NewLogger(cfg.Log.Level)
	if err != nil {
		log.Fatalf("Ошибка инициализации логгера: %v", err)
	}
	defer l.Sync()

	// Создание и запуск приложения
	application, err := app.NewApp(cfg, l)
	if err != nil {
		l.Fatal("Ошибка создания приложения", logger.Error(err))
	}

	// Запуск HTTP сервера
	server := &http.Server{
		Addr:    cfg.HTTP.Port,
		Handler: application.Router(),
	}

	// Канал для сигналов завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Запуск сервера в горутине
	go func() {
		l.Info("Запуск HTTP сервера", logger.String("address", cfg.HTTP.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatal("Ошибка запуска HTTP сервера", logger.Error(err))
		}
	}()

	// Ожидание сигнала завершения
	<-quit
	l.Info("Получен сигнал завершения, начинаю graceful shutdown...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		l.Fatal("Ошибка при graceful shutdown", logger.Error(err))
	}

	// Закрытие соединения с базой данных
	if err := application.Close(); err != nil {
		l.Fatal("Ошибка при закрытии приложения", logger.Error(err))
	}

	l.Info("Сервер успешно остановлен")
}
