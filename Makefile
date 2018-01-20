# ref. http://postd.cc/auto-documented-makefile/
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build `thracia` binary
	rm -f assets.go
	go-assets-builder -p thracia assets > assets.go
	go build cmd/thracia/thracia.go

.PHONY: install
install: ## Install requirements
	glide install
	go get -u github.com/jessevdk/go-assets-builder
