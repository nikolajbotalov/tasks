## Проект: **"TaskFlow" — Таск-трекер для команд**

### Описание
Сервис управления задачами с ролевой моделью. Команда может создавать задачи, назначать их на исполнителей, комментировать и отслеживать прогресс.

---

### Основные фичи (MVP)

1. **Аутентификация + роли**
    - JWT или session-based auth
    - Роли: `admin`, `manager`, `member`, `viewer`
    - Разные права доступа к CRUD задачам

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

## Что предстоит реализовать (по MVP)

1. **Аутентификация и роли** — не начато: нет `POST /api/auth/register`,
   `POST /api/auth/login`, нет JWT, нет хеширования паролей, нет middleware
   проверки роли (`admin`, `manager`, `member`, `viewer`). Структура `User` —
   заглушка и нигде не используется.
2. **Расширение модели задачи** — нет полей `priority`, `status`, `assignee_id`,
   `creator_id`, тегов. Сейчас только `name` + `description`.
3. **Миграции БД** — файлов миграций нет вообще (`golang-migrate` не подключён),
   SQL-схема проекта отсутствует, таблица `tasks` создана вручную в локальной БД.
4. **Комментарии** — не начато: ни домена, ни таблицы, ни эндпоинтов
   `POST/GET /api/tasks/:id/comments`, ни логирования изменений.
5. **Фильтрация, поиск, пагинация** — параметры `status`, `priority`,
   `assignee_id`, `search` игнорируются, пагинация работает некорректно.
6. **Swagger / OpenAPI** — отсутствует, как и валидация через
   `go-playground/validator`.
7. **Инфраструктура** — нет слоёв repository/usecase, нет конфига (секреты в
   коде), нет graceful shutdown, CORS, Docker/compose, `.gitignore` и
   git-репозитория, нет тестов.

---

## Явные баги и слабые места в текущем коде

- `internal/handlers/tasks/get_all.go:17` — `LIMIT $1 OFFSET $2` получает
  **строки**, а не числа; `OFFSET` равен `page` вместо `(page-1)*limit`;
  `Total` — длина текущей страницы, а не общее число записей. Ошибка БД только
  логируется, наружу уходит `200` с пустым списком.
- `internal/handlers/tasks/get_by_id.go:19-23` — при ошибке `Scan`
  (несуществующий ID) всё равно возвращается `200` и пустой объект вместо `404`.
- `internal/handlers/tasks/update.go` — ошибка БД возвращается как `500` с
  текстом ошибки БД (утечка внутренних деталей).
- `internal/handlers/tasks/delete.go` — паттерн «SELECT для проверки
  существования, затем DELETE» создаёт лишний запрос и гонку.
- `internal/handlers/tasks/create.go` — проверка `to_regclass` на существование
  таблицы выполняется в каждом запросе (лишняя логика в хендлере).
- `main.go` — при ошибке `sql.Open`/`Ping` программа продолжает работу;
  `r.Run()` без обработки ошибки.
- Во всех хендлерах SQL размазан по слою доставки — перед добавлением auth и
  комментариев потребуется рефакторинг в repository/usecase.