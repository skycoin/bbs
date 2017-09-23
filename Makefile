.DEFAULT_GOAL := help
.PHONY: run run-help test lint check release clean help

# Static files directory
STATIC_DIR = static/dist

run:  ## Run the skycoin node. To add arguments, do 'make ARGS="--foo" run'.
	go run cmd/bbsnode/bbsnode.go --http-gui-dir="./${STATIC_DIR}" ${ARGS}

run-help: ## Show skycoin node help
	@go run cmd/bbsnode/bbsnode.go --help

test: ## Run tests
	go test ./cmd/...
	go test ./src/...

lint: ## Run linters. requires vendorcheck, gometalinter, golint, goimports
	gometalinter --disable-all -E goimports --tests --vendor ./...
	vendorcheck ./...

check: lint test ## Run tests and linters

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
