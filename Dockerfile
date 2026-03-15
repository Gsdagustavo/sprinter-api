FROM golang:1.26 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM scratch

WORKDIR /app

ENV CONFIGS_PATH=/app/settings/dev-settings.toml

COPY --from=build /app/server .
COPY --from=build /app/settings ./settings
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8000

CMD ["./server"]