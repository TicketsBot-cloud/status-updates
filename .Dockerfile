# Build container
FROM golang:alpine AS builder

RUN apk update && apk upgrade && apk add git zlib-dev gcc musl-dev

COPY . /go/src/github.com/TicketsBot-cloud/status-updates
WORKDIR /go/src/github.com/TicketsBot-cloud/status-updates

RUN set -Eeux && \
    go mod download && \
    go mod verify

RUN GOOS=linux GOARCH=amd64 \
    go build \
    -tags=jsoniter \
    -trimpath \
    -o main cmd/status-updates/main.go

# Prod container
FROM alpine:latest

RUN apk update && apk upgrade && apk add curl

COPY --from=builder /go/src/github.com/TicketsBot-cloud/status-updates/main /srv/status-updates/main

RUN chmod +x /srv/status-updates/main

RUN adduser container --disabled-password --no-create-home
USER container
WORKDIR /srv/status-updates

CMD ["/srv/status-updates/main"]