# Build stage
FROM golang:1.22.3-alpine AS builder
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /myapp

# Production stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /myapp .
COPY .env .

EXPOSE 8080
CMD ["./myapp"]
