.PHONY: help
.DEFAULT_GOAL := help

release: ## release build
	@GOOS=linux GOARCH=amd64 go build -o jig_linux_amd64
	@GOOS=darwin GOARCH=amd64 go build -o jig_maxosx_amd64

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
