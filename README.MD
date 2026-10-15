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

Все маршруты `/api/tasks…` требуют заголовок `Authorization: Bearer <token>` и
соответствующее право (`task:*`/`comment:*`). Пагинация: `page ≥ 1`,
`1 ≤ limit ≤ 100` (по умолчанию 10).

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

Сервер поднимается на `http://<BIND_IP>:<PORT>`. При старте fail-fast завершают
ошибки загрузки конфига, подключения к БД (`Open`/`Ping`) и загрузки
`role_permissions`; сбой `ListenAndServe` логируется, но процесс продолжает
работу до сигнала остановки.

### Подключение к БД

> Текущее состояние: создание `*sql.DB`, параметры пула и `Ping` находятся в
> `internal/app/app.go`; вынос в `internal/repository/postgres/connect.go` —
> планируемый рефакторинг.

Инициализацию и настройку подключения к PostgreSQL следует держать в отдельном
файле (например, `internal/repository/postgres/connect.go`): создание `*sql.DB`,
параметры пула (`SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime`) и
проверка `Ping`. В `internal/app/app.go` этот конструктор вызывается, а
полученный `*sql.DB` передаётся в репозитории. Пул открывается один раз и
закрывается тоже один раз — в `App.Close()`, после остановки HTTP-сервера
(не в конструкторе подключения).

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
6. **Логирование** — `log/slog` (JSON в stdout), access-лог `SlogMiddleware`
   (`method`, `path`, `query`, `status`, `latency`, `client_ip`, `body_size`),
   логирование неожиданных сбоев на уровне `error` в хендлерах. Репозитории и
   use-case не логируют — только оборачивают ошибку через `%w`.
7. **Отказоустойчивость** — graceful shutdown по `SIGINT`/`SIGTERM`
   (`internal/app/server.go`); fail-fast при старте: загрузка конфига,
   подключение к БД, загрузка `role_permissions`.

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
   «Конфигурация приложения»), graceful shutdown и fail-fast на старте
   реализованы (`internal/app/server.go`); осталось: CORS, Docker/compose,
   тесты, вынос настройки пула соединений в
   `internal/repository/postgres/connect.go`.
8. **Rate limiting / `429 Too Many Requests`** — ограничение частоты запросов
   (в первую очередь на `/api/auth/login` — защита от перебора), ответ `429`
   с заголовком `Retry-After` и единым телом `{"error": "…"}`; клиент при
   `429` делает повтор с backoff, ориентируясь на `Retry-After`.

---

## Известные слабые места

- `config`: `env-required` не гарантирует непустоту значения, пришедшего из
  `.env`; поэтому `JWT_SECRET` и `DB_PASSWORD` дополнительно проверяются явно
  после загрузки (`internal/config/config.go`).
- `config` реализован как singleton (`sync.Once`) — осознанный выбор.
- `EnsureSchema` (`internal/repository/postgres/schema.go`) нигде не вызывается —
  мёртвый код; схема создаётся вручную, `golang-migrate` не подключён.
- `context.Context` не прокинут в use-case/repository.
- `SlogMiddleware` (`internal/app/server.go`): ветка `c.Errors` не используется
  (никто не вызывает `c.Error`), не проставляется request id, `query`
  логируется без фильтрации.
- Настройка пула соединений находится в `internal/app/app.go`, а не в
  `internal/repository/postgres/connect.go` (см. раздел «Подключение к БД»).
- Нет тестов; Docker/compose не подключены.