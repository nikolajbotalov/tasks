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

### Структура API (эндпоинты)

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
| `JWT_ACCESS_TTL` | нет | `24h` | время жизни access-токена; на актуальность роли не влияет |
| `DB_HOST` | нет | `localhost` | хост PostgreSQL |
| `DB_PORT` | нет | `5432` | порт PostgreSQL |
| `DB_NAME` | нет | `taskflow` | имя БД |
| `DB_USER` | нет | `postgres` | пользователь БД |
| `DB_PASSWORD` | да | — | пароль БД |
| `DB_SSL_MODE` | нет | `disable` | режим SSL (`disable`/`require`/`verify-full`) |

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

### Миграции БД

Схема БД описана миграциями `golang-migrate`; ручное создание таблиц не
используется. SQL-файлы лежат в `migrations/` и встраиваются в бинарник через
`//go:embed` (`migrations/embed.go`), поэтому каталог с миграциями рядом с
бинарником не нужен. Применение — при старте приложения:
`postgres.RunMigrations` (`internal/repository/postgres/migrate.go`) вызывает
`Up()`, история версий хранится в таблице `schema_migrations`.

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

Правила:

- **Применённую миграцию не редактируют.** Изменения схемы — новой версией
  (`000009_*.up.sql`). `golang-migrate` не хранит контрольных сумм, поэтому
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
  `SELECT version, dirty FROM schema_migrations;` (ожидается `8`, `false`).
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

1. **Аутентификация и роли** — `POST /api/v1/auth/register`, `POST /api/v1/auth/login`,
   JWT (HS256), хеширование паролей (bcrypt), middleware проверки токена.
   Роли `admin`/`user`/`viewer`; роль читается из БД на каждом запросе (в токене
   её нет); проверка разрешений (`task:*`, `comment:*`, `user:*`, `profile:*`)
   на маршрутах через `RequirePermission`. Регистрация принимает только `email`
   и `password` и создаёт аккаунт и пустой профиль одной транзакцией; логин
   возвращает `id` и пару токенов плоским объектом.
2. **CRUD задач** — создание, список, получение по id, частичное обновление,
   удаление; фильтрация по `status`/`priority`/`assignee_id` и поиск `search`
   по `name`/`description` (`ILIKE`, до 100 символов). Слои
   `handler → use-case → repository → domain`. Есть создатель задачи (`author_id`).
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
11. **Роли** — каталог: `GET /api/v1/users/roles` (право `user:read`, только
    `admin`) читает таблицу `roles` и отдаёт `{"roles": [...]}` с сортировкой по
    имени; назначение: `PUT /api/v1/users/:id/roles` (право `user:update`, только
    `admin`); слои `handler → use-case → repository`. Запрещено менять
    собственную роль (`403`, `domain.ErrSelfRoleChange`) и понижать последнего
    администратора (`409`, `domain.ErrLastAdmin`); проверка и обновление идут в
    одной транзакции с блокировкой набора администраторов. Путь `/roles` —
    множественное число оставлено под будущие операции со списком ролей.
    Смена роли применяется немедленно и без перелогина: роль читается из БД на
    каждом защищённом запросе, в токене её нет.
12. **Миграции БД** — `golang-migrate`; SQL в `migrations/` встроен в бинарник
    (`//go:embed`, `migrations/embed.go`) и применяется при старте
    (`internal/repository/postgres/migrate.go`); история — в `schema_migrations`;
    DDL всех таблиц и сид ролей/прав — `000001…000008` (см. «Миграции БД»).
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

1. **Управление ролями** — базовое назначение роли и защита инвариантов
   реализованы (`PUT /api/v1/users/:id/roles`, право `user:update`, только
   `admin`; самосмена → `403`, последний администратор → `409`). Осталось:
   создание и удаление ролей через API не реализованы: каталог `roles` и права
   создаются сид-миграциями `000004`/`000008`, а валидация назначения опирается
   на хардкод `domain.Role.IsValid()` — расхождение описано в «Известных слабых
   местах».
