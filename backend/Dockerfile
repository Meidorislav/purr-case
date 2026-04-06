FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY internal/migrations ./internal/migrations
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/purr-case/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/internal/migrations ./internal/migrations
EXPOSE 8080
CMD ["./server"]
