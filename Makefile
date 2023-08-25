GO_INSTALLED := $(shell which go)
GO_FILES = $(shell go list ./... | grep -v /vendor/ | grep -v /api/ | grep -v /cmd/)

MOCKERY = domain/transaction.go

.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

check: ## Check if all required binaries are installed correctly
ifndef GO_INSTALLED
	$(error "go is not installed, please run 'brew install go'")
endif

tidy: ## Vendor and Tidy go modules
	@rm -rf vendor
	@go mod tidy
	@go mod vendor

gen: check ## Run code generator, this generate server and clients code from spec in /api/ folder
	@go mod tidy
	@go mod vendor
	@echo "Success! Generated"

build: gen ## Build application binary
	@rm -rf bin
	@mkdir -p bin
	@go build -o bin/app ./app
	@echo "Success! Binaries can be found in 'bin' dir"

run: ## Run application on default port 8080 and print console format
	@go run app/main.go

mock-gen: $(MOCKERY)
	go generate ./...

test: ## Run go test with coverage and race detection
	@go test $(GO_FILES) -cover --race

lint:
	@golangci-lint run ./...