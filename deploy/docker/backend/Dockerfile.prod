# Backend Dockerfile
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy source code
COPY backend/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/api

# Final stage
FROM alpine:3.22

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy binary
COPY --from=builder /app/bin/api .

# Note: Config is loaded from environment variables via Viper
# No config.yaml file needed - all config comes from env vars

EXPOSE 8080

CMD ["./api"]
