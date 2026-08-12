.PHONY: build fmt lint test check ssh

build:
	docker build -f deploy/Dockerfile . -t sdk-go:latest --no-cache

fmt:
	go fmt .

lint:
	go vet ./...

test:
	go test ./... -v

check: fmt lint test

ssh:
	docker run -it --rm --name XYO_financial_SDK_Golang \
		--add-host api.xyo.financial:127.0.0.1 \
		sdk-go:latest sh -c "go test ./..."
