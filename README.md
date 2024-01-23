# go-musthave-metrics-tpl

Шаблон репозитория для трека «Сервер сбора метрик и алертинга».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

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

## Тесты

1. ./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -test.run=^TestIteration1$

2. ./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -test.run=^TestIteration2A$

./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -test.run=^TestIteration2A$

3. ./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -test.run=^TestIteration3A$

./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -test.run=^TestIteration3B$

4. ./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -server-port=8080 -test.run=^TestIteration4$

5. ./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -server-port=8080 -test.run=^TestIteration5$

6. ./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -server-port=8080 -test.run=^TestIteration6$

7. ./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -server-port=8080 -test.run=^TestIteration7$

8. ./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -server-port=8080 -test.run=^TestIteration8$

9. ./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -server-port=8080 -file-storage-path=/tmp/metrics-db.json -test.run=^TestIteration9$

10. ./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -server-port=8080 -file-storage-path=/tmp/metrics-db.json -database-dsn='postgres://postgres:password@localhost:5433/postgres?sslmode=disable' -test.run=^TestIteration10A$

./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -server-port=8080 -file-storage-path=/tmp/metrics-db.json -database-dsn='postgres://postgres:password@localhost:5433/postgres?sslmode=disable' -test.run=^TestIteration10B$

11. ./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -server-port=8080 -file-storage-path=/tmp/metrics-db.json -database-dsn='postgres://postgres:password@localhost:5433/postgres?sslmode=disable' -test.run=^TestIteration11$

12. ./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -server-port=8080 -file-storage-path=/tmp/metrics-db.json -database-dsn='postgres://postgres:password@localhost:5433/postgres?sslmode=disable' -test.run=^TestIteration12$

13. ./metricstest-darwin-arm64 -test.v -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -server-port=8080 -file-storage-path=/tmp/metrics-db.json -database-dsn='postgres://postgres:password@localhost:5433/postgres?sslmode=disable' -test.run=^TestIteration13$

## Запуск автотестов 1

/Users/kotvkompe/Desktop/YP/yandex-go-musthave-metrics-tpl
./metricstest-darwin-arm64 -test.v -test.run=^TestIteration7$ -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -source-path=. -server-port=8081

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).
