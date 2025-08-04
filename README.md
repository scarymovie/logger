# `scarymovie/logger`

**Go-пакет для логирования на базе `log/slog` с поддержкой контекста, групп, глобальных атрибутов и форматирования.**

## 🚀 Установка

```
go get github.com/scarymovie/logger
```

## 🧠 Возможности

*   Инициализация логгера с глобальными и сгруппированными атрибутами
*   Middleware-обработчик для вставки значений из `context.Context`
*   Поддержка `JSON` и `Text` форматов
*   Функции `WithRequestID` и `WithLogMessage` для прокидывания данных
*   Функция `WrapError` для обёртывания ошибок с контекстом

## 🛠 Пример использования

### Инициализация

```
logger.NewLogger(logger.Config{
    Level: slog.LevelInfo,
    JSONFormat: true,
    DefaultAttrs: []slog.Attr{
        slog.String("version", "1.0.2"),
    },
    GroupedAttrs: []logger.GroupedAttrs{
        {
            Group: "service",
            Attrs: []slog.Attr{
                slog.String("name", "payment-api"),
                slog.String("instance_id", "srv-01"),
            },
        },
        {
            Group: "env",
            Attrs: []slog.Attr{
                slog.String("env", "prod"),
                slog.String("region", "eu-central"),
            },
        },
    },
})
```

### Логирование с контекстом

```
ctx := logger.WithRequestID(context.Background(), "abc-123")
ctx = logger.WithLogMessage(ctx, "initializing user session")

logger.GetLogger().InfoContext(ctx, "session started")
```

### Работа с ошибками

```
err := someOperation()
if err != nil {
    err = logger.WrapError(ctx, err)
    logger.GetLogger().ErrorContext(ctx, "operation failed", slog.Any("error", err))
}
```

## 📦 API

### `type Config`

Конфигурация логгера:

```
type Config struct {
    Level        slog.Level
    JSONFormat   bool
    DefaultAttrs []slog.Attr
    GroupedAttrs []GroupedAttrs
}
```

### `func NewLogger(cfg Config)`

Инициализирует глобальный логгер. Без повторной инициализации.

### `func GetLogger() *slog.Logger`

Возвращает глобальный `slog.Logger`.

### `func WithRequestID(ctx context.Context, id string) context.Context`

Добавляет request ID в контекст.

### `func WithLogMessage(ctx context.Context, msg string) context.Context`

Добавляет сообщение в контекст.

### `func WrapError(ctx context.Context, err error) error`

Оборачивает ошибку с лог-контекстом.

## 📋 Лицензия

MIT

- - -

Автор: **scarymovie** 

Email: **vino.zeka@gmail.com**

Репозиторий: [github.com/scarymovie/logger](https://github.com/scarymovie/logger)