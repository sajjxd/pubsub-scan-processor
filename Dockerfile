FROM golang:1.20-alpine AS builder

WORKDIR /src

# Needed for SQLite
RUN apk add --no-cache gcc musl-dev

# Build
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o /app/processor ./cmd/processor

# Copy binary into slim image
FROM alpine
WORKDIR /app

# To make secure TLS connections
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/processor .
CMD ["/app/processor"]