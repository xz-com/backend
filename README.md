# Приложение для заметок - Бэкенд

Бэкенд-часть приложения для заметок, написанная на Go с использованием Gin и GORM.

## Особенности

- Аутентификация пользователей с использованием JWT
- CRUD операции для заметок
- Защита маршрутов с помощью middleware
- Работа с базой данных PostgreSQL через GORM
- Тесты для основных компонентов

## Требования

- Go 1.21 или выше
- PostgreSQL
- Переменные окружения (см. файл `.env.example`)

## Установка и запуск

1. Клонируйте репозиторий:

```bash
git clone https://github.com/yourusername/notes-app.git
cd notes-app/backend
```

2. Установите зависимости:

```bash
go mod download
```

3. Создайте базу данных PostgreSQL:

```bash
createdb notes_app
```

4. Настройте переменные окружения:

```bash
cp .env.example .env
# Отредактируйте файл .env с вашими настройками
```

5. Запустите приложение:

```bash
go run cmd/api/main.go
```

## API Endpoints

### Аутентификация

- `POST /api/auth/register` - Регистрация нового пользователя
- `POST /api/auth/login` - Вход пользователя

### Пользователи

- `GET /api/user/profile` - Получение профиля текущего пользователя (требуется JWT)

### Заметки

- `POST /api/notes` - Создание новой заметки (требуется JWT)
- `GET /api/notes` - Получение всех заметок пользователя (требуется JWT)
- `GET /api/notes/:id` - Получение заметки по ID (требуется JWT)
- `PUT /api/notes/:id` - Обновление заметки (требуется JWT)
- `DELETE /api/notes/:id` - Удаление заметки (требуется JWT)

## Запуск тестов

```bash
go test -v ./tests/...
```

## Структура проекта

```
backend/
├── cmd/
│   └── api/
│       └── main.go           # Точка входа в приложение
├── internal/
│   ├── auth/
│   │   └── jwt.go            # Работа с JWT-токенами
│   ├── database/
│   │   └── database.go       # Подключение к базе данных
│   ├── handlers/
│   │   ├── note_handlers.go  # Обработчики для заметок
│   │   └── user_handlers.go  # Обработчики для пользователей
│   ├── middleware/
│   │   └── auth.go           # Middleware для аутентификации
│   ├── models/
│   │   ├── note.go           # Модель заметки
│   │   └── user.go           # Модель пользователя
│   └── routes/
│       └── routes.go         # Настройка маршрутов
├── tests/
│   ├── handlers_test.go      # Тесты для обработчиков
│   └── models_test.go        # Тесты для моделей
├── .env                      # Переменные окружения
├── .env.example              # Пример файла с переменными окружения
├── .gitignore                # Файлы, игнорируемые Git
├── go.mod                    # Зависимости Go
└── README.md                 # Документация
```
