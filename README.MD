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

## Текущее состояние

### Реализовано

1. **Аутентификация** — `POST /api/auth/register`, `POST /api/auth/login`,
   JWT (HS256), хеширование паролей (bcrypt), middleware проверки токена.
   Ролей пока нет — только аутентификация.
2. **CRUD задач** — создание, список, получение по id, частичное обновление,
   удаление. Слои `handler → use-case → repository → domain`. Есть создатель
   задачи (`author_id`).
3. **Инфраструктура** — слои repository/usecase, пакет `config`, `.gitignore`,
   git-репозиторий.

### Осталось (по MVP)

1. **Роли** (`admin`, `manager`, `member`, `viewer`) и middleware проверки роли.
2. **Расширение модели задачи** — `priority`, `status`, `assignee_id`, теги.
   Создатель уже реализован (`author_id`).
3. **Миграции БД** — `golang-migrate` не подключён, схема создана вручную.
4. **Комментарии** — ни домена, ни таблицы, ни эндпоинтов
   `POST/GET /api/tasks/:id/comments`, ни логирования изменений.
5. **Фильтрация и поиск** — параметры `status`, `priority`, `assignee_id`,
   `search` пока не реализованы; пагинация уже исправлена.
6. **Swagger / OpenAPI** — отсутствует (валидация через
   `go-playground/validator` уже используется).
7. **Инфраструктура** — секреты из окружения (сейчас DSN в коде), graceful
   shutdown, CORS, Docker/compose, тесты.

---

## Известные слабые места

- `main.go` — при ошибке `sql.Open`/`Ping` программа продолжает работу;
  `r.Run()` без обработки ошибки.
- `context.Context` не прокинут в use-case/repository.
- Логирование через `fmt.Printf`; переход на `log/slog` — отдельная задача.