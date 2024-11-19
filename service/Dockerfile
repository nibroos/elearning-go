# Production Dockerfile
FROM golang:1.23-alpine as builder

# Install necessary packages including Make and migrate
RUN apk add --no-cache make curl \
    && go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application as a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main . && ls -l /app/main

# Final minimal image
FROM alpine:3.18 as final

# Set the timezone to GMT+7 and install necessary certificates
RUN apk --no-cache add ca-certificates tzdata \
    && cp /usr/share/zoneinfo/Asia/Bangkok /etc/localtime \
    && echo "Asia/Bangkok" > /etc/timezone \
    && apk del tzdata

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/main /app/main

# Copy the migrate CLI
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

# Set the working directory
WORKDIR /app

# Create a public directory for temporary files
RUN mkdir -p /app/public/pdf /app/public/csv /app/public/pictures /app/public/videos \
    && chmod -R 755 /app/public

# Drop privileges by creating a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Expose necessary ports
EXPOSE 4001 50051

# Command to run the application
CMD ["/app/main"]