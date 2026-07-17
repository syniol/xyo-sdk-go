FROM golang:1.18-alpine

RUN mkdir -p /var/usr/local/xyo/sdk

COPY . /var/usr/local/xyo/sdk

WORKDIR /var/usr/local/xyo/sdk

ENV CGO_ENABLED=0

RUN apk add --update-cache \
    curl

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.50.1 \
    && golangci-lint --version

# Format check
RUN if [ -n "$(gofmt -l .)" ]; then echo "Go code is not formatted. Please run go fmt."; exit 1; fi

# Static Analysis
RUN go vet ./...

# Running linting and Tests
RUN golangci-lint run ./... -v && go test ./... -v

# Testing import and instantiation with an example package
WORKDIR /var/usr/local/xyo/sdk/example
RUN go mod download && go run main.go

WORKDIR /var/usr/local/xyo/sdk
