# Chopper

Backend сервис для отслеживания настроения и сна.


## Функционал
 - Регистрация и авторизация (JWT)
 - Создание ежедневных записей
 - Анализ последних 7 дней
 - Alert система
 - Rate limiting
 - Graceful shutdown
 - Dockerized deployment


## Стек технологий
 - Go 1.22+
 - Gin
 - JWT
 - PostgreSQL
 - Docker
 - pgx
 - golang-migrate
 - bcrypt


## Структура
Проект реализован по принципам Clean Architecture

Зависимости направлены внутрь (Dependency Inversion Principle)

Client -> Gin -> Middleware -> Handler -> Service -> Repository -> PostgreSQL

### Архитектура проекта
```
internal/
├── delivery/
│   └── http/
├── usecase/
├── repository/
├── middleware/
├── config/
└── build/
```


## Схема базы данных

### Users
- id (uuid)
- username
- email
- password_hash
- role
- created_at
- deleted_at

### DailyEntries
- id (uuid)
- user_id (uuid)
- date
- mood
- sleep_hours
- load


## Безопасность
 - JWT авторизация
 - Хэширование пароля
 - Rate Limiting - ограничение количества запросов по IP
 - Graceful shutdown с корректным завершением соединений


## API
На всех эндпоинтах используется rate limiter (ограничение запросов по IP)

### POST /users/register
регистрация пользователя

#### Пример запроса
```json
{
    "username": "dexter",
    "email": "tonightsthenight@email.com",
    "password": "bayharbour"
}
```

### POST /users/login
вход и получение токена

#### Пример запроса
```json
{
    "username": "dexter",
    "password": "bayharbour"
}
```

### GET /users/me
получение информации о себе (используется токен аутентификации)

### POST /notes/new
создание записи (используется токен аутентификации)

#### Пример запроса
```json
{
    "mood": 5,
    "sleep_hours": 5.5,
    "load": 5
}
```

### GET /alert/get
получение информации о состоянии (используется токен аутентификации)


## Установка

### Клонировать репозиторий
```bash
git clone https://github.com/HellsKitchen99/Chopper.git
```

### Заменить .env.example на .env
```bash
cp .env.example .env
```

### Запустить контейнеры в Docker
```bash
make rebuild run
```
