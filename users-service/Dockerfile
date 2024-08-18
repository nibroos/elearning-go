# Start with a minimal base image
FROM golang:1.20-alpine as builder

# Install Air
RUN go install github.com/cosmtrek/air@latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the working directory inside the container
COPY . .

# Build the Go app
RUN go build -o main ./cmd/main.go

# Start a new stage from scratch
FROM alpine:3.18

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main /app/main

# Copy Air config file and install Air
COPY --from=builder /go/bin/air /usr/bin/air
COPY .air.toml /app/.air.toml

WORKDIR /app

# Command to run the executable
CMD ["air"]
