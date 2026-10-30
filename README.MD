## Проект: **"TaskFlow" — Таск-трекер для команд**

### Описание
Сервис управления задачами с ролевой моделью. Команда может создавать задачи, назначать их на исполнителей, комментировать и отслеживать прогресс.

---

### Основные фичи (MVP)

1. **Аутентификация + роли**
    - JWT (access + refresh), обновление токенов и logout с отзывом
    - Роли: `admin`, `user`, `viewer`
    - Разные права доступа к задачам, комментариям и назначению ролей
      (проверка прав через `role_permissions`)

2. **CRUD задач**
    - Заголовок, описание, приоритет (low/middle/high), статус
      (`to_do`, `in_progress`, `in_review`, `done`, `on_hold`, `cancelled`)
    - Назначение на исполнителя (user_id)
    - Создатель задачи — автосоздание при создании

3. **Комментарии к задачам**
    - Автор, текст, timestamp
    - Создание и чтение (`POST/GET /api/v1/tasks/:id/comments`)

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

### Структура API (эндпоинты)

```
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh
POST   /api/v1/auth/logout
GET    /api/v1/tasks?page=1&limit=20&status=to_do&priority=high&assignee_id=&search=
GET    /api/v1/tasks/:id
POST   /api/v1/tasks
PUT    /api/v1/tasks/:id
DELETE /api/v1/tasks/:id
POST   /api/v1/tasks/:id/comments
GET    /api/v1/tasks/:id/comments
GET    /api/v1/profiles/:id
PUT    /api/v1/profiles/:id
GET    /api/v1/users/roles
PUT    /api/v1/users/:id/roles

GET    /swagger/index.html
GET    /swagger/doc.json
```

Swagger UI — `/swagger/index.html`, спецификация — `/swagger/doc.json`; оба
маршрута находятся вне `/api/v1` и не требуют токена (см. «Swagger / OpenAPI»).

Все маршруты `/api/v1/tasks…` требуют заголовок `Authorization: Bearer <token>` и
соответствующее право (`task:*`/`comment:*`). Пагинация: `page ≥ 1`,
`1 ≤ limit ≤ 100` (по умолчанию 10).

