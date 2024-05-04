# go-musthave-metrics-tpl

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.

### Запуск сервера

```shell
go run cmd/server/main.go
```

###### Конфигурация

Флаги

- -h - показать все команды
- -a - адрес и порт сервера
- -i - интевал сохранения метрик в файл
- -f - файл для созранения метрик. По деволту: /tmp/metrics-db.json
- -r - загружать ли при запуске метрики из файла
- -d - строка подключения к базе данных
- -r - ключ sha256
- -c - путь к файлу конфигов

Env

- ADDRESS - адрес и порт сервера
- STORE_INTERVAL - интевал сохранения метрик в файл
- FILE_STORAGE_PATH - файл для созранения метрик. По деволту: /tmp/metrics-db.json
- RESTORE - загружать ли при запуске метрики из файла
- DATABASE_DSN - строка подключения к базе данных
- KEY - ключ sha256
- CONFIG - путь к файлу конфигов

```go
postgres://postgres:password@localhost:5433/postgres?sslmode=disable
```

### Запуск агента

```shell
go run cmd/agent/main.go
```

###### Конфигурация

Флаги

- -a - адрес и порт сервера
- -r - интервал отправки метрик на сервер
- -p - интервал сбора метрик
- -k - ключ sha256
- -l - лимит запросов

Env

- ADDRESS - адрес и порт сервера
- REPORT_INTERVAL - интервал отправки метрик на сервер
- POLL_INTERVAL - интервал сбора метрик
- KEY - ключ sha256
- RATE_LIMIT - лимит запросов

### Запуск статического анализатора

```shell
go run cmd/staticlint/main.go ./...
```

Анализатор включает в себя:

- стандартные анализаторы из пакета golang.org/x/tools/go/analysis/passes;
- все анализаторы класса SA пакета staticcheck.io;
- анализатор ST1000 Incorrect or missing package comment
- анализатор github.com/fatih/errwrap/errwrap Проверка формата вывода ошибки
- анализатор github.com/timakin/bodyclose/passes/bodyclose Проверка закрытия тела запроса
- анализатор github.com/kisielk/errcheck/errcheck Проверки непроверенных ошибок в коде
- анализатор для проверки os.Exit в пакете main в функции main
- анализатор для проверки fmt.Print в коде
- анализатор для проверки закомментированного кода

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-metrics-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск тестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).

Далее представлены команды для запуска поочередного тестирования

```shell
./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -test.run=^TestIteration1$
```

```shell
./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -test.run=^TestIteration2A$
```

```shell
./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -test.run=^TestIteration2B$
```

```shell
./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -test.run=^TestIteration3A$
```

```shell
./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -test.run=^TestIteration3B$
```

```shell
./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -server-port=8080 -test.run=^TestIteration4$
```

```shell
./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -server-port=8080 -test.run=^TestIteration5$
```

```shell
./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -server-port=8080 -test.run=^TestIteration6$
```

```shell
./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -server-port=8080 -test.run=^TestIteration7$
```

```shell
./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -server-port=8080 -test.run=^TestIteration8$
```

```shell
./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -server-port=8080 -file-storage-path=/tmp/metrics-db.json -test.run=^TestIteration9$
```

```shell
./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -server-port=8080 -file-storage-path=/tmp/metrics-db.json -database-dsn='postgres://postgres:password@localhost:5433/postgres?sslmode=disable' -test.run=^TestIteration10A$
```

```shell
./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -server-port=8080 -file-storage-path=/tmp/metrics-db.json -database-dsn='postgres://postgres:password@localhost:5433/postgres?sslmode=disable' -test.run=^TestIteration10B$
```

```shell
./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -server-port=8080 -file-storage-path=/tmp/metrics-db.json -database-dsn='postgres://postgres:password@localhost:5433/postgres?sslmode=disable' -test.run=^TestIteration11$
```

```shell
./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -server-port=8080 -file-storage-path=/tmp/metrics-db.json -database-dsn='postgres://postgres:password@localhost:5433/postgres?sslmode=disable' -test.run=^TestIteration12$
```

```shell
./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -server-port=8080 -file-storage-path=/tmp/metrics-db.json -database-dsn='postgres://postgres:password@localhost:5433/postgres?sslmode=disable' -test.run=^TestIteration13$
```

```shell
./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -server-port=8080 -file-storage-path=/tmp/metrics-db.json -database-dsn='postgres://postgres:password@localhost:5433/postgres?sslmode=disable' -key='test123' -test.run=^TestIteration14$
```

## GoDoc

```shell
godoc -http=:8080 -play
```

Ссылка для просмотра документации приложения:

http://localhost:8080/pkg/github.com/havilcorp/yandex-go-musthave-metrics-tpl/?m=all

где

- /pkg/github.com/havilcorp/yandex-go-musthave-metrics-tpl - пакет с исходным кодом
- /?m=all флаг просмотра скрытых пакетов, таких как internal

## mockery

Предоставляет возможность легко создавать макеты для интерфейсов Golang

```shell
docker run -v "$PWD":/src -w /src vektra/mockery --all
```

## Покрытие

Команды для проверки процента покрытия тестами

```shell
go test -coverprofile=coverage.out ./internal...
go tool cover -func=coverage.out
```

## Swagger

```shell
swag init --dir ./internal/handlers
```
