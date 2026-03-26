FROM golang:1.26.1-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git

COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o pack-api ./cmd/main.go

FROM alpine:latest
WORKDIR /app

# Copy the built binary from builder
COPY --from=builder /app/pack-api .

# Copy migrations if needed
COPY migrations ./migrations

EXPOSE 8080

# Run the binary
CMD ["./pack-api"]
