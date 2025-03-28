# STAGE 1 
FROM golang:1.24-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/bin/ecommerce-go-api ./cmd/main.go


# STAGE 2
FROM alpine:3.18.4

RUN apk add --no-cache tzdata

ARG APP_VERSION
ENV APP_VERSION=$APP_VERSION \
    TZ=Asia/Jakarta \
    APP_PORT=80

WORKDIR /app

COPY --from=builder /app/bin/ecommerce-go-api ./api

EXPOSE 80

ENTRYPOINT ["./api", "rest"]
