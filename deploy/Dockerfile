FROM golang:1.18-alpine

RUN apk add --update-cache \
    build-base \
    ca-certificates \
    openssl

RUN mkdir -p /var/usr/local/xyo/sdk

COPY . /var/usr/local/xyo/sdk

WORKDIR /var/usr/local/xyo/sdk

RUN openssl req -x509 -newkey rsa:4096 -nodes -out server.crt -keyout server.key \
    -days 365 \
    -subj "/C=GB/O=localhost/OU=Domain Control Validated/CN=localhost"

ENV CGO_ENABLED=0

#RUN go test ./... -v
