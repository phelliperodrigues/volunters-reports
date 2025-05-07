# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o ccb-report cmd/api/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/ccb-report .

# Copy example files
COPY --from=builder /app/files/books.csv.example /app/files/books.csv
COPY --from=builder /app/files/input.csv.example /app/files/input.csv

# Create necessary directories
RUN mkdir -p /app/files/output /app/data/reports

# Set environment variables
ENV BASE_PATH=/app
ENV SERVER_PORT=8080

# Expose the port
EXPOSE 8080

# Run the binary
CMD ["./ccb-report"] 