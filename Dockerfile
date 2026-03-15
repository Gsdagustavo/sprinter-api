# Create server executable
FROM golang:1.26 AS build

WORKDIR /app

COPY . /app

RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# Run server
FROM scratch

WORKDIR /app

COPY --from=build /app/server ./

EXPOSE 8000

CMD [ ".server" ]