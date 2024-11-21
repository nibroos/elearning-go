FROM golang:1.23-alpine as builder

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
RUN apk add --no-cache ca-certificates

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /workdir/app /bin/app

USER 65534

ENTRYPOINT ["/bin/app"]