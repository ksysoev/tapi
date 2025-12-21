.DEFAULT_GOAL := help

help: ## Show this help message
	@awk 'BEGIN {FS = ":.*## "; printf "\nUsage:\n  make <target>\n\nTargets:\n"} \
		/^([a-zA-Z_-]+):.*## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary
	go build -o bin/tapi ./cmd/tapi

install: ## Install the binary
	go install ./cmd/tapi

test: ## Run unit tests with race detector
	go test -race ./...

lint: ## Run golangci-lint
	golangci-lint run

tidy: ## Run go mod tidy
	go mod tidy

fmt: ## Format code with gofmt
	gofmt -w .

run: build ## Build and run the TUI
	./bin/tapi explore

clean: ## Clean build artifacts
	rm -rf bin/
