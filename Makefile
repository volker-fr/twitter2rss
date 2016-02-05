build: clean
	go build

run: build
	./twitter2rss -config twitter2rss.hcl

clean:
	rm -rf twitter2rss

linux: clean
	GOOS=linux GOARCH=amd64 go build

fmt:
	go fmt *.go
	go fmt config/*.go
	go fmt filter/*.go
