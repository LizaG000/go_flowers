# Flower Shop API

REST API для интернет-магазина цветов, реализованный на Go с использованием Gin, PostgreSQL, Docker Compose, JWT-аутентификации и Swagger-документации.

API предоставляет работу с каталогом цветов, пользователями и избранными цветами.

---

## Содержание

* [Предметная область](#предметная-область)
* [Технологии](#технологии)
* [Функциональные возможности](#функциональные-возможности)
* [Структура проекта](#структура-проекта)
* [Запуск проекта](#запуск-проекта)
* [Миграции базы данных](#миграции-базы-данных)
* [Swagger-документация](#swagger-документация)
* [Аутентификация](#аутентификация)
* [Ограничение частоты запросов](#ограничение-частоты-запросов)
* [Идемпотентность добавления в избранное](#идемпотентность-добавления-в-избранное)
* [Конечные точки API](#конечные-точки-api)
* [Примеры запросов](#примеры-запросов)
* [Выводы](#выводы)

---

## Предметная область

Предметной областью проекта является интернет-магазин цветов.

Система предназначена для хранения информации о цветах, доступных для продажи. Для каждого цветка указываются его название, описание, цена, высота и количество единиц в наличии.

Пользователь может зарегистрироваться, авторизоваться, получить данные своего профиля и формировать список избранных цветов. Избранное используется для сохранения интересующих товаров.

Взаимодействие с системой выполняется через REST API. Клиент передаёт запросы по HTTP, а сервер возвращает ответы в формате JSON.

---

## Технологии

* Go;
* Gin;
* PostgreSQL;
* Docker;
* Docker Compose;
* Goose;
* JWT;
* RSA-ключи;
* bcrypt;
* Swagger UI;
* Swaggo;
* `slog`;
* SHA-256.

---

## Функциональные возможности

API поддерживает:

* получение списка цветов;
* создание цветов;
* обновление данных цветка;
* удаление цветов;
* регистрацию пользователей;
* авторизацию пользователей;
* получение профиля текущего пользователя;
* добавление цветов в избранное;
* получение списка избранных цветов;
* удаление записей из избранного;
* ограничение частоты запросов;
* идемпотентность при добавлении цветка в избранное;
* Swagger-документацию API.

---

## Структура проекта

```text
.
├── cmd/
│   └── server.go
├── config/
│   └── local.yaml
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal/
│   ├── config/
│   ├── controller/
│   ├── dto/
│   ├── entity/
│   ├── middleware/
│   ├── repository/
│   ├── router/
│   ├── security/
│   │   └── keys/
│   ├── service/
│   └── storage/
│       └── postgres/
├── migrations/
├── deploy/
│   ├── Dockerfile
│   └── docker-compose.yml
├── go.mod
└── go.sum
```

Основные уровни приложения:

* `controller` — обработка HTTP-запросов;
* `service` — бизнес-логика;
* `repository` — выполнение запросов к PostgreSQL;
* `entity` — сущности предметной области;
* `dto` — структуры входящих запросов;
* `middleware` — авторизация, логирование и rate limit;
* `security` — JWT, RSA-ключи, bcrypt и SHA-256;
* `router` — регистрация маршрутов.

---

## Запуск проекта

### Требования

Для запуска необходимы:

* Docker;
* Docker Compose;
* Go;
* Goose;
* OpenSSL.

### Создание RSA-ключей

Перед первым запуском необходимо создать пару RSA-ключей.

```bash
mkdir -p internal/security/keys

openssl genrsa -out internal/security/keys/private.pem 2048

openssl rsa \
  -in internal/security/keys/private.pem \
  -pubout \
  -out internal/security/keys/public.pem
```

`private.pem` используется для создания JWT-токенов.

`public.pem` используется для проверки JWT-токенов.

### Конфигурация

Пример файла `config/local.yaml`:

```yaml
env: local

http_server:
  address: "0.0.0.0:8080"

database:
  host: go_postgres
  port: 5432
  user: postgres
  password: postgres
  dbname: go_app_db
  sslmode: disable

auth:
  private_key_path: "./internal/security/keys/private.pem"
  public_key_path: "./internal/security/keys/public.pem"
  token_ttl: 24h
```

### Запуск контейнеров

Из корня проекта:

```bash
docker compose -f deploy/docker-compose.yml up --build
```

Запуск в фоновом режиме:

```bash
docker compose -f deploy/docker-compose.yml up --build -d
```

Просмотр логов:

```bash
docker logs go_gin_app -f
```

Остановка контейнеров:

```bash
docker compose -f deploy/docker-compose.yml down
```

---

## Миграции базы данных

Применение миграций:

```bash
goose -dir ./migrations postgres \
"host=localhost port=5432 user=postgres password=postgres dbname=go_app_db sslmode=disable" up
```

Проверка статуса миграций:

```bash
goose -dir ./migrations postgres \
"host=localhost port=5432 user=postgres password=postgres dbname=go_app_db sslmode=disable" status
```

Создание новой миграции:

```bash
goose -dir ./migrations create migration_name sql
```

---

## Swagger-документация

Swagger UI доступен по адресу:

```text
http://localhost:8080/swagger/index.html
```

Генерация Swagger-документации:

```bash
swag init -g cmd/server.go --parseInternal
```

После изменения Swagger-аннотаций необходимо повторно выполнить эту команду.

JSON-описание API:

```text
http://localhost:8080/swagger/doc.json
```

---

## Аутентификация

Для защищённых маршрутов используется JWT-аутентификация.

После регистрации или авторизации сервер возвращает JWT-токен. Его необходимо передавать в заголовке:

```http
Authorization: Bearer <access_token>
```

Пример:

```http
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Обоснование выбора JWT

JWT выбран по следующим причинам:

* не требуется хранить пользовательские сессии на сервере;
* токен содержит идентификатор пользователя;
* токен удобно передавать через HTTP-заголовок;
* JWT подходит для REST API;
* RSA позволяет подписывать токен приватным ключом и проверять публичным;
* пароль не требуется передавать при каждом запросе.

Пароли пользователей не хранятся в открытом виде. Перед сохранением они хешируются с помощью bcrypt.

---

## Ограничение частоты запросов

Для API реализовано ограничение количества запросов с одного IP-адреса.

При превышении лимита сервер возвращает:

```text
429 Too Many Requests
```

Также в ответ передаются заголовки:

```http
X-Limit-Remaining: 0
Retry-After: 35
```

* `X-Limit-Remaining` — число оставшихся запросов;
* `Retry-After` — время в секундах до следующей доступной попытки.

Пример ответа:

```json
{
  "error": "слишком много запросов, попробуйте позже"
}
```

---

## Идемпотентность добавления в избранное

Ключ идемпотентности используется только для добавления цветка в избранное:

```text
POST /api/favorites
```

Ключ передаётся в заголовке:

```http
Idempotency-Key: 7e9d3b4a-6f84-4d1a-9c52-3c7e4b8f1d20
```

При повторной отправке одинакового запроса с тем же ключом сервер не создаёт вторую запись в таблице избранного, а возвращает сохранённый результат первого запроса.

Для проверки тела повторного запроса используется SHA-256-хеш.

Если с тем же ключом отправлено другое тело запроса, сервер возвращает:

```json
{
  "error": "ключ идемпотентности уже использован для другого запроса"
}
```

---

## Конечные точки API

Базовый путь:

```text
/api
```

### Цветы

| Метод  | URL            | Описание               |
| ------ | -------------- | ---------------------- |
| GET    | `/api/flowers` | Получить список цветов |
| POST   | `/api/flowers` | Создать цветок         |
| PUT    | `/api/flowers` | Обновить цветок        |
| DELETE | `/api/flowers` | Удалить цветок         |

### Аутентификация

| Метод | URL                      | Описание                 |
| ----- | ------------------------ | ------------------------ |
| POST  | `/api/auth/registration` | Регистрация пользователя |
| POST  | `/api/auth/login`        | Авторизация пользователя |

### Пользователь

| Метод | URL          | Описание                               |
| ----- | ------------ | -------------------------------------- |
| GET   | `/api/users` | Получить профиль текущего пользователя |

Для маршрута `/api/users` требуется JWT.

### Избранное

| Метод  | URL                                | Описание                     |
| ------ | ---------------------------------- | ---------------------------- |
| GET    | `/api/v1/favorites`                | Получить избранное, версия 1 |
| GET    | `/api/v2/favorites`                | Получить избранное, версия 2 |
| POST   | `/api/favorites`                   | Добавить цветок в избранное  |
| DELETE | `/api/favorites?favoriteID={uuid}` | Удалить запись из избранного |

Для всех маршрутов избранного требуется JWT-токен.

Для `POST /api/favorites` дополнительно требуется `Idempotency-Key`.

---

## Примеры запросов

### Получение списка цветов

```http
GET /api/flowers
```

Пример ответа:

```json
[
  {
    "id": "30db86b6-bd4e-4d50-b1d9-55d0803ce96c",
    "title": "Роза красная",
    "description": "Свежая красная роза на длинном стебле.",
    "price": 350,
    "height": 60,
    "count": 1,
    "created_at": "2026-06-24T10:00:00Z",
    "updated_at": "2026-06-24T10:00:00Z"
  }
]
```

### Создание цветка

```http
POST /api/flowers
Content-Type: application/json
```

```json
{
  "title": "Эустома сиреневая",
  "description": "Изящная сиреневая эустома с несколькими бутонами на стебле.",
  "price": 380,
  "height": 58,
  "count": 5
}
```

Пример ответа:

```json
{
  "id": "0e2e2ed6-a7a2-40a8-9b15-1e469f3fc2f0",
  "title": "Эустома сиреневая",
  "description": "Изящная сиреневая эустома с несколькими бутонами на стебле.",
  "price": 380,
  "height": 58,
  "count": 5,
  "created_at": "2026-06-24T10:20:00Z",
  "updated_at": "2026-06-24T10:20:00Z"
}
```

### Обновление цветка

```http
PUT /api/flowers
Content-Type: application/json
```

```json
{
  "id": "30db86b6-bd4e-4d50-b1d9-55d0803ce96c",
  "update_data": {
    "price": 450,
    "count": 7
  }
}
```

Пример ответа:

```json
{
  "id": "30db86b6-bd4e-4d50-b1d9-55d0803ce96c",
  "title": "Роза красная",
  "description": "Свежая красная роза на длинном стебле.",
  "price": 450,
  "height": 60,
  "count": 7,
  "created_at": "2026-06-24T10:00:00Z",
  "updated_at": "2026-06-24T11:00:00Z"
}
```

### Удаление цветка

```http
DELETE /api/flowers
Content-Type: application/json
```

```json
{
  "id": "30db86b6-bd4e-4d50-b1d9-55d0803ce96c"
}
```

### Регистрация пользователя

```http
POST /api/auth/registration
Content-Type: application/json
```

```json
{
  "first_name": "Елизавета",
  "second_name": "Александровна",
  "last_name": "Иванова",
  "email": "elizaveta@example.com",
  "birth_date": "2004-05-18T00:00:00Z",
  "password": "MyStrongPassword123"
}
```

Пример ответа:

```json
{
  "user": {
    "id": "8c572e51-61d9-4476-b192-17e5a8b1040e",
    "first_name": "Елизавета",
    "second_name": "Александровна",
    "last_name": "Иванова",
    "email": "elizaveta@example.com",
    "birth_date": "2004-05-18T00:00:00Z",
    "created_at": "2026-06-24T10:00:00Z",
    "updated_at": "2026-06-24T10:00:00Z"
  },
  "access_token": "jwt_token"
}
```

### Авторизация

```http
POST /api/auth/login
Content-Type: application/json
```

```json
{
  "email": "elizaveta@example.com",
  "password": "MyStrongPassword123"
}
```

Пример ответа:

```json
"jwt_token"
```

### Получение профиля

```http
GET /api/users
Authorization: Bearer <access_token>
```

Пример ответа:

```json
{
  "id": "8c572e51-61d9-4476-b192-17e5a8b1040e",
  "first_name": "Елизавета",
  "second_name": "Александровна",
  "last_name": "Иванова",
  "email": "elizaveta@example.com",
  "birth_date": "2004-05-18T00:00:00Z",
  "created_at": "2026-06-24T10:00:00Z",
  "updated_at": "2026-06-24T10:00:00Z"
}
```

### Добавление цветка в избранное

```http
POST /api/favorites
Authorization: Bearer <access_token>
Idempotency-Key: 7e9d3b4a-6f84-4d1a-9c52-3c7e4b8f1d20
Content-Type: application/json
```

```json
{
  "flowerID": "30db86b6-bd4e-4d50-b1d9-55d0803ce96c"
}
```

Пример ответа:

```json
{
  "id": "2c93cb7b-27e1-448e-8894-bb1c9ba9b3fb",
  "user_id": "8c572e51-61d9-4476-b192-17e5a8b1040e",
  "flower_id": "30db86b6-bd4e-4d50-b1d9-55d0803ce96c",
  "created_at": "2026-06-24T10:30:00Z",
  "updated_at": "2026-06-24T10:30:00Z"
}
```

### Получение избранных цветов

```http
GET /api/v2/favorites
Authorization: Bearer <access_token>
```

Пример ответа:

```json
[
  {
    "id": "2c93cb7b-27e1-448e-8894-bb1c9ba9b3fb",
    "user_id": "8c572e51-61d9-4476-b192-17e5a8b1040e",
    "flower_id": "30db86b6-bd4e-4d50-b1d9-55d0803ce96c",
    "title": "Роза красная",
    "description": "Свежая красная роза на длинном стебле.",
    "price": 350,
    "height": 60,
    "count": 1,
    "created_at": "2026-06-24T10:30:00Z",
    "updated_at": "2026-06-24T10:30:00Z"
  }
]
```

### Удаление записи из избранного

```http
DELETE /api/favorites?favoriteID=2c93cb7b-27e1-448e-8894-bb1c9ba9b3fb
Authorization: Bearer <access_token>
```

Пример ответа:

```json
{
  "id": "2c93cb7b-27e1-448e-8894-bb1c9ba9b3fb",
  "user_id": "8c572e51-61d9-4476-b192-17e5a8b1040e",
  "flower_id": "30db86b6-bd4e-4d50-b1d9-55d0803ce96c",
  "created_at": "2026-06-24T10:30:00Z",
  "updated_at": "2026-06-24T10:30:00Z"
}
```

---

## Выводы

В результате был разработан REST API интернет-магазина цветов.

Приложение разделено на контроллеры, сервисы и репозитории, что позволяет отделить HTTP-логику, бизнес-правила и работу с базой данных.

Для аутентификации используются JWT-токены с RSA-подписью. Пароли пользователей хранятся в виде bcrypt-хешей.

В проекте реализованы Swagger-документация, ограничение частоты запросов и идемпотентность операции добавления цветка в избранное. Это повышает надёжность API и предотвращает повторное создание одинаковых записей избранного.
