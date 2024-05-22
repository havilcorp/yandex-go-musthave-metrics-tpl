# go-musthave-metrics-tpl

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.

### Запуск сервера

```shell
go run cmd/server/main.go
```

Для запуска сервера в режиме gRPC необходимо:

1. либо в конфиге указать address_grpc
2. либо передать флаг address_grpc
3. либо указать переменную окружения (env) ADDRESS_GRPC

###### Конфигурация

Флаги

- -h - показать все команды
- -a - адрес и порт сервера
- -address_grpc - адрес grpc
- -i - интевал сохранения метрик в файл
- -f - файл для созранения метрик. По деволту: /tmp/metrics-db.json
- -r - загружать ли при запуске метрики из файла
- -d - строка подключения к базе данных
- -k - ключ sha256
- -crypto-key - путь к файлу с приватным ключем для расшифрования сообщения
- -c - путь к файлу конфигов
- -t - доверенная маска подсети, например: 192.168.0.0/24

Env

- ADDRESS - адрес и порт сервера
- ADDRESS_GRPC - адрес и порт GRPC сервера
- STORE_INTERVAL - интевал сохранения метрик в файл
- FILE_STORAGE_PATH - файл для созранения метрик. По деволту: /tmp/metrics-db.json
- RESTORE - загружать ли при запуске метрики из файла
- DATABASE_DSN - строка подключения к базе данных
- KEY - ключ sha256
- CRYPTO_KEY - путь до файла с приватным ключем для расшифрования сообщения
- CONFIG - путь к файлу конфигов
- TRUSTED_SUBNET - доверенная маска подсети, например: 192.168.0.0/24

```go
postgres://postgres:password@localhost:5433/postgres?sslmode=disable
```

### Запуск агента

```shell
go run cmd/agent/main.go
```

Для запуска агента в режиме gRPC необходимо:

1. либо в конфиге указать address_grpc
2. либо передать флаг address_grpc
3. либо указать переменную окружения (env) ADDRESS_GRPC

###### Конфигурация

Флаги

- -a - адрес и порт сервера
- -address_grpc - адрес grpc
- -r - интервал отправки метрик на сервер
- -p - интервал сбора метрик
- -k - ключ sha256
- -l - лимит запросов
- -crypto-key - путь к файлу с публичным ключем для шифрования сообщения
- -crypto-crt - путь к файлу с сертификатом
- -c - путь к файлу конфигов

Env

- ADDRESS - адрес и порт сервера
- ADDRESS_GRPC - адрес и порт GRPC сервера
- REPORT_INTERVAL - интервал отправки метрик на сервер
- POLL_INTERVAL - интервал сбора метрик
- KEY - ключ sha256
- RATE_LIMIT - лимит запросов
- CRYPTO_KEY - путь до файла с публичным ключем для шифрования сообщения
- CRYPTO_CRT - путь к файлу с сертификатом
- CONFIG - путь к файлу конфигов

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

## TLS защищенное соединение

Зайдите в корневую директорию, затем в папке tls и там пропишите:

```shell
openssl genrsa -out ca.key 2048
openssl req -new -x509 -days 365 -key ca.key -subj "/C=CN/ST=GD/L=SZ/O=Acme, Inc./CN=Acme Root CA" -out ca.crt
openssl req -newkey rsa:2048 -nodes -keyout server.key -subj "/C=CN/ST=GD/L=SZ/O=Acme, Inc./CN=localhost" -out server.csr
openssl x509 -req -extfile <(printf "subjectAltName=DNS:localhost") -days 365 -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt
openssl rsa -in server.key -pubout > key.pub
```

Для запуска gRPC в режиме tls необходимо:

1. Для агента указать путь crypto_crt как "./tls/ca.crt"
2. Для сервера указать путь crypto_key как "./tls/server.key"

Для шифрования трафика по протоколу REST между агентом и сервером по ключу RSA необходимо:

1. Для агента указать путь crypto_key как "./tls/key.pub"
2. Для сервера указать путь crypto_key как "./tls/server.key"

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
go build -o cmd/agent/agent cmd/agent/main.go
go build -o cmd/server/server cmd/server/main.go
```

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

## GRPC Генерация прото файлов

protoc --go_out=. --go_opt=paths=source_relative \
 --go-grpc_out=. --go-grpc_opt=paths=source_relative \
 pkg/proto/metric/metric.proto
