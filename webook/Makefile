.PHONY: docker
docker:
	@rm webook || true
	@GOOS=linux GOARCH=amd64 go build -o webook .
	@docker rmi -f liyilo/webook-live:v0.0.1
	@docker build -t liyilo/webook-live:v0.0.1 .