# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build all three binaries
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/migrations ./cmd/migrations
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/wheelgpt ./cmd/wheelgpt

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy all binaries
COPY --from=builder /app/bin/migrations ./
COPY --from=builder /app/bin/api ./
COPY --from=builder /app/bin/wheelgpt ./

RUN chmod +x ./migrations ./api ./wheelgpt

# Default to API
CMD ["./api"]