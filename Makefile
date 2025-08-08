build:
	docker build -f deploy/Dockerfile . -t xyo-sdk-go:latest --no-cache


run:
	docker run -it --rm --name XYO_financial_SDK_Golang \
		--add-host api.xyo.financial:127.0.0.1 \
		xyo-sdk-go:latest sh
