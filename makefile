.PHONY: test requirements

test:
	docker compose run --rm package bash -c "go get && golangci-lint run -v && go test -v ./..."

requirements:
	curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh
	go get -v -t ./...
