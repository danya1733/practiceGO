FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем только файлы, необходимые для сборки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/app

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарный файл из предыдущего этапа
COPY --from=builder /app/app .

# Копируем директорию с документацией Swagger
COPY --from=builder /app/docs/swagger /root/docs/swagger

EXPOSE 8080

CMD ["./app"]
