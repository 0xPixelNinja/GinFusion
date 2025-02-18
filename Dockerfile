# Build stage using golang alpine image with required packages
FROM golang:1.22.6-alpine AS builder

# Install gcc, musl-dev, and sqlite-dev for CGO and go-sqlite3 to work
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copy go.mod and go.sum and download dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY . .

# Build the binary with CGO enabled.
RUN CGO_ENABLED=1 go build -o webapp ./cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates and SQLite runtime library if needed.
RUN apk --no-cache add ca-certificates sqlite

WORKDIR /root/

# Copy the binary and configuration file from the builder stage.
COPY --from=builder /app/webapp .
COPY config.yaml .

# Expose ports for both the API and Admin servers.
EXPOSE 8080 8081

# Run the application.
CMD ["./webapp"]