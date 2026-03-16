FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o evm-api ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/evm-api .

EXPOSE 8080

CMD ["./evm-api"]
