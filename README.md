## Delayed Notifier (Go)

Сервис отложенных уведомлений с HTTP API, RabbitMQ, PostgreSQL, Redis (опционально) и простым UI.

### Переменные окружения

- `APP_HTTP_ADDR` — адрес HTTP-сервера (по умолчанию `:8080`).
- `APP_DB_DSN` — DSN PostgreSQL (по умолчанию `postgres://user:password@localhost:5432/delayed_notifier?sslmode=disable`).
- `APP_RABBIT_URL` — URL RabbitMQ (по умолчанию `amqp://guest:guest@localhost:5672/`).
- `APP_RABBIT_QUEUE` — имя очереди RabbitMQ (по умолчанию `notifications`).
- `APP_REDIS_ADDR` — адрес Redis (по умолчанию `localhost:6379`).
- `APP_REDIS_PASSWORD` — пароль Redis (по умолчанию пусто).
- `APP_REDIS_DB` — номер базы Redis (по умолчанию `0`).
- `APP_BASE_RETRY_DELAY` — базовая задержка ретраев (например, `1m`).
- `APP_MAX_RETRY` — максимальное количество ретраев (по умолчанию `5`).
- **Email (SMTP):** `APP_SMTP_HOST`, `APP_SMTP_PORT` (по умолчанию 587), `APP_SMTP_USER`, `APP_SMTP_PASSWORD`, `APP_EMAIL_DEFAULT_RECIPIENT`. Если `APP_SMTP_HOST` не задан — письма не отправляются (заглушка).
- **Telegram:** `APP_TELEGRAM_TOKEN` (токен бота от [@BotFather](https://t.me/BotFather)), `APP_TELEGRAM_DEFAULT_RECIPIENT` (Telegram ID получателя, можно получить через tg-бота @myidbot). Если токен не задан — сообщения в Telegram не отправляются (заглушка).

### Сборка и запуск

```bash
go mod tidy
go build -o delayed-notifier ./cmd/app
./delayed-notifier
```

### HTTP API

- `POST /api/notify` — создать уведомление.
- `GET /api/notify` — получить список последних уведомлений.
- `GET /api/notify/{id}` — получить статус уведомления по его ID.
- `DELETE /api/notify/{id}` — отменить уведомление.

Пример запроса создания:

```bash
curl -X POST http://localhost:8080/api/notify \
  -H "Content-Type: application/json" \
  -d '{
    "channel": "email",
    "recipient": "user@example.com",
    "message": "Hello!",
    "send_at": "2026-03-12T10:00:00Z"
  }'
```

### UI

Откройте в браузере `http://localhost:8080/` — форма создания уведомления и таблица текущих уведомлений и их статусов.

