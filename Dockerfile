# Этап сборки
FROM golang:1.24-alpine AS builder

# Установка необходимых зависимостей
RUN apk add --no-cache git make

# Установка рабочей директории
WORKDIR /app

# Копирование модулей и зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN make build

# Этап финального образа
FROM alpine:latest

# Установка необходимых пакетов
RUN apk add --no-cache ca-certificates tzdata

# Создание непривилегированного пользователя
RUN adduser -D -g '' devhelper

# Установка рабочей директории
WORKDIR /home/devhelper

# Копирование бинарного файла из этапа сборки
COPY --from=builder /app/build/devhelper .

# Установка прав доступа
RUN chown -R devhelper:devhelper .

# Переключение на непривилегированного пользователя
USER devhelper

# Определение точки входа
ENTRYPOINT ["./devhelper"]