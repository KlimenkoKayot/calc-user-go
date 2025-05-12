# Распределенный вычислитель арифметических выражений с аутентификацией

Проект представляет собой распределенную систему для вычисления арифметических выражений с JWT аутентификацией и API Gateway. Система использует микросервисную архитектуру с отдельными сервисами:
- API Gateway - единая точка входа
- Auth Service - управление пользователями и JWT токенами
- Вычислительный сервис (Оркестратор + Агенты)

## Содержание
- [Архитектура системы](#архитектура-системы)
  - [API Gateway](#api-gateway)
  - [Сервис аутентификации](#сервис-аутентификации)
  - [Вычислительный сервис](#вычислительный-сервис)
- [Аутентификация и авторизация](#аутентификация-и-авторизация)
- [API Endpoints](#api-endpoints)
  - [Auth Service](#auth-service)
    - [Регистрация](#1-регистрация)
    - [Логин](#2-логин)
  - [Вычислительный сервис](#вычислительный-сервис-1)
    - [Добавление выражения](#1-добавление-выражения)
    - [Получение списка выражений](#2-получение-списка-выражений)
    - [Получение выражения по ID](#3-получение-выражения-по-id)
- [Настройка окружения](#настройка-окружения)
- [Запуск проекта](#запуск-проекта)
- [Примеры запросов](#примеры-запросов)
- [Лицензия](#лицензия)

---

## Архитектура системы

![Архитектура системы](https://github.com/KlimenkoKayot/calc-net-go/blob/main/web/static/img/pattern.jpg)

### API Gateway
- Единая точка входа для всех запросов
- Маршрутизация запросов к соответствующим сервисам
- Проверка JWT токенов
- Балансировка нагрузки

### Сервис аутентификации
- Регистрация и аутентификация пользователей
- Генерация JWT токенов
- Хранение учетных данных пользователей
- Обновление токенов

### Вычислительный сервис
- Оркестратор: прием и распределение задач
- Агенты: выполнение вычислений
- Хранение истории вычислений

---

## Аутентификация и авторизация

Система использует JWT (JSON Web Tokens) для аутентификации. Для доступа к защищенным endpoint'ам необходимо:
1. Зарегистрироваться через `/api/v1/register`
2. Получить токен через `/api/v1/login`
3. Добавлять токен в заголовок запроса:
   ```
   Authorization: Bearer <ваш_токен>
   ```

Токены имеют срок жизни:
- Access Token: 15 минут
- Refresh Token: 24 часа

---

## API Endpoints

# ***Все endpoint доступны через API-шлюз!***

### Auth Service

#### 1. Регистрация
**Endpoint:** `POST /api/v1/register`

**Запрос:**
```json
{
  "login": "username",
  "password": "securepassword"
}
```

**Ответы:**
- `200 OK` - успешная регистрация
- `400 Bad Request` - невалидные данные
- `409 Conflict` - пользователь уже существует

#### 2. Логин
**Endpoint:** `POST /api/v1/login`

**Запрос:**
```json
{
  "login": "username",
  "password": "securepassword"
}
```

**Ответ:**
```json
{
  "access_token": "eyJhbGciOi...",
  "refresh_token": "eyJhbGciOi...",
  "expires_in": 900
}
```

### Вычислительный сервис

Все endpoint'ы требуют `access_token`

#### 1. Добавление выражения
**Endpoint:** `POST /api/v1/calculate`

**Запрос:**
```json
{
  "expression": "(2+2)*5"
}
```

**Ответ:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000"
}
```

#### 2. Получение списка выражений
**Endpoint:** `GET /api/v1/expressions`

**Ответ:**
```json
{
  "expressions": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "expression": "(2+2)*5",
      "status": "completed",
      "result": 20
    }
  ]
}
```

#### 3. Получение выражения по ID
**Endpoint:** `GET /api/v1/expressions/{id}`

**Ответ:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "expression": "(2+2)*5",
  "status": "completed",
  "result": 20
}
```

---

## Настройка окружения

Основные настройки в `config/config.yaml`:

```yaml
api_gateway:
  http:
    host: "127.0.0.1"
    port: 8080
  services:
    auth:
      url: "http://127.0.0.1:8081"
    calc:
      url: "http://127.0.0.1:8082"
  router: "gorilla"
  logger: "zap" # zap | logrus

auth:
  http:
    host: "127.0.0.1"
    port: 8081
    read_timeout: "15s"
    write_timeout: "15s"
  database:
    dsn: "file:data/auth.db?cache=shared&mode=rwc"
  jwt:
    secret: "your-jwt-secret-key"
    access_token_expiry: "15m"
    refresh_token_expiry: "24h"
  logger: "zap" # zap | logrus
  router: "gorilla"

calc:
  orchestrator:
    port: 8082
    time_addition_ms: 0
    time_subtraction_ms: 0
    time_multiplication_ms: 0
    time_division_ms: 0
  agent:
    workers: 4
    timeout: "5s"
```

Переменные окружения (переопределяют yaml):
- `AUTH_JWT_SECRET` - секрет для подписи JWT
- `AUTH_JWT_ACCESS_EXPIRY` - срок жизни access токена
- `AUTH_JWT_REFRESH_EXPIRY` - срок жизни refresh токена

---

## Запуск проекта

1. Клонировать репозиторий:
   ```bash
   git clone https://github.com/KlimenkoKayot/calc-user-go.git
   cd calc-user-go
   ```

2. Запустить сервисы:
   ```bash
   # API Gateway
   go run ./api-gateway/cmd/.
   
   # Auth Service
   go run ./auth/cmd/.
   
   # Оркестратор
   go run ./calc/cmd/orchestrator/.
   
   # Агент
   go run ./calc/cmd/agent/.
   ```

---

## Примеры запросов (api-gateway port 8080)

1. Регистрация:
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"login":"user1", "password":"pass123"}'
```

2. Логин:
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"login":"user1", "password":"pass123"}'
```

3. Добавление выражения (с токеном):
```bash
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"expression": "(2+2)*5"}'
```


4. Получение списка выражений (с токеном):
```bash
curl -X GET http://localhost:8080/api/v1/expressions
```

---

## Лицензия

Этот проект распространяется под лицензией MIT. Подробности см. в файле [LICENSE](LICENSE).