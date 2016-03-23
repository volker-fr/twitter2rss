FOLDERS = config filter parser feed

# Via http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: clean ## Compile go for this platform
	go build

run: build ## Compile & run the application with the default coniguration file
	./twitter2rss -config twitter2rss.hcl

clean: ## Wipe the compiled executable
	rm -rf twitter2rss

linux: clean ## (Cross) compile it for Linux amd64 platform
	GOOS=linux GOARCH=amd64 go build

freebsd: clean ## (Cross) compile it for FreeBSD amd64 platform
	GOOS=freebsd GOARCH=amd64 go build

fmt: ## Format source
	@go fmt *.go
	@for dir in $(FOLDERS); do \
		cd $$dir && go fmt *.go; \
		cd ..; \
	done

debug: build ## Build and run it in debug mode
	./twitter2rss -config twitter2rss.hcl -max-tweets 5 -debug

vet: ## Run "go tool vet"
	@go tool vet .

test: vet ## Run tests
	@go test -cover -v
	@for dir in $(FOLDERS); do \
		cd $$dir && echo && echo Running test in $$dir && go test -cover; \
		cd ..; \
	done
