# Система управления складами

Система управления складами, товарами и инвентаризацией с аналитикой продаж. Приложение разработано на Go с использованием PostgreSQL в качестве базы данных.

## Функциональность

- Управление складами (создание, получение списка)
- Управление товарами (создание, обновление, получение списка)
- Инвентаризация товаров на складах (добавление товаров, обновление количества, установка скидок)
- Система покупок с учетом скидок
- Аналитика продаж по складам и товарам
- RESTful API для всех операций

## Требования

- Go 1.24 или выше
- PostgreSQL 16
- Docker и Docker Compose (для запуска в контейнерах)

## Установка и запуск

### С использованием Docker (рекомендуется)

1. Клонировать репозиторий:
   ```bash
   git clone https://github.com/danya1733/practiceGO.git
   cd practiceGO
   ```

2. Запустить проект с помощью Docker Compose:
   ```bash
   docker-compose up -d
   ```

   Это автоматически:
   - Создаст базу данных PostgreSQL
   - Применит миграции схем
   - Запустит приложение на порту 8080
   - Запустит pgAdmin 4 для управления базой данных (опционально)

3. Приложение будет доступно по адресу: http://localhost:8080
   pgAdmin 4 будет доступен по адресу: http://localhost:8588 (логин: admin@example.com, пароль: adminpass)

> **Примечание**: Если вам не требуется pgAdmin 4, вы можете удалить соответствующий сервис из файла docker-compose.yml перед запуском.

### Локальный запуск (без Docker)

1. Клонировать репозиторий:
   ```bash
   git clone https://github.com/danya1733/practiceGO.git
   cd practiceGO
   ```

2. Установить и запустить PostgreSQL локально.

3. Создать базу данных:
   ```bash
   createdb warehouse
   ```

4. Применить миграции (требуется установленный [golang-migrate](https://github.com/golang-migrate/migrate)):
   ```bash
   migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/warehouse?sslmode=disable" up
   ```

5. Запустить приложение:
   ```bash
   go run cmd/app/main.go
   ```

## Конфигурация

Конфигурация приложения осуществляется через переменные окружения:

- `DATABASE_URL` - URL подключения к PostgreSQL (по умолчанию: `postgres://postgres:postgres@localhost:5432/warehouse`)
- `HTTP_PORT` - порт для HTTP сервера (по умолчанию: `:8080`)
- `LOG_LEVEL` - уровень логирования (`debug`, `info`, `warn`, `error`) (по умолчанию: `info`)

## API документация

Документация API доступна через Swagger UI по адресу: http://localhost:8080/swagger/

### Основные эндпоинты

#### Склады
- `GET /api/warehouses` - получить список всех складов
- `POST /api/warehouses` - создать новый склад

#### Товары
- `GET /api/products` - получить список всех товаров
- `POST /api/products` - создать новый товар
- `PUT /api/products/{id}` - обновить товар

#### Инвентаризация
- `POST /api/inventory` - создать запись инвентаризации (добавить товар на склад)
- `PUT /api/inventory/quantity` - обновить количество товара на складе
- `PUT /api/inventory/discount` - обновить скидку на товар
- `GET /api/warehouses/{id}/products` - получить список товаров на складе
- `GET /api/warehouses/{warehouse_id}/products/{product_id}` - получить информацию о товаре на складе

#### Покупки
- `POST /api/warehouses/calculate` - рассчитать стоимость покупки с учетом скидок
- `POST /api/warehouses/purchase` - выполнить покупку товаров

#### Аналитика
- `GET /api/analytics/warehouses/{id}` - получить аналитику по складу
- `GET /api/analytics/warehouses/top` - получить топ складов по выручке

## Примеры запросов

### Создание склада

```bash
curl -X POST http://localhost:8080/api/warehouses \
  -H "Content-Type: application/json" \
  -d '{"address": "ул. Складская, 123"}'
```

### Создание товара

```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Ноутбук",
    "description": "Ноутбук Dell XPS 13",
    "characteristics": {"processor": "Intel i7", "ram": "16GB", "storage": "512GB SSD"},
    "weight": 1.3,
    "barcode": "1234567890123"
  }'
```

### Добавление товара на склад

```bash
curl -X POST http://localhost:8080/api/inventory \
  -H "Content-Type: application/json" \
  -d '{
    "warehouse_id": "полученный_uuid_склада",
    "product_id": "полученный_uuid_товара",
    "quantity": 10,
    "price": 75000,
    "discount": 5
  }'
```

### Покупка товаров

```bash
curl -X POST http://localhost:8080/api/warehouses/purchase \
  -H "Content-Type: application/json" \
  -d '{
    "warehouse_id": "полученный_uuid_склада",
    "products": [
      {
        "product_id": "полученный_uuid_товара",
        "quantity": 1
      }
    ]
  }'
```

## Структура базы данных

### warehouses
- `id` - UUID, первичный ключ
- `address` - TEXT, адрес склада

### products
- `id` - UUID, первичный ключ
- `name` - TEXT, название товара
- `description` - TEXT, описание товара
- `characteristics` - JSONB, характеристики товара
- `weight` - FLOAT, вес товара
- `barcode` - TEXT, штрих-код товара (уникальный)

### inventory
- `id` - UUID, первичный ключ
- `warehouse_id` - UUID, внешний ключ на warehouses
- `product_id` - UUID, внешний ключ на products
- `quantity` - INTEGER, количество товара на складе
- `price` - FLOAT, цена товара
- `discount` - FLOAT, скидка на товар в процентах

### analytics
- `id` - UUID, первичный ключ
- `warehouse_id` - UUID, внешний ключ на warehouses
- `product_id` - UUID, внешний ключ на products
- `sold_quantity` - INTEGER, количество проданных товаров
- `total_sum` - FLOAT, общая сумма продаж

## Разработка

### Структура проекта

```
project/
├── cmd/
│   └── app/
│       └── main.go          # Точка входа в приложение
├── internal/
│   ├── app/
│   │   └── app.go           # Инициализация приложения
│   ├── config/
│   │   └── config.go        # Конфигурация приложения
│   ├── domain/
│   │   └── models.go        # Модели данных
│   ├── handler/
│   │   └── handler.go       # HTTP обработчики
│   └── repository/
│       ├── postgres.go      # Подключение к базе данных
│       ├── warehouse_repository.go # Репозиторий для складов
│       ├── product_repository.go   # Репозиторий для товаров
│       ├── inventory_repository.go # Репозиторий для инвентаризации
│       └── analytics_repository.go # Репозиторий для аналитики
├── migrations/
│   ├── 000001_init_schema.up.sql   # Миграция вверх
│   └── 000001_init_schema.down.sql # Миграция вниз
├── pkg/
│   └── logger/
│       └── logger.go        # Пакет для логирования
├── docker-compose.yml       # Docker Compose конфигурация
├── Dockerfile               # Докер-файл для сборки образа
├── go.mod                   # Go модули
└── go.sum                   # Хеши зависимостей
```

### Миграции

Для логирования используется библиотека [zap](https://github.com/uber-go/zap). Все HTTP запросы автоматически логируются с уникальным идентификатором запроса.

## Лицензия

MIT