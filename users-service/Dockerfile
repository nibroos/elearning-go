# Production Dockerfile

FROM golang:1.22-alpine as builder

# Install necessary packages including Make
RUN apk add --no-cache make

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application as a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# Final minimal image
FROM alpine:3.18 as final

# Install necessary certificates and Make
RUN apk --no-cache add ca-certificates make

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/main /app/main

# Copy the migrate CLI
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate

# Set the working directory
WORKDIR /app

# Drop privileges by creating a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Expose necessary ports
EXPOSE 4001 50051

# Command to run the application
CMD ["/app/main"]
