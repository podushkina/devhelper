.PHONY: build clean test lint install run docker-build docker-run help

# Переменные сборки
BINARY_NAME=devhelper
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "1.0.0")
BUILD_TIME=$(shell date +%FT%T%z)
GIT_COMMIT=$(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}"

# Go команды
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOVET=$(GOCMD) vet
GOLINT=golangci-lint

# Пути
SRC_DIRS=./cmd/... ./internal/... ./pkg/...
CMD_DIR=./cmd/devhelper
BUILD_DIR=./build
DIST_DIR=./dist

# Определение OS и ARCH
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

# Цели
help: ## Показать список доступных команд
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Собрать приложение для текущей платформы
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "Сборка завершена: $(BUILD_DIR)/$(BINARY_NAME)"

run: build ## Собрать и запустить приложение
	$(BUILD_DIR)/$(BINARY_NAME)

clean: ## Очистить артефакты сборки
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -rf $(DIST_DIR)

test: ## Запустить тесты
	$(GOTEST) -v -race -coverprofile=coverage.out $(SRC_DIRS)
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

lint: ## Запустить линтер
	$(GOLINT) run $(SRC_DIRS)

vet: ## Запустить go vet
	$(GOVET) $(SRC_DIRS)

fmt: ## Форматировать код с gofmt
	gofmt -s -w .

tidy: ## Обновить go.mod
	$(GOMOD) tidy

install: build ## Установить приложение в GOPATH/bin
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/

# Кросс-компиляция
build-all: build-linux build-windows build-macos ## Собрать приложение для всех платформ

build-linux: ## Собрать приложение для Linux
	mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 $(CMD_DIR)

build-windows: ## Собрать приложение для Windows
	mkdir -p $(DIST_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_DIR)

build-macos: ## Собрать приложение для macOS
	mkdir -p $(DIST_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_DIR)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 $(CMD_DIR)

# Docker
docker-build: ## Собрать Docker образ
	docker build -t $(BINARY_NAME):$(VERSION) .

docker-run: ## Запустить Docker образ
	docker run --rm -it $(BINARY_NAME):$(VERSION)

# Релиз
release: build-all ## Создать релизный архив
	mkdir -p $(DIST_DIR)/archives
	cd $(DIST_DIR) && tar -czvf archives/$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	cd $(DIST_DIR) && tar -czvf archives/$(BINARY_NAME)-$(VERSION)-linux-arm64.tar.gz $(BINARY_NAME)-linux-arm64
	cd $(DIST_DIR) && tar -czvf archives/$(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	cd $(DIST_DIR) && tar -czvf archives/$(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64
	cd $(DIST_DIR) && zip -r archives/$(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe
	@echo "Релизные архивы созданы в $(DIST_DIR)/archives/"