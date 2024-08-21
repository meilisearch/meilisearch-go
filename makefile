.PHONY: test easyjson requirements

easyjson:
	easyjson -all types.go

test:
	docker compose run --rm package bash -c "go get && golangci-lint run -v && go test -v"

requirements:
	go get github.com/mailru/easyjson && go install github.com/mailru/easyjson/...@latest
	curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh
	go get -v -t ./...
