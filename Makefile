build:
	docker build -f deploy/Dockerfile . -t xyo-sdk-go:latest --no-cache


run:
	docker run -it --rm --name XYoFinancialSDKGolang xyo-sdk-go:latest sh
