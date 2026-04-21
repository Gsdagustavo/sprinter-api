FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o sprinter-api .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/sprinter-api .

ENV SERVER_PORT=8000
ENV SERVER_HOST=localhost
ENV SERVER_DOMAIN=http://localhost:8000

ENV DB_HOST=localhost
ENV DB_PORT=3306
ENV DB_NAME=sprinter_db
ENV DB_USER=root
ENV DB_PASSWORD=root

ENV ENVIRONMENT_TYPE=development

ENV LOG_DIR=./runtime/logs

ENV PASETO_SECURITY_KEY=b7e16bf2-ab0f-4a68-a16b-f2ab0f6a

ENV CORS_ORIGINS=localhost:8000

ENV FILE_STORAGE_FOLDER=./runtime/storage

ENV SMTP_HOST=localhost
ENV SMTP_PORT=1025
ENV SMTP_USER=test
ENV SMTP_PASSWORD=test
ENV SMTP_FROM=dev@test.com
ENV SMTP_MAX_CONNECTIONS=5
ENV SMTP_IDLE_TIMEOUT=30s
ENV SMTP_POOL_WAIT_TIMEOUT=10s

EXPOSE 8000

CMD ["./sprinter-api"]