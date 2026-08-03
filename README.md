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
httpapi.Handler ──(интерфейс TaskService)──► service.Service ──(интерфейс Storage)──► postgres.Storage ──► PostgreSQL
```

Интерфейсы объявлены на стороне потребителя: `TaskService` — в `httpapi`, `Storage` — в `service`.
Благодаря этому реализации хранилища взаимозаменяемы (PostgreSQL или файл), а верхние слои
о конкретной реализации не знают.

- **`internal/api/httpapi`** — HTTP-слой. Разбор запроса, коды ответов, JSON, роутинг и запуск
  сервера. Ответы через `respondJSON` / `respondError`.
- **`internal/service`** — бизнес-логика (`Service`). Валидация (например, пустой текст)
  и проброс в хранилище. Объявляет интерфейс `Storage`.
- **`internal/repository/postgres`** — реализация `Storage` поверх PostgreSQL (`pgx`).
- **`internal/repository/storage`** — альтернативная реализация `Storage` поверх JSON-файла
  (`FileStorage`).
- **`internal/model`** — доменные модели (`Task`) и sentinel-ошибки
  `ErrNotFound` / `ErrEmptyText`, различаемые через `errors.Is`.
- **`internal/config`** — загрузка конфигурации из окружения / `.env`.
- **`cmd/httptodo`** — точка входа: сборка зависимостей и запуск сервера.

### Структура проекта

```
cmd/httptodo/main.go            точка входа (сборка зависимостей)
internal/
  api/httpapi/                  HTTP-слой: хендлеры, роутинг, сервер
    handler.go                  тип Handler + интерфейс TaskService
    tasks.go                    хендлеры ручек /tasks
    respond.go                  хелперы ответа (JSON / ошибки)
    routes.go                   регистрация маршрутов (NewRouter)
    httpserver.go               запуск HTTP-сервера (Server)
  service/service.go            бизнес-логика + интерфейс Storage
  repository/postgres/          реализация Storage на PostgreSQL
  repository/storage/           файловая (JSON) реализация Storage
  model/                        модели и ошибки
  config/                       конфигурация
migrations/                     SQL-миграции схемы БД (goose)
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

```bash
docker compose up -d
```

### 2. Настроить окружение

```bash
cp .env.example .env
# затем открыть .env и подставить реальные DATABASE_URL и PORT
```

### 3. Применить миграции

Схема БД версионируется через [goose](https://github.com/pressly/goose).
Установите утилиту и накатите миграции (goose читает `DATABASE_URL` из окружения,
`.env` он сам не подхватывает — экспортируйте переменную заранее):

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest

export $(grep -v '^#' .env | xargs)
goose -dir migrations postgres "$DATABASE_URL" up
goose -dir migrations postgres "$DATABASE_URL" status
```

### 4. Запустить приложение

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
