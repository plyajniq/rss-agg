# Этап сборки
FROM golang:1.23-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Устанавливаем необходимые утилиты, включая make, swag
RUN apk add --no-cache make && \
    go install github.com/swaggo/swag/cmd/swag@latest

# Копируем файлы go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod tidy

# Копируем весь исходный код
COPY . .

# Выполняем сборку с помощью make
RUN make build

# Финальный образ для запуска
FROM alpine:3.18

# Устанавливаем сертификаты для работы с HTTPS
RUN apk --no-cache add ca-certificates

# Копируем скомпилированный бинарник из этапа сборки
COPY --from=builder /app/bin/rssagg /app/rssagg

COPY --from=builder /app/docs /app/docs

COPY --from=builder /app/.env /app/.env

# Устанавливаем команду для запуска приложения
CMD ["/app/rssagg"]
