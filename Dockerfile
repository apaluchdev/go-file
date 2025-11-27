FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

# Create storage directory
RUN mkdir -p /app/storage

# Expose port
EXPOSE 8080

# Set environment variable for storage path
ENV STORAGE_PATH=/app/storage

CMD ["./main"]
