.PHONY: all fmt fmt-check lint test check build ssh

all: fmt lint test

fmt:
	gofmt -w .

fmt-check:
	@if [ -n "$$(gofmt -l .)" ]; then echo "Go code is not formatted. Please run 'make fmt' or 'go fmt ./...'."; exit 1; fi

lint:
	go vet ./...

test:
	go test ./... -v

check: fmt-check lint test

build:
	docker build -f deploy/Dockerfile . -t sdk-go:latest --no-cache

ssh:
	docker run -it --rm --name XYO_financial_SDK_Golang \
		--add-host api.xyo.financial:127.0.0.1 \
		sdk-go:latest sh -c "go test ./..."

