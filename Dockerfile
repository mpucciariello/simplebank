# Build stage
FROM golang:1.20-alpine3.17 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY config.env .
COPY db/migration ./migration
COPY start.sh .

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]

# docker build -t simplebank:latest .
# docker run --name simplebank --network bank_network -p 8080:8080 -e DB_SOURCE="postgresql://root:secret@postgres12:5432/simple_bank?sslmode=disable" simplebank:latest
