FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o consumer ./cmd/consumer/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/consumer .

RUN mkdir -p /app/logs && chmod 755 /app/logs

CMD ["./consumer"]
