.PHONY: run build test clean createdb dropdb resetdb init api-test api-register api-login api-profile api-create-note api-get-notes

# Переменные
APP_NAME=notes-app
BUILD_DIR=./build
DB_NAME=$(shell grep DB_NAME .env | cut -d '=' -f2)
API_URL=http://localhost:8080/api
TOKEN_FILE=.token

# Запуск приложения
run:
	go run cmd/api/main.go

# Сборка приложения
build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) cmd/api/main.go

# Создание базы данных
createdb:
	@echo "Создание базы данных $(DB_NAME)..."
	@if [ -z "$(DB_NAME)" ]; then \
		echo "Ошибка: DB_NAME не найдено в файле .env"; \
		exit 1; \
	fi
	@createdb $(DB_NAME) || echo "База данных $(DB_NAME) уже существует"
	@echo "База данных $(DB_NAME) готова к использованию"

# Удаление базы данных
dropdb:
	@echo "Удаление базы данных $(DB_NAME)..."
	@if [ -z "$(DB_NAME)" ]; then \
		echo "Ошибка: DB_NAME не найдено в файле .env"; \
		exit 1; \
	fi
	@dropdb $(DB_NAME) || echo "База данных $(DB_NAME) не существует"

# Сброс базы данных (удаление и создание заново)
resetdb: dropdb createdb
	@echo "База данных $(DB_NAME) сброшена"

# Запуск тестов
test:
	go test -v ./tests/...

# Запуск тестов с покрытием
test-coverage:
	go test -v -coverprofile=coverage.out ./tests/...
	go tool cover -html=coverage.out -o coverage.html

# Загрузка зависимостей
deps:
	go mod download

# Обновление зависимостей
deps-update:
	go get -u ./...
	go mod tidy

# Очистка
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Запуск линтера
lint:
	golangci-lint run

# Миграция базы данных (запускается автоматически при запуске приложения)
migrate:
	go run cmd/api/main.go --migrate-only

# Инициализация проекта (создание базы данных и загрузка зависимостей)
init: createdb deps
	@echo "Проект инициализирован и готов к использованию"

# Тестирование API с помощью curl

# Регистрация пользователя
api-register:
	@echo "Регистрация пользователя..."
	@curl -s -X POST $(API_URL)/auth/register \
		-H "Content-Type: application/json" \
		-d '{"username":"testuser","email":"test@example.com","password":"password123"}' | jq

# Вход пользователя и сохранение токена
api-login:
	@echo "Вход пользователя..."
	@curl -s -X POST $(API_URL)/auth/login \
		-H "Content-Type: application/json" \
		-d '{"email":"test@example.com","password":"password123"}' | jq
	@curl -s -X POST $(API_URL)/auth/login \
		-H "Content-Type: application/json" \
		-d '{"email":"test@example.com","password":"password123"}' | jq -r '.token' > $(TOKEN_FILE)
	@echo "Токен сохранен в файле $(TOKEN_FILE)"

# Получение профиля пользователя
api-profile:
	@echo "Получение профиля пользователя..."
	@if [ ! -f $(TOKEN_FILE) ]; then \
		echo "Ошибка: файл с токеном не найден. Сначала выполните 'make api-login'"; \
		exit 1; \
	fi
	@curl -s -X GET $(API_URL)/user/profile \
		-H "Authorization: Bearer $$(cat $(TOKEN_FILE))" | jq

# Создание заметки
api-create-note:
	@echo "Создание заметки..."
	@if [ ! -f $(TOKEN_FILE) ]; then \
		echo "Ошибка: файл с токеном не найден. Сначала выполните 'make api-login'"; \
		exit 1; \
	fi
	@curl -s -X POST $(API_URL)/notes \
		-H "Authorization: Bearer $$(cat $(TOKEN_FILE))" \
		-H "Content-Type: application/json" \
		-d '{"title":"Тестовая заметка","content":"Содержимое тестовой заметки"}' | jq

# Получение всех заметок
api-get-notes:
	@echo "Получение всех заметок..."
	@if [ ! -f $(TOKEN_FILE) ]; then \
		echo "Ошибка: файл с токеном не найден. Сначала выполните 'make api-login'"; \
		exit 1; \
	fi
	@curl -s -X GET $(API_URL)/notes \
		-H "Authorization: Bearer $$(cat $(TOKEN_FILE))" | jq

# Тестирование всех API эндпоинтов
api-test: api-register api-login api-profile api-create-note api-get-notes
	@echo "Тестирование API завершено"

# Помощь
help:
	@echo "Доступные команды:"
	@echo "  make run              - Запуск приложения"
	@echo "  make build            - Сборка приложения"
	@echo "  make createdb         - Создание базы данных"
	@echo "  make dropdb           - Удаление базы данных"
	@echo "  make resetdb          - Сброс базы данных (удаление и создание заново)"
	@echo "  make init             - Инициализация проекта (создание базы данных и загрузка зависимостей)"
	@echo "  make test             - Запуск тестов"
	@echo "  make test-coverage    - Запуск тестов с покрытием"
	@echo "  make deps             - Загрузка зависимостей"
	@echo "  make deps-update      - Обновление зависимостей"
	@echo "  make clean            - Очистка"
	@echo "  make lint             - Запуск линтера"
	@echo "  make migrate          - Миграция базы данных"
	@echo "  make api-register     - Тестирование регистрации пользователя"
	@echo "  make api-login        - Тестирование входа пользователя"
	@echo "  make api-profile      - Тестирование получения профиля пользователя"
	@echo "  make api-create-note  - Тестирование создания заметки"
	@echo "  make api-get-notes    - Тестирование получения всех заметок"
	@echo "  make api-test         - Тестирование всех API эндпоинтов"
	@echo "  make help             - Показать эту справку"

# По умолчанию
default: help 