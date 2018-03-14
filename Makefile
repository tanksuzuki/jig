.PHONY: help
.DEFAULT_GOAL := help

release: ## release build
	mkdir -p out
	rm -rf out/*
	GOOS=linux GOARCH=amd64 go build -o jig
	zip out/jig_linux_amd64.zip jig
	rm -rf jig
	GOOS=linux GOARCH=386 go build -o jig
	zip out/jig_linux_386.zip jig
	rm -rf jig
	GOOS=darwin GOARCH=amd64 go build -o jig
	zip out/jig_macosx_amd64.zip jig
	rm -rf jig
	GOOS=darwin GOARCH=386 go build -o jig
	zip out/jig_macosx_386.zip jig
	rm -rf jig

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