`GET /api/v1/tasks` поддерживает фильтры `status`, `priority`, `assignee_id` и
`search` — регистронезависимый поиск по `name` и `description` (до 100 символов;
`%`, `_`, `\` трактуются буквально). Пустой `search` фильтр не применяет.

Все маршруты `/api/v1/users…` требуют `Authorization: Bearer <token>`.
`GET /api/v1/users/roles` — каталог ролей, право `user:read` (есть только у
`admin`), ответ `200 {"roles": ["admin", "user", "viewer"]}` (сортировка по имени).
`PUT /api/v1/users/:id/roles` — назначение роли пользователю, право `user:update`
(только `admin`); тело — `{"role": "viewer|user|admin"}`. Ответы: `200` — роль
обновлена, `403` — нет права или попытка изменить собственную роль, `404` —
пользователь не найден, `409` — попытка понизить последнего администратора,
`422` — неизвестная роль. Проверка и обновление выполняются в одной транзакции с
блокировкой строк администраторов (`FOR NO KEY UPDATE`), поэтому конкурентные
понижения не оставляют систему без администратора. Смена роли применяется
немедленно: роль читается из БД на каждом защищённом запросе, claim `role` в
токене не используется и не выпускается — клиент не должен вычислять права из
payload токена.

`POST /api/v1/auth/register` принимает `{"email": "…", "password": "…"}` и
создаёт аккаунт (`users`) и пустой профиль (`profiles`) в одной транзакции;
ответ — `201 {"id": "…"}`. `POST /api/v1/auth/login` возвращает
`{"id": "…", "access_token": "…", "refresh_token": "…"}`; `id` можно использовать
для запроса собственного профиля.

Оба токена — JWT (HS256) с claims `id`, `ver`, `typ` и `exp`; роль в токены не
кладётся. Access-токен (`typ="access"`) живёт `JWT_ACCESS_TTL` (по умолчанию
24h), refresh (`typ="refresh"`) — `JWT_REFRESH_TTL` (по умолчанию 30 дней).
`ver` — текущий `users.token_version`, который middleware сверяет с БД на каждом
защищённом запросе.

`POST /api/v1/auth/refresh` принимает `{"refresh_token": "…"}` и возвращает новую
пару `{"access_token": "…", "refresh_token": "…"}`. Refresh-токен проверяется по
подписи, `typ`, `exp` и `ver`; серверного хранилища refresh-токенов нет, поэтому
он не ротируется и живёт до `exp` или до logout. Невалидный или просроченный
токен → `401`.

`POST /api/v1/auth/logout` требует `Authorization: Bearer <access_token>` и
увеличивает `users.token_version`; все ранее выданные токены пользователя — и
access, и refresh — становятся невалидными сразу, следующий логин получает
токены с новой версией. Logout **глобальный**: завершает сессию на всех
устройствах. Ответ — `200 {"message": "logged out"}`.

Все маршруты `/api/v1/profiles…` требуют `Authorization: Bearer <token>`.
`GET /api/v1/profiles/:id` — профиль пользователя, право `profile:read` (есть у
всех ролей); ответ — `200 {"profile": {"user_id": "…", "username": null,
"first_name": null, "last_name": null, "created_at": "…", "updated_at": "…"}}`.
`PUT /api/v1/profiles/:id` — частичное обновление профиля, право
`profile:update` (есть у всех ролей), но только собственного: `id` должен
совпадать с id из токена, иначе `403`. Поля `username`, `first_name`,
`last_name` — от 2 до 100 символов после `TrimSpace`; пустая или пробельная
строка → `422`. Ответ — `200 {"message": "profile was updated"}`. В
`profiles` поля `username`, `first_name`, `last_name` nullable: при регистрации
профиль создаётся пустым, `username` заполняется позже через `PUT`.

---

### Что это даёт тебе

- **БД**: PostgreSQL, миграции, транзакции
- **Auth**: JWT (access/refresh), middleware ролей, отзыв токенов через `token_version`
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
| `JWT_ACCESS_TTL` | нет | `24h` | время жизни access-токена; logout отзывает его досрочно, на актуальность роли не влияет |
| `JWT_REFRESH_TTL` | нет | `720h` | время жизни refresh-токена (30 дней); logout отзывает его досрочно |
| `DB_HOST` | нет | `localhost` | хост PostgreSQL |
| `DB_PORT` | нет | `5432` | порт PostgreSQL |
| `DB_NAME` | нет | `taskflow` | имя БД |
| `DB_USER` | нет | `postgres` | пользователь БД |
| `DB_PASSWORD` | да | — | пароль БД |
| `DB_SSL_MODE` | нет | `disable` | режим SSL (`disable`/`require`/`verify-full`) |
| `DB_MAX_OPEN_CONNS` | нет | `25` | максимальное число открытых соединений в пуле |
| `DB_MAX_IDLE_CONNS` | нет | `2` | максимальное число простаивающих соединений в пуле |

### Запуск локально

1. Создать пустую БД (по умолчанию `taskflow`): `CREATE DATABASE taskflow;`
2. Скопировать шаблон: `cp .env.example .env`
3. Заполнить `JWT_SECRET` и `DB_PASSWORD` (остальное — при необходимости).
4. Запустить: `go run .` — схема и справочники создаются миграциями при старте,
   отдельный шаг с `migrate` не нужен.
5. Зарегистрировать пользователя и назначить первого администратора SQL-ом
   (см. «Миграции БД»): API назначения ролей требует `user:update`, которое есть
   только у `admin`.

Сервер поднимается на `http://<BIND_IP>:<PORT>`. При старте fail-fast завершают
ошибки загрузки конфига, подключения к БД (`Open`/`Ping`), применения миграций и
загрузки `role_permissions`; сбой `ListenAndServe` логируется, но процесс
продолжает работу до сигнала остановки.

### Подключение к БД

Создание и настройка подключения к PostgreSQL вынесены в
`internal/repository/postgres/connect.go` (`ConnectDB`): DSN собирается из
переменных `DB_*`, открывается `*sql.DB`, задаются параметры пула
(`SetMaxOpenConns` — `DB_MAX_OPEN_CONNS`, `SetMaxIdleConns` —
`DB_MAX_IDLE_CONNS`, `SetConnMaxLifetime` — 1 минута), выполняется `Ping` и
применяются миграции
(`RunMigrations`). `internal/app/app.go` только вызывает конструктор и передаёт
полученный `*sql.DB` в репозитории; пул открывается один раз и закрывается тоже
один раз — в `App.Close()`, после остановки HTTP-сервера.

### Миграции БД

Схема БД описана миграциями `golang-migrate`; ручное создание таблиц не
используется. SQL-файлы лежат в `migrations/` и встраиваются в бинарник через
`//go:embed` (`migrations/embed.go`), поэтому каталог с миграциями рядом с
бинарником не нужен. Применение — при старте приложения: `postgres.ConnectDB`
(`internal/repository/postgres/connect.go`) вызывает `postgres.RunMigrations`
(`internal/repository/postgres/migrate.go`), а тот — `Up()`; история версий
хранится в таблице `schema_migrations`.

| Версия | Содержимое |
|---|---|
| `000001` | `users` (роль по умолчанию — `user`) |
| `000002` | каталог `roles` |
| `000003` | `role_permissions` (M:N, FK на `roles`) |
| `000004` | базовый сид ролей и прав: `admin` — 8, `user` — 6, `viewer` — 2 |
| `000005` | `tasks` (FK на `users`, индексы по `author_id`/`assignee_id`) |
| `000006` | `task_comments` (FK на `tasks` и `users`) |
| `000007` | `profiles` (1:1 с `users`, FK с `ON DELETE CASCADE`, `username` nullable); `users.username` переносится в `profiles` и удаляется |
| `000008` | сид `profile:read`/`profile:update` для `admin`, `user`, `viewer`; итоговая матрица — `admin` 10, `user` 8, `viewer` 4 |
| `000009` | `users.token_version` (`INTEGER NOT NULL DEFAULT 0`) — версия токенов для мгновенного logout |

Правила:

- **Применённую миграцию не редактируют.** Изменения схемы — новой версией
  (`000010_*.up.sql`). `golang-migrate` не хранит контрольных сумм, поэтому
  правка старого файла на уже мигрированной БД не заметна, а на чистой даст
  другую схему.
- **Пара `up`/`down` обязательна.** Для данных `down` удаляет только то, что
  вставил `up` (в `000004` — строки трёх ролей, а не весь справочник).
- **`IF NOT EXISTS` в `CREATE` не используется**: расхождение схемы должно падать
  явно (`42P07 relation already exists`), а не маскироваться.
- **Целостность — в БД:** FK с осмысленным `ON DELETE` (`CASCADE` для
  задач/комментариев, `SET NULL` для исполнителя) и индексы под фактические
  фильтры. `CHECK` на статусы, приоритеты и роли не ставим — это дублировало бы
  валидацию `IsValid()` из `domain` (см. «Известные слабые места»).
- **Чистая установка** — пустая БД + запуск приложения.
- **Существующая БД, созданная до миграций**, автоматически не принимается:
  либо пересоздать схему (`DROP SCHEMA public CASCADE; CREATE SCHEMA public;`),
  либо принять её как baseline (`pg_dump --schema-only` → одна baseline-миграция
  → `migrate force <version>`). Признак непринятой БД — `42P07` на первой
  миграции; признак упавшей — `dirty = true` в `schema_migrations`: сначала
  разбирают фактическую схему, потом снимают флаг.
- **Первый администратор** после чистой установки назначается SQL-ом
  (`UPDATE users SET role = 'admin' WHERE email = '…';`): регистрация всегда даёт
  `user`, а API назначения ролей требует `user:update` (есть только у `admin`).
  Перелогин не нужен — роль подхватывается из БД на следующем запросе.
- **Проверка** — только на пустой БД: `go build`/`go vet` SQL не валидируют.
  Быстрый цикл: `DROP SCHEMA public CASCADE; CREATE SCHEMA public;` → `go run .` →
  `SELECT version, dirty FROM schema_migrations;` (ожидается `9`, `false`).
- Каталог `migrations/` остаётся на диске, поэтому для ручных операций применим и
  CLI: `migrate -path migrations -database "postgres://…" up|down`.

### Swagger / OpenAPI

Спецификация генерируется `swag` из godoc-аннотаций на хендлерах и отдаётся
`gin-swagger`:

- UI — `http://<BIND_IP>:<PORT>/swagger/index.html`;
- спецификация (Swagger 2.0) — `/swagger/doc.json`;
- сгенерированные файлы лежат в `docs/` (`docs.go`, `swagger.json`,
  `swagger.yaml`) и коммитятся в репозиторий; каталог рядом с бинарником не нужен,
  `docs.go` импортируется в `main.go` и регистрирует спеку при старте.

Генерация после изменения аннотаций:

```
swag init
go build ./...
```

Зависимости (`swaggo/swag`, `swaggo/gin-swagger`, `swaggo/files`) уже в `go.mod`;
CLI ставится отдельно: `go install github.com/swaggo/swag/cmd/swag@latest`.
Версию CLI стоит держать близкой к модулю `swaggo/swag`, иначе шаблоны
сгенерированных файлов могут отличаться.

Аннотации: `@BasePath /api/v1` и `@securityDefinitions.apikey BearerAuth` — в
`main.go`, операции описываются над хендлерами (`@Summary`, `@Param`, `@Success`,
`@Failure`, `@Router`). Для защищённых операций — `@Security BearerAuth`; в UI
токен вводится через **Authorize** в виде `Bearer <access_token>`.

Маршруты `/swagger/*` публичны: они не проходят JWT-middleware и rate limiter.
Для продакшена их нужно закрыть или отключить.

---

## Текущее состояние

### Реализовано

1. **Аутентификация и роли** — `POST /api/v1/auth/register`,
   `POST /api/v1/auth/login`, `POST /api/v1/auth/refresh`,
   `POST /api/v1/auth/logout`; JWT (HS256), хеширование паролей (bcrypt),
   middleware проверки access-токена. Access и refresh различаются claim `typ`;
   `ver` сверяется с `users.token_version`, поэтому logout мгновенно отзывает все
   ранее выданные токены пользователя. Роли `admin`/`user`/`viewer`; роль читается
   из БД на каждом запросе (в токене её нет); проверка разрешений (`task:*`,
   `comment:*`, `user:*`, `profile:*`) на маршрутах через `RequirePermission`.
   Регистрация принимает только `email` и `password` и создаёт аккаунт и пустой
   профиль одной транзакцией; логин возвращает `id` и пару токенов плоским
   объектом.
2. **CRUD задач** — создание, список, получение по id, частичное обновление,
   удаление; фильтрация по `status`/`priority`/`assignee_id` и поиск `search`
   по `name`/`description` (`ILIKE`, до 100 символов). Слои
   `handler → use-case → repository → domain`. Есть создатель задачи (`author_id`).
3. **Комментарии** — `POST/GET /api/v1/tasks/:id/comments`. Слои
   `handler → use-case → repository → domain`. Редактирование, удаление и
   аудит изменений вне объёма MVP: комментарий неизменяем.
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
   на остальных маршрутах — 10 rps (burst 20). На каждую группу маршрутов
   создаётся один инстанс лимитера, поэтому бюджет на IP общий для всех
   маршрутов группы; middleware подключён на уровне группы до JWT-проверки.
   Ответ `429` содержит тело `{"error": "rate limit exceeded"}` и заголовок
   `Retry-After` (целые секунды до следующего токена); неактивные записи
   чистятся по TTL.
10. **Trusted proxies** — `SetTrustedProxies(nil)`: `ClientIP()` берётся из
    соединения, а не из `X-Forwarded-For`, поэтому IP в логах и лимитере не
    подделать. При появлении reverse-proxy понадобится список CIDR.
11. **Роли** — каталог: `GET /api/v1/users/roles` (право `user:read`, только
    `admin`) читает таблицу `roles` и отдаёт `{"roles": [...]}` с сортировкой по
    имени; назначение: `PUT /api/v1/users/:id/roles` (право `user:update`, только
    `admin`); слои `handler → use-case → repository`. Запрещено менять
    собственную роль (`403`, `domain.ErrSelfRoleChange`) и понижать последнего
    администратора (`409`, `domain.ErrLastAdmin`); проверка и обновление идут в
    одной транзакции с блокировкой набора администраторов. Путь `/roles` —
    множественное число текущего контракта. Смена роли применяется немедленно и
    без перелогина: роль читается из БД на каждом защищённом запросе, в токене
    её нет. Создание/удаление ролей и управление `role_permissions` через API в
    MVP не входят: роли и матрица прав задаются сид-миграциями.
12. **Миграции БД** — `golang-migrate`; SQL в `migrations/` встроен в бинарник
    (`//go:embed`, `migrations/embed.go`) и применяется при старте
    (`internal/repository/postgres/migrate.go`); история — в `schema_migrations`;
    DDL всех таблиц и сид ролей/прав — `000001…000009` (см. «Миграции БД»).
13. **Swagger/OpenAPI** — `swag` + `gin-swagger`: аннотации на хендлерах,
    сгенерированные `docs/` (`docs.go`, `swagger.json`, `swagger.yaml`)
    вкоммичены; UI — `/swagger/index.html`, спека — `/swagger/doc.json` (вне
    `/api/v1`, без авторизации); описаны все текущие операции, для защищённых —
    `BearerAuth` (см. «Swagger / OpenAPI»).
14. **Профили** — таблица `profiles` (1:1 с `users`, FK с `ON DELETE CASCADE`);
    `GET/PUT /api/v1/profiles/:id`; права `profile:read`/`profile:update` есть у
    всех ролей, обновление разрешено только для собственного id. При регистрации
    профиль создаётся пустым в одной транзакции с пользователем. `username`,
    `first_name`, `last_name` nullable; комментарии используют
    `COALESCE(profiles.username, '')` для `author_name`.

### Осталось (по MVP)

1. **Инфраструктура** — секреты и параметры БД (включая пул соединений)
   вынесены в конфиг (см. раздел «Конфигурация приложения»), graceful shutdown и
   fail-fast на старте реализованы (`internal/app/server.go`); осталось:
   Docker/compose и тесты.

---

## Вне текущего объёма MVP

- **Динамические роли.** Роли `admin`, `user`, `viewer` фиксированы: каталог
  `roles` доступен только на чтение (`GET /api/v1/users/roles`), назначение
  валидируется `domain.Role.IsValid()`. Создание, удаление и переименование
  ролей, а также управление правами ролей через API не реализованы и в текущий
  объём не входят.
- **Редактирование, удаление и аудит комментариев.** Комментарий неизменяем:
  есть только `POST` и `GET`, в `task_comments` нет `updated_at`/`deleted_at`,
  истории версий и журнала действий нет. Осознанное сужение MVP.

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
- Миграции применяются внутри `postgres.ConnectDB`
  (`internal/repository/postgres/connect.go`), который вызывается из `NewApp`, на
  каждом старте: параллельные реплики сериализуются на advisory-lock
  `golang-migrate`, а падение миграции роняет всё приложение. Для прода миграции
  лучше вынести в отдельный шаг деплоя (init-контейнер или CLI); для локальной
  разработки текущая схема удобна.
- `golang-migrate` не хранит контрольных сумм применённых файлов: правка старой
  миграции не обнаруживается, а расхождение схем проявится позже и в другом
  месте. Защита — только дисциплина «применённую миграцию не редактируем».
- Bootstrap первого администратора — только SQL: регистрация даёт `user`, а
  назначение ролей требует `user:update` (есть у `admin`). Нужен сид или
  CLI-команда для создания администратора без ручного `UPDATE`.
- Автоматического adoption миграций для БД, созданных до их появления, нет:
  нужно пересоздание схемы либо baseline-миграция + `migrate force`.
- Права описаны дважды: константы `domain.Permission` и строки сид-миграций
  `000004`/`000008`. Новое право в коде без строки в сиде даёт `403` на маршруте,
  и наоборот. Матрица прав меняется новой миграцией, а не ручным `INSERT`.
- `context.Context` прокинут сквозь все слои для auth, задач, комментариев, ролей
  и профилей (`handler → use-case → repository`; `QueryContext`/`ExecContext`/
  `BeginTx`). Контекст приходит из `c.Request.Context()`; use-case и репозитории
  не создают `context.Background()` в обычном потоке.
- Роль и `token_version` читаются из БД на каждом защищённом запросе одним
  `SELECT` по первичному ключу (`GetAuthState`): смена роли, включая снижение
  прав, действует сразу и без перелогина, logout отзывает токены мгновенно, а
  токен удалённого пользователя перестаёт работать (`401`). При недоступной БД
  защищённые маршруты отвечают `500` (fail-closed), а не работают по устаревшей
  роли.
- Logout глобальный: `users.token_version` увеличивается для пользователя целиком,
  поэтому выход на одном устройстве завершает сессии на всех. Per-device сессий
  нет — для них нужна таблица сессий.
- Refresh-токены stateless: сервер их не хранит и не ротирует, поэтому старый
  refresh остаётся валиден до `exp` (`JWT_REFRESH_TTL`, 30 дней) или до logout.
  При краже refresh-токена окно компрометации — до истечения срока или выхода
  пользователя; отозвать отдельный refresh-токен и обнаружить его повторное
  использование нельзя (см. «Осталось»).
- Токены, выданные до появления claims `typ`/`ver` (до миграции `000009`), после
  деплоя отклоняются с `401` — нужен повторный логин.
- `/auth/refresh` и `/auth/logout` входят в группу `/api/v1/auth/*` и используют
  строгий лимитер (0.1 rps, burst 5 на IP): при исчерпанном бюджете logout вернёт
  `429`, и ранее выданные токены останутся валидны до `exp`.
- TTL токенов (`JWT_ACCESS_TTL`/`JWT_REFRESH_TTL`) не валидируются при старте:
  нулевое или отрицательное значение даст сразу истёкшие токены.
- `GET /api/v1/users/roles` возвращает каталог всех ролей, хотя по пути читается
  как «роли пользователя»; `PUT /api/v1/users/:id/roles` меняет роль конкретного
  пользователя. Профили уже вынесены в отдельную таблицу, но путь каталога ролей
  пока сохранён: при добавлении `GET /api/v1/users/:id` его придётся перенести в
  `/api/v1/roles`, иначе gin не сможет зарегистрировать статический сегмент и
  `:id` на одном уровне.
- Собственного профиля по токену (`GET/PUT /api/v1/profile`) нет: свой профиль
  запрашивается по `id`, полученному при логине/регистрации, через
  `/api/v1/profiles/:id`. Переход на `/profile` и `/users/:id` — отдельная задача.
- `profiles.username` nullable: при регистрации профиль создаётся пустым, поэтому
  до заполнения профиля `author_name` в комментариях равен пустой строке.
- Откат `000007` восстанавливает `users.username` (для `NULL` и пустых значений —
  локальную часть email), но теряет `first_name`/`last_name` из `profiles`.
- Роли статические: `admin`/`user`/`viewer` заданы в `domain.Role.IsValid()` и
  сид-миграциях, таблица `roles` используется как каталог для `GET`. Роль,
  добавленная в `roles` вручную, появится в каталоге, но не пройдёт валидацию
  назначения, а пользователь с такой ролью получит `500` в JWT-middleware.
  Убирать хардкод и вводить динамические роли — отдельная задача (см. «Вне
  текущего объёма MVP»).
- Поиск задач — `ILIKE '%...%'` по `name`/`description`: ведущий `%` не использует
  индексы, на больших объёмах это seq scan. Нужен `pg_trgm` + GIN или FTS.
- `SlogMiddleware` (`internal/app/server.go`): ветка `c.Errors` не используется
  (никто не вызывает `c.Error`), не проставляется request id, `query`
  логируется без фильтрации.
- Параметры пула читаются из `DB_MAX_OPEN_CONNS`/`DB_MAX_IDLE_CONNS`;
  `SetConnMaxLifetime` зафиксирован константой `1m`. При ошибке `Ping` или
  миграций `ConnectDB` не закрывает пул.
- Rate limiting хранит карту клиентов в памяти процесса: записи чистятся по TTL
  только при поступлении запросов, а бюджет считается по `ClientIP()` — за общим
  NAT или reverse-proxy без списка доверенных CIDR все клиенты делят один лимит.
- Swagger UI (`/swagger/index.html`) и спецификация (`/swagger/doc.json`)
  публичны: маршрут не проходит JWT-middleware и rate limiter; в проде его нужно
  закрыть или отключить.
- `docs/` генерируются вручную (`swag init`), проверки актуальности в CI нет:
  после правки аннотаций легко забыть перегенерировать файлы. Версия CLI `swag`
  может расходиться с модулем `swaggo/swag` — фиксировать при обновлении.
- Нет тестов; Docker/compose не подключены.