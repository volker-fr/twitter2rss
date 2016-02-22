FOLDERS = config filter parser feed

build: clean
	go build

run: build
	./twitter2rss -config twitter2rss.hcl

clean:
	rm -rf twitter2rss

linux: clean
	GOOS=linux GOARCH=amd64 go build

freebsd: clean
	GOOS=freebsd GOARCH=amd64 go build

fmt:
	@go fmt *.go
	@for dir in $(FOLDERS); do \
		cd $$dir && go fmt *.go; \
		cd ..; \
	done

debug: build
	./twitter2rss -config twitter2rss.hcl -max-tweets 5 -debug

vet:
	@go tool vet .

test: vet
	@go test -cover -v
	@for dir in $(FOLDERS); do \
		cd $$dir && echo && echo Running test in $$dir && go test -cover; \
		cd ..; \
	done
