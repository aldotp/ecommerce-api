_Ecommerce-Go-Api_

### Command for running this program:

Copy file .env.example to .env and edit:

```shell
$ cp .env.example .env
```

### RUN WITH DOCKER

This command build and run container

## Run Docker (setup database and redis manually)

```shell
// build process
$ docker build --rm --tag ecommerce-go-api:latest -f Dockerfile .
// run
$ docker run --env-file .env --rm -p 8080:8080 --name ecommerce-go-api ecommerce-go-api:latest
```

## Run Docker using Docker Compose

```shell
// running docker compose in background
$ docker compose up -d
```

## Run Migrations

```shell
# Up the migration
migrate -path ./internal/adapter/storage/postgres/migrations \
        -database "postgres://postgres:12345678a@127.0.0.1:5432/ecommerce-go-api?sslmode=disable" \
        -verbose up


# Down the migration
migrate -path ./internal/adapter/storage/postgres/migrations \
        -database "postgres://postgres:12345678a@127.0.0.1:5432/ecommerce-go-api?sslmode=disable" \
        -verbose down
```

## New Migrartions

```shell
// add new migration
migrate create -ext sql -dir  ./internal/adapter/storage/postgres/migrations -format "20060102150405" add_table_carts
```

#### Generate Swagger

```shell
swag init -g cmd/main.go
```

### Open Swagger Documentation

http://localhost:8080/docs/index.html
