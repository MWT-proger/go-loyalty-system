# Накопительная система лояльности «Гофермарт»

[Техническое задание](docs/SPECIFICATION.md) для ознакомления с проектом

[Документация OpenAPI mwt-proger.github.io/go-loyalty-system](https://mwt-proger.github.io/go-loyalty-system/) запущена для работы с localhost


## Развертывание проекта

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.

```bash
git clone https://github.com/MWT-proger/go-loyalty-system.git
```


2. Скопируйте шаблон файла с переменным окружения

```bash
  cp deployments/.env.example deployments/.env
```

3. Укажите верные переменные окружения в только что созданный файл [.env](deployments/.env)

*Доступны следующие переменные*
```bash
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=testDB
POSTGRES_PORT=5432
```
4. Запустите БД Postgres следующей командой

```bash
  docker compose -f deployments/docker-compose.yaml --env-file deployments/.env up
```
5. Запустите систему лояльности
```
go run ./cmd/gophermart -a "localhost:7000" -d "user=postgres password=postgres host=localhost port=5432 dbname=testDB sslmode=disable"
```