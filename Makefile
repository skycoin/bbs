.DEFAULT_GOAL := help
.PHONY: run run-help test build package clean help

# Static files directory
STATIC_DIR = static
STATIC_DIST_DIR = $(STATIC_DIR)/dist

# Package files directory
PACKAGE_DIR = pkg
PACKAGE_BUILD_DIR = $(PACKAGE_DIR)/build

run:  ## Runs the BBS node. To add arguments, do 'make run ARGS="--foo" run'.
	go run cmd/bbsnode/bbsnode.go --http-gui-dir="./${STATIC_DIST_DIR}" ${ARGS}

run-help: ## Shows BBS node help.
	@go run cmd/bbsnode/bbsnode.go --help

test: ## Run tests.
	go test ./cmd/...
	go test ./src/...

build: ## Build static files.
	cd $(STATIC_DIR) && npm install -g @angular/cli@latest && yarn install && ng build --target=production
	
package: ## Builds static and binaries into zip files located in pkg/build directory.
	cd $(PACKAGE_DIR) && bash package.sh

clean: ## Cleans static and built package files.
	rm -rf $(STATIC_DIST_DIR)
	rm -rf $(PACKAGE_BUILD_DIR)

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
