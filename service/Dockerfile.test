FROM golang:1.23-alpine AS builder

RUN apk add --no-cache make curl

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

# Copy the go binary and set GOROOT
COPY --from=builder /usr/local/go /usr/local/go
ENV GOROOT=/usr/local/go
ENV PATH=$PATH:/usr/local/go/bin

# Create a public directory for temporary files
RUN mkdir -p /apps/public/pdf /apps/public/csv /apps/public/pictures /apps/public/videos \
    && chmod -R 755 /apps/public

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /workdir/app /bin/app
COPY --from=builder /workdir /apps

# Copy the .env file
COPY .env /apps/.env

# Copy the make command
COPY --from=builder /usr/bin/make /usr/bin/make

# Create a writable home directory for the non-root user
RUN mkdir -p /home/appuser && chown -R appuser:appgroup /home/appuser
ENV HOME=/home/appuser

RUN mkdir -p /go && chown -R appuser:appgroup /go
# USER appuser

ENTRYPOINT ["/apps/app"]