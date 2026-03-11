# transya

Сервис перевода текстов на базе Yandex Cloud Translate API.

Получает задачи на перевод из очереди NATS, выполняет перевод через Yandex Cloud Translate и публикует результаты в Kafka.

## Архитектура

```
NATS (subject: notifications) --> transya --> Kafka (topic: myTopic)
```

1. Сервис подписывается на NATS-subject `notifications` в queue group `workers`.
2. Из входящего сообщения извлекает текст и целевой язык.
3. Отправляет запрос в Yandex Cloud Translate API.
4. Публикует обогащённое сообщение (с полем `translatedText`) в Kafka.

## Конфигурация

Настройки задаются через переменные окружения или файл `.env` в корне проекта.

| Переменная    | Обязательная |  По умолчанию    | Описание                                              |
|---------------|:------------:|------------------|-------------------------------------------------------|
| `FOLDER_ID`   |      да      | —                | ID каталога в Yandex Cloud                            |
| `API_KEY`     |      да      | —                | API-ключ для аутентификации в Yandex Cloud Translate  |
| `NATS_URL`    |      да      | `localhost:4222` | Адрес NATS-сервера                                    |
| `KAFKA_ADDR`  |      да      | `localhost:9093` | Адрес Kafka-брокера                                   |
| `KAFKA_TOPIC` |      да      | `myTopic`        | Kafka-топик для публикации результатов                |

Пример `.env` файла (`.env.example`):

```env
FOLDER_ID=""
API_KEY=""
NATS_URL="localhost:4222"
KAFKA_ADDR="localhost:9093"
KAFKA_TOPIC="myTopic"
```

## Формат сообщений

### Входящее сообщение (NATS)

```json
{
  "requestHash": "abc123",
  "language": "en",
  "textHash": "def456",
  "text": "Привет, мир!",
  "translatedText": "",
  "statusCode": false,
  "errorText": ""
}
```

| Поле            | Тип    | Описание                                              |
|-----------------|--------|-------------------------------------------------------|
| `requestHash`   | string | Идентификатор запроса                                 |
| `language`      | string | Целевой язык перевода (код ISO 639-1, например `en`, `ru`) |
| `textHash`      | string | Идентификатор текста                                  |
| `text`          | string | Текст для перевода                                    |

### Исходящее сообщение (Kafka)

Те же поля, что и во входящем, плюс заполненное поле `translatedText`.

Kafka-заголовки:
- `reqId` — значение `requestHash`
- `textId` — значение `textHash`

## Запуск

```bash
cp .env.example .env
# заполните .env своими значениями

go run ./cmd/main.go
```

## Зависимости

- [nats.go](https://github.com/nats-io/nats.go) — клиент NATS
- [confluent-kafka-go](https://github.com/confluentinc/confluent-kafka-go) — клиент Kafka
- [zerolog](https://github.com/rs/zerolog) — структурированное логирование
- [godotenv](https://github.com/joho/godotenv) — загрузка `.env` файла
- [env](https://github.com/caarlos0/env) — парсинг переменных окружения
