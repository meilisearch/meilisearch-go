.PHONY: test requirements mock

mock:
	@echo "Generating mocks..."
	@mkdir -p mocks
	@mockery
	@echo "✓ Mock generation complete"

test:
	docker compose run --rm package bash -c "go get && golangci-lint run -v && go test -v ./..."

requirements:
	curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh
	go get -v -t ./...
	go install github.com/vektra/mockery/v3@v3.6.1
