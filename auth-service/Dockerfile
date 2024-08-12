# auth-service/Dockerfile
FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/cosmtrek/air@v1.29.0

CMD ["air"]