2. **Комментарии** — реализованы; осталось логирование изменений (audit).
3. **Инфраструктура** — секреты и параметры БД вынесены в конфиг (см. раздел
   «Конфигурация приложения»), graceful shutdown и fail-fast на старте
   реализованы (`internal/app/server.go`); осталось: Docker/compose,
   тесты, вынос настройки пула соединений в
   `internal/repository/postgres/connect.go`.
4. **Rate limiting / `429 Too Many Requests`** — базовое ограничение частоты
   реализовано (`internal/handlers/middleware/rate_limiter.go`); осталось:
   заголовок `Retry-After` и общий лимит поверх маршрутов (сейчас лимитер
   навешен на каждый маршрут отдельно).
5. **Сессии и отзыв токенов** — `refresh_token` генерируется при логине, но нигде
   не хранится и не проверяется, эндпоинта обновления нет. Отозвать уже выданный
   access-токен нельзя: он действует до `exp` (`JWT_ACCESS_TTL`, по умолчанию
   24h). Нужны таблица сессий, ротация refresh-токена, `logout` и `token_version`
   для мгновенного отзыва.

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
- Миграции применяются внутри `NewApp` на каждом старте: параллельные реплики
  сериализуются на advisory-lock `golang-migrate`, а падение миграции роняет всё
  приложение. Для прода миграции лучше вынести в отдельный шаг деплоя
  (init-контейнер или CLI); для локальной разработки текущая схема удобна.
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
- `context.Context` прокинут в auth (handler → use-case → repository); в
  use-case/repository комментариев и задач контекст пока не внедрён.
- Роль читается из БД на каждом защищённом запросе: смена роли, включая снижение
  прав, действует сразу и без перелогина, а токен удалённого пользователя
  перестаёт работать (`401`). Цена — дополнительный `SELECT` по первичному ключу
  на каждый запрос; при недоступной БД защищённые маршруты отвечают `500`
  (fail-closed), а не работают по устаревшей роли.
- Отзыва токенов нет: смена пароля, «выйти на всех устройствах» и кража
  access-токена не инвалидируют уже выданный токен — он действует до `exp`
  (`JWT_ACCESS_TTL`, по умолчанию 24h). Для этого нужны сессии и `token_version`
  (см. «Осталось»).
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
- Роли описаны в трёх местах: таблица `roles` (каталог для `GET`),
  `domain.Role.IsValid()` (хардкод `viewer/user/admin` для валидации назначения) и
  `role_permissions` (права). Создание/удаление ролей через API не реализованы,
  поэтому расхождение возможно: роль, добавленная в `roles`, попадёт в `GET`, но
  `PUT` вернёт `422`, пока её не добавят в `IsValid()`. Устранить при переходе на
  динамические роли с единым источником в `roles`.
- Поиск задач — `ILIKE '%...%'` по `name`/`description`: ведущий `%` не использует
  индексы, на больших объёмах это seq scan. Нужен `pg_trgm` + GIN или FTS.
- `SlogMiddleware` (`internal/app/server.go`): ветка `c.Errors` не используется
  (никто не вызывает `c.Error`), не проставляется request id, `query`
  логируется без фильтрации.
- Настройка пула соединений находится в `internal/app/app.go`, а не в
  `internal/repository/postgres/connect.go` (см. раздел «Подключение к БД»).
- Swagger UI (`/swagger/index.html`) и спецификация (`/swagger/doc.json`)
  публичны: маршрут не проходит JWT-middleware и rate limiter; в проде его нужно
  закрыть или отключить.
- `docs/` генерируются вручную (`swag init`), проверки актуальности в CI нет:
  после правки аннотаций легко забыть перегенерировать файлы. Версия CLI `swag`
  может расходиться с модулем `swaggo/swag` — фиксировать при обновлении.
- Нет тестов; Docker/compose не подключены.