MODULE_NAME=$(shell grep ^module go.mod | cut -d " " -f2)
GIT_COMMIT_HASH=$(shell git rev-parse HEAD)
LD_FLAGS=-ldflags="-X $(MODULE_NAME)/internal/config.gitCommitHash=$(GIT_COMMIT_HASH)"

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: api
api:
	@go run $(LD_FLAGS) cmd/api/main.go 

.PHONY: build
build:
	@go build -o ./bin/api $(LD_FLAGS) ./cmd/api

.PHONY: docker-stop
docker-stop:
	@docker-compose -f ./build/package/docker-compose.yaml stop

.PHONY: docker
docker:
	@docker-compose -f ./build/docker-compose.yaml up

.PHONY: generate
generate:
	@go generate ./...