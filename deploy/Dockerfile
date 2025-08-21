FROM golang:1.18-alpine

RUN mkdir -p /var/usr/local/xyo/sdk

COPY . /var/usr/local/xyo/sdk

WORKDIR /var/usr/local/xyo/sdk

ENV CGO_ENABLED=0

RUN go test ./... -v

# Testing example package for import and instantiation
WORKDIR /var/usr/local/xyo/sdk/example
RUN go mod download && go run main.go

WORKDIR /var/usr/local/xyo/sdk
