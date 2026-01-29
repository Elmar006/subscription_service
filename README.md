# subscription_service

Сервис управления подписками на различные сервисы.
Реализован на Go, использует PostgreSQL для хранения данных и предоставляет REST API с документацией Swagger.

Сервис позволяет:

1) Создавать, обновлять, получать и удалять подписки.
2) Получать список подписок пользователя.
3) Считать общую сумму расходов пользователя за определённый период.

Технологии:

1) GO v1.25
2) PostgreSQL
3) Docker / Docker Compose
4) Chi Router в качестве маршрутизатора
5) Swagger для документации
6) И миграции



Установка и запуск

1. Клонируем репозиторий:

```bash
git clone https://github.com/Elmar006/subscription_service.git
cd subscription_service
```

2. Устанавливаем зависимости:

```bash
go mod tidy
```

3. Создаем `.env` файл с конфигурацией базы данных если нужно:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=subscriptions
APP_ENV=dev
```

 Для продакшена рекомендуется использовать `APP_ENV=prod` для удобного логирования в формате JSON.

4. Запускаем миграции:

```bash
docker compose up -d postgres
```

5. Запускаем сервис:

```bash
go run cmd/main.go
```

Сервис будет доступен на: `http://localhost:8080`



Docker Compose
Используем Docker Compose для локальной разработки:

Запуск:

```bash
docker compose up -d
```



API Эндпоинты

| Метод  | Путь                                                                                                 | Описание                     |
| ------ | ---------------------------------------------------------------------------------------------------- | ---------------------------- |
| POST   | /subscriptions                                                                                       | Создать подписку             |
| GET    | /subscriptions/{id}                                                                                  | Получить подписку по ID      |
| GET    | /subscriptions?user_id={user_id}                                                                     | Список подписок пользователя |
| PUT    | /subscriptions/{id}                                                                                  | Обновить подписку            |
| DELETE | /subscriptions/{id}                                                                                  | Удалить подписку             |
| GET    | /subscriptions/total?user_id={user_id}&service_name={service_name}&from={yyyy-mm-dd}&to={yyyy-mm-dd} | Общая сумма по фильтрам      |




Swagger:

Swagger-документация доступна по адресу: `http://localhost:8080/swagger/index.html`

Пример тестового запроса через Swagger:

POST /subscriptions

```json
{
  "service_name": "Netflix",
  "price": 550,
  "user_id": "e4f1c2a7-9b3d-4f5e-a2d1-8c7f6b9d2e3a",
  "start_date": "2026-01-01",
  "end_date": "2026-12-31"
}
```

PUT /subscriptions/{id}

```json
{
  "service_name": "Netflix Premium",
  "price": 600,
  "start_date": "2026-01-01",
  "end_date": "2026-12-31"
}
```

GET /subscriptions?user_id={user_id}

```
user_id=e4f1c2a7-9b3d-4f5e-a2d1-8c7f6b9d2e3a
```

GET /subscriptions/{id}

```
id=d34728e9-6c01-4e74-ab3a-1bb2bef2a342
```

GET /subscriptions/total

```
user_id=e4f1c2a7-9b3d-4f5e-a2d1-8c7f6b9d2e3a&service_name=Netflix&from=2026-01-01&to=2026-12-31
```

DELETE /subscriptions/{id}

```
id=d34728e9-6c01-4e74-ab3a-1bb2bef2a342
```



Тестирование:

Тесты находятся в `internal/repository/subscriptions_test.go` и `internal/handler/handler_test.go`
Запуск:

```bash
go test ./internal/repository -v
```

```bash
go test ./internal/handler -v
```

Тесты проверяют:

Создание подписки
Получение по ID
Обновление подписки
Удаление подписки
Список подписок пользователя
Общую сумму подписок по фильтрам


CI/CD
Каждый push запускает unit-тесты 
Пуш с тегом запускаем процесс деплоя

Логирование
Для логирования в проекте используется библиотека Logrus
В режиме разработки .env/ APP_ENV=dev логи выводятся подробно в консоль, в удобном формате для чтения и отладки.
Для продакшена рекомендуется выставить `APP_ENV=prod` вместо `APP_ENV=dev` `.env` для удобного логирования в JSON формате


Контакты:

Автор: Эльмар
Электронная почта: birembekove@gmail.com
Telegram: @viy1ix
GitHub: https://github.com/Elmar006


