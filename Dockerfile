FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ai-proxy ./cmd/ai-router/main.go

FROM gcr.io/distroless/base:latest

WORKDIR /app

COPY --from=builder /app/ai-proxy /app/ai-proxy

EXPOSE 8080

ENTRYPOINT ["/app/ai-proxy"]
