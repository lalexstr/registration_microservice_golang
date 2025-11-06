## Auth Service

Лёгкий сервис аутентификации и управления пользователями на Go (Gin + GORM + SQLite, JWT).

### Возможности
- Регистрация и вход по email/паролю
- Выдача JWT токена (HS256)
- Эндпоинт текущего пользователя `/me`
- CRUD пользователей (с разграничением прав: admin/user)

### Технологии
- Gin (HTTP API)
- GORM + SQLite (хранение данных)
- golang-jwt/jwt (подпись и проверка токенов)
- bcrypt (хэширование паролей)

---

### Быстрый старт
1) Установите Go 1.21+
2) Скопируйте пример переменных окружения и отредактируйте при необходимости:

```bash
cp .env.example .env  # если файла нет — создайте .env вручную, см. ниже
```

Минимально необходимые переменные (значения по умолчанию в скобках):

- `PORT` ("8080") — порт HTTP сервера
- `JWT_SECRET` — секрет для подписи JWT (обязательно)
- `JWT_TTL_MINUTES` ("60") — срок жизни токена в минутах
- `SQLITE_PATH` ("auth.db") — путь к файлу БД SQLite

3) Запуск в dev-режиме:

```bash
go run ./main.go
```

Сервис поднимется на `http://localhost:8080` (или на указанном `PORT`). При первом запуске выполнится авто-миграция и сид админа:

- Email: `admin@example.com`
- Пароль: `admin123`

4) Сборка бинаря:

```bash
go build -o auth-service
./auth-service
```

> Конфиг ищется рядом с бинарём: при старте пытается загрузить `.env` из директории исполняемого файла, затем из текущей директории.

---

### Структура
- `main.go` — точка входа
- `config/` — загрузка конфигурации и переменных окружения
- `db/` — инициализация GORM/SQLite, автмиграции, сид админа
- `models/` — модели БД (User)
- `repositories/` — доступ к данным (UserRepo)
- `services/` — бизнес-логика (AuthService)
- `handlers/` — HTTP-обработчики (Auth/User)
- `middleware/` — JWT аутентификация
- `routes/` — регистрация маршрутов
- `utils/` — утилиты (ответы, хэширование)

---

### API
Базовый префикс: `http://localhost:8080/api/v1`

Тело запросов/ответов — JSON. Защищённые маршруты требуют заголовок `Authorization: Bearer <token>`.

#### Публичные
- POST `/register`
  - Request:
    ```json
    { "email": "user@example.com", "password": "secret123", "full_name": "John Doe" }
    ```
  - Response 201:
    ```json
    { "data": { "id": 1, "email": "user@example.com", "full_name": "John Doe", "role": "user", "created_at": "...", "updated_at": "..." } }
    ```

- POST `/login`
  - Request:
    ```json
    { "email": "user@example.com", "password": "secret123" }
    ```
  - Response 200:
    ```json
    { "data": { "token": "<JWT>" } }
    ```

#### Защищённые
- GET `/me` — профиль текущего пользователя

- GET `/users` — список пользователей (admin)
  - Query: `page` (int, по умолч. 1), `size` (int, по умолч. 20, максимум 100)
  - Response 200:
    ```json
    { "data": { "items": [ {"id":1, "email":"...", "full_name":"...", "role":"..."} ], "total": 1, "page": 1, "size": 20 } }
    ```

- POST `/users` — создать пользователя (admin)
  - Request:
    ```json
    { "email": "new@example.com", "password": "secret123", "full_name": "New User", "role": "user" }
    ```

- GET `/users/:id` — получить пользователя (admin или владелец)

- PUT `/users/:id` — обновить пользователя (admin или владелец)
  - Request (любые поля опциональны):
    ```json
    { "full_name": "New Name", "password": "newpass", "role": "admin" }
    ```
    Примечание: менять `role` может только admin.

- DELETE `/users/:id` — удалить пользователя (admin или владелец)

#### Формат ошибок
```json
{ "error": "message" }
```

---

### Примеры cURL
Регистрация:
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","password":"secret123","full_name":"John"}'
```

Логин (получение токена):
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","password":"secret123"}'
```

Текущий пользователь:
```bash
curl http://localhost:8080/api/v1/me \
  -H "Authorization: Bearer $TOKEN"
```

Создание пользователя админом:
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"email":"new@example.com","password":"secret123","full_name":"New User","role":"user"}'
```

---

### Заметки по безопасности
- Не храните `JWT_SECRET` в репозитории — используйте `.env`
- Пароли хранятся исключительно в виде bcrypt-хэшей
- Убедитесь, что сервис запускается за обратным прокси с HTTPS в проде

---

### Лицензия
MIT (или укажите свою)


