# HTTPtodo

Учебный REST API для задач на Go (net/http, без фреймворков).

## Архитектура
- handlers → service (TaskService) → storage (интерфейс Storage) → FileStorage (JSON-файл)
- Sentinel-ошибки ErrNotFound / ErrEmptyText, различаются через errors.Is
- Хендлеры отвечают через respondJSON / respondError

## Планы
- Разнести по файлам, тесты с MemoryStorage, потом PostgreSQL как вторая реализация Storage

## Как со мной работать
Я учу Go — объясняй концепциями, не пиши готовый код за меня.
Показывай направление и разбирай ошибки, писать буду сам.
