# Step 1: Build the application
FROM golang:1.20.2-alpine3.17 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bin/challenge ./cmd/challenge

# Step 2: Create a smaller container for deployment
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/bin/challenge .
COPY --from=builder /app/config.toml .
COPY --from=builder /app/data/ ./data/

ENV CHALLENGE_REDIS_SEARCH_URL=""

EXPOSE 8080

CMD ["./challenge"]
