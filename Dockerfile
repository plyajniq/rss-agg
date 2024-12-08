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

RUN make build-linux

# Финальный образ для запуска
FROM golang:1.23-alpine

# Устанавливаем сертификаты для работы с HTTPS
RUN apk --no-cache add ca-certificates && \
    go install github.com/pressly/goose/v3/cmd/goose@latest && \
    # Устанавливаем curl для тестирования
    apk add curl


COPY --from=builder /app/bin/rssagg /app/rssagg

COPY --from=builder /app/sql/schema /app/sql/schema

COPY --from=builder /app/static /app/static

COPY --from=builder /app/templates /app/templates

EXPOSE 80

WORKDIR /app

CMD ["/app/rssagg"]
