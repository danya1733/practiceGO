CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Таблица складов
CREATE TABLE IF NOT EXISTS warehouses (
    id UUID PRIMARY KEY,
    address TEXT NOT NULL
);

-- Таблица товаров
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    characteristics JSONB DEFAULT '{}',
    weight FLOAT NOT NULL,
    barcode TEXT NOT NULL UNIQUE
);

-- Таблица инвентаризации (связь между товарами и складами)
CREATE TABLE IF NOT EXISTS inventory (
    id UUID PRIMARY KEY,
    warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    product_id UUID NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL DEFAULT 0,
    price FLOAT NOT NULL,
    discount FLOAT NOT NULL DEFAULT 0,
    UNIQUE (warehouse_id, product_id)
);

-- Таблица аналитики продаж
CREATE TABLE IF NOT EXISTS analytics (
    id UUID PRIMARY KEY,
    warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    product_id UUID NOT NULL REFERENCES products(id),
    sold_quantity INTEGER NOT NULL DEFAULT 0,
    total_sum FLOAT NOT NULL DEFAULT 0,
    UNIQUE (warehouse_id, product_id)
);

-- Индексы для улучшения производительности
CREATE INDEX IF NOT EXISTS idx_inventory_warehouse ON inventory(warehouse_id);
CREATE INDEX IF NOT EXISTS idx_inventory_product ON inventory(product_id);
CREATE INDEX IF NOT EXISTS idx_analytics_warehouse ON analytics(warehouse_id);
CREATE INDEX IF NOT EXISTS idx_analytics_product ON analytics(product_id);
