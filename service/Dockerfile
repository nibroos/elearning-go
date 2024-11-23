FROM golang:1.23-alpine AS builder

RUN apk add --no-cache curl

# Install golang-migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.2

COPY . /workdir
WORKDIR /workdir

ENV CGO_CPPFLAGS="-D_FORTIFY_SOURCE=2 -fstack-protector-all"
ENV GOFLAGS="-buildmode=pie"

# Specify output binary name explicitly
RUN go build -ldflags "-s -w" -trimpath -o app .

FROM alpine:latest

# Install necessary packages
RUN apk --no-cache add ca-certificates tzdata \
    && cp /usr/share/zoneinfo/Asia/Bangkok /etc/localtime \
    && echo "Asia/Bangkok" > /etc/timezone \
    && apk del tzdata

# Copy the migrate CLI
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

# Create a public directory for temporary files
RUN mkdir -p /app/public/pdf /app/public/csv /app/public/pictures /app/public/videos \
    && chmod -R 755 /app/public

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /workdir/app /bin/app
COPY --from=builder /workdir/app /test/app
COPY --from=builder /workdir/app /app/app
COPY --from=builder /workdir /apps

# Copy the .env file
COPY .env /apps/.env
COPY .env /app/.env

USER 65534

ENTRYPOINT ["/app/app"]