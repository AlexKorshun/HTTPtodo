# HTTPtodo

Учебный REST API для управления задачами (todo) на Go.
Написан на стандартном `net/http` (роутинг Go 1.22+, без сторонних фреймворков),
хранение — PostgreSQL через `pgx`.

## Возможности

- Просмотр списка задач
- Создание задачи
- Переключение статуса «выполнено / не выполнено»
- Удаление задачи

## Архитектура

Слоёная архитектура, зависимости направлены внутрь через интерфейсы:

```
HTTP-запрос
   │
   ▼
api.Handler ──(интерфейс TaskService)──► service.TaskService ──(интерфейс Storage)──► postgres.PostgresStorage ──► PostgreSQL
```

- **`internal/api`** — HTTP-хендлеры. Разбор запроса, коды ответов, JSON.
  Ответы через `respondJSON` / `respondError`.
- **`internal/service`** — бизнес-логика (`TaskService`). Валидация (например, пустой текст)
  и проброс в хранилище. Объявляет интерфейс `Storage`.
- **`internal/repository/postgres`** — реализация `Storage` поверх PostgreSQL (`pgx`).
- **`internal/repository/storage`** — вспомогательные типы/маппинг хранилища.
- **`internal/model`** — доменные модели (`Task`) и sentinel-ошибки
  `ErrNotFound` / `ErrEmptyText`, различаемые через `errors.Is`.
- **`internal/config`** — загрузка конфигурации из окружения / `.env`.
- **`cmd/httptodo`** — точка входа: сборка зависимостей и запуск сервера.

### Структура проекта

```
cmd/httptodo/main.go            точка входа
internal/
  api/handlers.go               HTTP-хендлеры
  service/service.go            бизнес-логика + интерфейс Storage
  repository/postgres/          реализация Storage на PostgreSQL
  repository/storage/           маппинг/вспомогательное хранилище
  model/                        модели и ошибки
  config/                       конфигурация
schema.sql                      схема БД
docker-compose.yml              PostgreSQL для локального запуска
```

## Конфигурация

Приложение читает переменные окружения (и `.env`, если файл есть):

| Переменная     | Обязательна | По умолчанию | Описание                       |
|----------------|-------------|--------------|--------------------------------|
| `DATABASE_URL` | да          | —            | строка подключения к PostgreSQL |
| `PORT`         | нет         | `8080`       | порт HTTP-сервера              |

Заготовку возьмите из [`.env.example`](.env.example) и подставьте свои значения:

```env
DATABASE_URL=postgres://<user>:<password>@localhost:5432/<db>
PORT=<port>
```


## Запуск

### 1. Поднять PostgreSQL

`docker-compose` поднимает Postgres и автоматически применяет `schema.sql`:

```bash
docker compose up -d
```

### 2. Настроить окружение

```bash
cp .env.example .env
# при необходимости отредактировать DATABASE_URL / PORT
```

### 3. Запустить приложение

```bash
go run ./cmd/httptodo
```

Сервер поднимется на `http://localhost:8080` (или на указанном `PORT`).

### Сборка бинарника

```bash
go build -o httptodo ./cmd/httptodo
./httptodo
```

## API

Базовый URL: `http://localhost:8080`

| Метод    | Путь          | Описание                          | Тело запроса        | Успешный ответ         |
|----------|---------------|-----------------------------------|---------------------|------------------------|
| `GET`    | `/tasks`      | Список всех задач                 | —                   | `200 OK` — массив задач |
| `POST`   | `/tasks`      | Создать задачу                    | `{"text": "..."}`   | `201 Created` — задача  |
| `PATCH`  | `/tasks/{id}` | Переключить статус `done`         | —                   | `200 OK` — задача       |
| `DELETE` | `/tasks/{id}` | Удалить задачу                    | —                   | `204 No Content`        |

### Модель задачи

```json
{
  "id": 1,
  "text": "купить хлеб",
  "done": false
}
```

### Ошибки

Ошибки возвращаются в виде JSON: `{"error": "описание"}`.

| Код   | Когда                                                        |
|-------|-------------------------------------------------------------|
| `400` | некорректный JSON, пустой `text`, некорректный `id`         |
| `404` | задача не найдена                                            |
| `500` | внутренняя ошибка сервера                                    |

### Примеры (curl)

```bash
# создать
curl -X POST localhost:8080/tasks -d '{"text":"купить хлеб"}'

# список
curl localhost:8080/tasks

# переключить done у задачи с id=1
curl -X PATCH localhost:8080/tasks/1

# удалить
curl -X DELETE localhost:8080/tasks/1
```

## Планы

- Тесты с in-memory реализацией `Storage`
- Расширение бизнес-логики (редактирование текста, фильтры)
