## Проект: **"TaskFlow" — Таск-трекер для команд**

### Описание
Сервис управления задачами с ролевой моделью. Команда может создавать задачи, назначать их на исполнителей, комментировать и отслеживать прогресс.

---

### Основные фичи (MVP)

1. **Аутентификация + роли**
    - JWT или session-based auth
    - Роли: `admin`, `user`, `viewer`
    - Разные права доступа к CRUD задачам (проверка прав через `role_permissions`)

2. **CRUD задач**
    - Заголовок, описание, приоритет (low/medium/high), статус (todo/in-progress/done)
    - Назначение на исполнителя (user_id)
    - Создатель задачи — автосоздание при создании

3. **Комментарии к задачам**
    - Автор, текст, timestamp
    - Логирование изменений

4. **Фильтрация и поиск**
    - По статусу, приоритету, исполнителю, тегам (опционально)
    - Поиск по тексту задачи

5. **REST API + Swagger/OpenAPI docs**
    - Автоматическая документация
    - Валидация запросов

---

### Стек технологий (рекомендация)

| Компонент | Технология |
|-----------|------------|
| Language/Framework | Go + Gin |
| Database | PostgreSQL |
| Auth | JWT (golang-jwt) |
| Validation | go-playground/validator |
| Docs | Swag / gin-swagger |
| Migrations | golang-migrate |

---

### Структура API (примерные эндпоинты)

```
POST   /api/auth/register
POST   /api/auth/login
GET    /api/tasks?page=1&limit=20&status=todo&priority=high&assignee_id=&search=
GET    /api/tasks/:id
PUT    /api/tasks/:id
DELETE /api/tasks/:id
POST   /api/tasks/:id/comments
GET    /api/tasks/:id/comments
```

---

### Что это даёт тебе

- **БД**: PostgreSQL, миграции, транзакции
- **Auth**: JWT, middleware для ролей
- **API design**: RESTful conventions, pagination, error handling
- **Validation**: input validation, custom errors
- **Практика**: полная реализация с нуля

---

## Конфигурация приложения

Конфигурация читается пакетом `internal/config` из переменных окружения и, для
локального запуска, из файла `.env` (библиотека `cleanenv`). Файл `.env` в
репозиторий не коммитится; шаблон — `.env.example`.

### Переменные окружения

| Переменная | Обязательна | По умолчанию | Назначение |
|---|---|---|---|
| `BIND_IP` | нет | `0.0.0.0` | адрес, который слушает сервер |
| `PORT` | нет | `8080` | порт HTTP-сервера |
| `JWT_SECRET` | да | — | ключ подписи JWT (HS256) |
| `JWT_ACCESS_TTL` | нет | `24h` | время жизни access-токена |
| `JWT_REFRESH_TTL` | нет | `168h` | время жизни refresh-токена |
| `DB_HOST` | нет | `localhost` | хост PostgreSQL |
| `DB_PORT` | нет | `5432` | порт PostgreSQL |
| `DB_NAME` | нет | `taskflow` | имя БД |
| `DB_USER` | нет | `postgres` | пользователь БД |
| `DB_PASSWORD` | да | — | пароль БД |
| `DB_SSL_MODE` | нет | `disable` | режим SSL (`disable`/`require`/`verify-full`) |

### Запуск локально

1. Скопировать шаблон: `cp .env.example .env`
2. Заполнить `JWT_SECRET` и `DB_PASSWORD` (остальное — при необходимости).
3. Запустить: `go run .`

Сервер поднимается на `http://<BIND_IP>:<PORT>`. При старте проверяются конфиг,
подключение к БД и запуск сервера — ошибка любого шага завершает процесс
(fail-fast).

---

## Текущее состояние

### Реализовано

1. **Аутентификация и роли** — `POST /api/auth/register`, `POST /api/auth/login`,
   JWT (HS256), хеширование паролей (bcrypt), middleware проверки токена.
   Роли `admin`/`user`/`viewer`; claim `role` в токене; проверка разрешений
   (`task:*`, `comment:*`) на маршрутах через `RequirePermission`.
2. **CRUD задач** — создание, список, получение по id, частичное обновление,
   удаление. Слои `handler → use-case → repository → domain`. Есть создатель
   задачи (`author_id`).
3. **Комментарии** — `POST/GET /api/tasks/:id/comments`. Слои
   `handler → use-case → repository → domain`.
4. **Роли и права** — `domain.Role`/`domain.Permission`, таблицы `roles` и
   `role_permissions`, middleware `RequirePermission`.
5. **Инфраструктура** — слои repository/usecase, пакет `config`, `.gitignore`,
   git-репозиторий.

### Осталось (по MVP)

1. **Управление ролями** — назначение ролей пользователям через API не
   реализовано: при регистрации роль по умолчанию `user`; `admin` назначается
   вне API. Модель ролей и проверка прав уже реализованы.
2. **Расширение модели задачи** — `priority`, `status`, `assignee_id`, теги.
   Создатель уже реализован (`author_id`).
3. **Миграции БД** — `golang-migrate` не подключён, схема создана вручную.
4. **Комментарии** — реализованы; осталось логирование изменений (audit).
5. **Фильтрация и поиск** — параметры `status`, `priority`, `assignee_id`,
   `search` пока не реализованы; пагинация уже исправлена.
6. **Swagger / OpenAPI** — отсутствует (валидация через
   `go-playground/validator` уже используется).
7. **Инфраструктура** — секреты и параметры БД вынесены в конфиг (см. раздел
   «Конфигурация приложения»); осталось: graceful shutdown, CORS,
   Docker/compose, тесты.

---

## Известные слабые места

- `config`: `env-required` не гарантирует непустоту значения, пришедшего из
  `.env`; обязательные секреты требуют явной проверки после загрузки.
- `config` реализован как singleton (`sync.Once`) — осознанный выбор.
- `role_permissions` загружаются не фатально: при сбое checker пуст, и все
  защищённые маршруты отдают 403.
- `context.Context` не прокинут в use-case/repository.
- Логирование через `fmt`/`log`; переход на `log/slog` — отдельная задача.
- Нет тестов; Docker/compose и graceful shutdown не подключены.