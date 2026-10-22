## Проект: **"TaskFlow" — Таск-трекер для команд**

### Описание
Сервис управления задачами с ролевой моделью. Команда может создавать задачи, назначать их на исполнителей, комментировать и отслеживать прогресс.

---

### Основные фичи (MVP)

1. **Аутентификация + роли**
    - JWT или session-based auth
    - Роли: `admin`, `user`, `viewer`
    - Разные права доступа к задачам, комментариям и управлению ролями
      (проверка прав через `role_permissions`)

2. **CRUD задач**
    - Заголовок, описание, приоритет (low/middle/high), статус
      (`to_do`, `in_progress`, `in_review`, `done`, `on_hold`, `cancelled`)
    - Назначение на исполнителя (user_id)
    - Создатель задачи — автосоздание при создании

3. **Комментарии к задачам**
    - Автор, текст, timestamp
    - Логирование изменений

4. **Фильтрация и поиск**
    - По статусу, приоритету, исполнителю
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
POST   /api/v1/auth/register
POST   /api/v1/auth/login
GET    /api/v1/tasks?page=1&limit=20&status=to_do&priority=high&assignee_id=&search=
GET    /api/v1/tasks/:id
POST   /api/v1/tasks
PUT    /api/v1/tasks/:id
DELETE /api/v1/tasks/:id
POST   /api/v1/tasks/:id/comments
GET    /api/v1/tasks/:id/comments
PUT    /api/v1/users/:id/roles
```

Все маршруты `/api/v1/tasks…` требуют заголовок `Authorization: Bearer <token>` и
соответствующее право (`task:*`/`comment:*`). Пагинация: `page ≥ 1`,
`1 ≤ limit ≤ 100` (по умолчанию 10).

Все маршруты `/api/v1/users…` требуют `Authorization: Bearer <token>` и право
`user:update` (есть только у роли `admin`). Тело — `{"role": "viewer|user|admin"}`.
Ответы: `200` — роль обновлена, `403` — нет права, `404` — пользователь не найден,
`422` — неизвестная роль.

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

1. **Аутентификация и роли** — `POST /api/v1/auth/register`, `POST /api/v1/auth/login`,
   JWT (HS256), хеширование паролей (bcrypt), middleware проверки токена.
   Роли `admin`/`user`/`viewer`; claim `role` в токене; проверка разрешений
   (`task:*`, `comment:*`, `user:*`) на маршрутах через `RequirePermission`.
2. **CRUD задач** — создание, список, получение по id, частичное обновление,
   удаление. Слои `handler → use-case → repository → domain`. Есть создатель
   задачи (`author_id`).
3. **Комментарии** — `POST/GET /api/v1/tasks/:id/comments`. Слои
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
8. **CORS** — `gin-contrib/cors`; список origin'ов и методов приходит из конфига
   (`CORS_ALLOWED_ORIGINS`, `CORS_ALLOWED_METHODS`), заголовок `Authorization`
   разрешён в preflight, preflight не расходует лимит запросов.
9. **Rate limiting** — token bucket (`golang.org/x/time/rate`) на IP: отдельный
   строгий лимит на `/api/v1/auth/*` (0.1 rps, burst 5 — защита от перебора),
   на остальных маршрутах — 10 rps (burst 20); ответ `429` с телом
   `{"error": "rate limit exceeded"}`, неактивные записи чистятся по TTL.
10. **Trusted proxies** — `SetTrustedProxies(nil)`: `ClientIP()` берётся из
    соединения, а не из `X-Forwarded-For`, поэтому IP в логах и лимитере не
    подделать. При появлении reverse-proxy понадобится список CIDR.
11. **Управление ролями** — `PUT /api/v1/users/:id/roles` (право `user:update`,
    только `admin`); слои `handler → use-case → repository`. Путь `/roles` —
    множественное число оставлено под будущие операции со списком ролей.
    Смена роли применяется после перелогина: claim `role` в токене из БД не
    перечитывается.

### Осталось (по MVP)

1. **Управление ролями** — базовое назначение роли реализовано
   (`PUT /api/v1/users/:id/roles`, право `user:update`, только `admin`).
   Осталось:
   - запрет менять собственную роль и защита от потери последнего администратора
     (сейчас администратор может снять права с себя);
   - применение смены роли без перелогина — перечитывать роль из БД или
     инвалидировать токены; сейчас роль из claim действует до истечения
     access-токена (`JWT_ACCESS_TTL`, по умолчанию 24h);
   - разделение учётной записи и профиля: вынести профильные данные в отдельную
     таблицу/ресурс (`/users`, `/profile`) вместо одной таблицы `users`,
     которая сейчас хранит логин, пароль и роль.
2. **Миграции БД** — `golang-migrate` не подключён, схема создана вручную; права
   в `role_permissions` (включая `admin → user:update`) заполняются вручную.
   DDL и сид будут добавлены вместе с миграциями.
3. **Комментарии** — реализованы; осталось логирование изменений (audit).
4. **Фильтрация и поиск** — параметры `status`, `priority`, `assignee_id`
   реализованы; остался `search`.
5. **Swagger / OpenAPI** — отсутствует (валидация через
   `go-playground/validator` уже используется).
6. **Инфраструктура** — секреты и параметры БД вынесены в конфиг (см. раздел
   «Конфигурация приложения»), graceful shutdown и fail-fast на старте
   реализованы (`internal/app/server.go`); осталось: Docker/compose,
   тесты, вынос настройки пула соединений в
   `internal/repository/postgres/connect.go`.
7. **Rate limiting / `429 Too Many Requests`** — базовое ограничение частоты
   реализовано (`internal/handlers/middleware/rate_limiter.go`); осталось:
   заголовок `Retry-After` и общий лимит поверх маршрутов (сейчас лимитер
   навешен на каждый маршрут отдельно).

---

## Известные слабые места

- `config`: `env-required` не гарантирует непустоту значения, пришедшего из
  `.env`; поэтому `JWT_SECRET` и `DB_PASSWORD` дополнительно проверяются явно
  после загрузки (`internal/config/config.go`).
- `config` реализован как singleton (`sync.Once`) — осознанный выбор.
- `config.LoadConfig` (`internal/config/config.go`): при отсутствии файла `.env`
  чтение переменных окружения не выполняется — ветка
  `errors.Is(err, os.ErrNotExist)` возвращает ошибку. Локально не проявляется
  (`.env` есть), но в деплое только на переменных окружения сервис не стартует.
  Поправить при подготовке Docker/compose: при `os.ErrNotExist` переходить к
  `cleanenv.ReadEnv`, а не возвращать ошибку.
- `EnsureSchema` (`internal/repository/postgres/schema.go`) нигде не вызывается —
  мёртвый код; схема создаётся вручную, `golang-migrate` не подключён.
- `context.Context` не прокинут в use-case/repository.
- Роль берётся из claim `role` и не перечитывается из БД: изменение роли, включая
  понижение прав, действует после перелогина — до истечения access-токена
  (`JWT_ACCESS_TTL`, по умолчанию 24h).
- `SlogMiddleware` (`internal/app/server.go`): ветка `c.Errors` не используется
  (никто не вызывает `c.Error`), не проставляется request id, `query`
  логируется без фильтрации.
- Настройка пула соединений находится в `internal/app/app.go`, а не в
  `internal/repository/postgres/connect.go` (см. раздел «Подключение к БД»).
- Нет тестов; Docker/compose не подключены.