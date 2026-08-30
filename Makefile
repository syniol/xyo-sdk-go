.PHONY: build fmt lint test check ssh

build:
	docker build -f deploy/Dockerfile . -t sdk-go:latest --no-cache

fmt:
	go fmt .

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		go vet ./...; \
	fi

test:
	go test -race ./... -v

check: fmt lint test

ssh:
	docker run -it --rm --name XYO_financial_SDK_Golang \
		--add-host api.xyo.financial:127.0.0.1 \
		sdk-go:latest sh -c "go test ./..."
