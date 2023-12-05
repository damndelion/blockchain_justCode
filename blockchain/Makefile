include .env.example
export

LOCAL_BIN:=$(CURDIR)/bin
PATH:=$(LOCAL_BIN):$(PATH)

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

compose-up: ### Run docker-compose
	docker-compose up
.PHONY: compose-up


swag-v1-auth: ### swag init for auth service
	 cd internal/auth/controller/http/v1 && swag init -g ../../../../../cmd/auth/main.go -o  ../../../../../docs/auth
.PHONY: swag-v1-auth

swag-v1-user: ### swag init for user service
	 cd internal/user/controller/http/v1 && swag init -g ../../../../../cmd/user/main.go -o  ../../../../../docs/user
.PHONY: swag-v1-user

swag-v1-blockchain: ### swag init for blockchain service
	 cd internal/blockchain/controller/http/v1 && swag init -g ../../../../../cmd/blockchain/main.go -o  ../../../../../docs/blockchain
.PHONY: swag-v1-blockchain

linter-golangci: ### check by golangci linter
	golangci-lint run --fix
.PHONY: linter-golangci

mockery-blockchain: ### create mock for blockchain
	mockery --name ChainRepo --dir internal/blockchain/usecase --output internal/blockchain/mocks
.PHONY: mockery-blockchain

test-blockchain: ### run test
	cd internal/blockchain/usecase && go test
.PHONY: test-blockchain


docker-build-user:
	docker build -t user -f Dockerfile-user .
.PHONY: docker-build-user

docker-build-auth:
	docker build -t auth -f Dockerfile-auth .
.PHONY: docker-build-auth

 docker-build-blockchain:
	docker build -t blockchain -f Dockerfile-blockchain .
 .PHONY: docker-build-blockchain

